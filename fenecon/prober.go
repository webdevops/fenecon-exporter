package fenecon

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	resty "github.com/go-resty/resty/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/remeh/sizedwaitgroup"
	"go.uber.org/zap"
)

type (
	FeneconProber struct {
		ctx      context.Context
		client   *resty.Client
		logger   *zap.SugaredLogger
		registry *prometheus.Registry

		parallelRequests int

		target string

		prometheus feneconMetrics
	}
)

func New(ctx context.Context, registry *prometheus.Registry, logger *zap.SugaredLogger) *FeneconProber {
	fp := FeneconProber{}
	fp.ctx = ctx
	fp.registry = registry
	fp.logger = logger
	fp.parallelRequests = 5
	fp.initResty()
	fp.initMetrics()

	return &fp
}

func (fp *FeneconProber) initResty() {
	fp.client = resty.New()
	fp.client.RetryCount = 3
	fp.client.RetryWaitTime = 2 * time.Second
	fp.client.RetryMaxWaitTime = 5 * time.Second
	fp.client.AddRetryAfterErrorCondition()

	fp.client.OnAfterResponse(func(client *resty.Client, response *resty.Response) error {
		switch statusCode := response.StatusCode(); statusCode {
		case 401:
			return errors.New(`fenecon requires authentication or credentials are invalid`)
		case 200:
			// all ok, proceed
			return nil
		default:
			return fmt.Errorf(`expected http status 200, got %v`, response.StatusCode())
		}
	})
}

func (fp *FeneconProber) SetUserAgent(val string) {
	fp.client.SetHeader("User-Agent", val)
}

func (fp *FeneconProber) SetParallelRequests(val int) {
	fp.parallelRequests = val
}

func (fp *FeneconProber) SetTimeout(timeout time.Duration) {
	fp.client.SetTimeout(timeout)
}

func (fp *FeneconProber) SetHttpAuth(username, password string) {
	fp.client.SetDisableWarn(true)
	fp.client.SetBasicAuth(username, password)
}

func (fp *FeneconProber) Run(target string) {
	fp.target = target
	fp.logger.With(zap.String("target", target))

	startTime := time.Now()
	fp.logger.Debugf(`start probe`)

	client := fp.client.SetBaseURL(
		fmt.Sprintf(`%s/rest/channel/`, strings.TrimRight(target, "/")),
	)

	commonLabels := prometheus.Labels{"target": target}
	phase1Labels := prometheus.Labels{"target": target, "phase": "1"}
	phase2Labels := prometheus.Labels{"target": target, "phase": "2"}
	phase3Labels := prometheus.Labels{"target": target, "phase": "3"}

	fp.prometheus.info.With(commonLabels).Set(1)

	wg := sizedwaitgroup.New(5)

	// general
	fp.queryCommon(&wg, client, "_sum/State", commonLabels, fp.prometheus.status)

	// battery
	fp.queryCommon(&wg, client, "_sum/EssSoc", commonLabels, fp.prometheus.batteryCharge)
	fp.queryCommon(&wg, client, "_sum/EssActivePower", commonLabels, fp.prometheus.batteryPower)
	fp.queryCommonIfNotZero(&wg, client, "_sum/EssActiveChargeEnergy", commonLabels, fp.prometheus.batteryPowerChargeTotal)
	fp.queryCommonIfNotZero(&wg, client, "_sum/EssActiveDischargeEnergy", commonLabels, fp.prometheus.batteryPowerDischargeTotal)
	fp.queryCommon(&wg, client, "_sum/EssActivePowerL1", phase1Labels, fp.prometheus.batteryPowerPhase)
	fp.queryCommon(&wg, client, "_sum/EssActivePowerL2", phase2Labels, fp.prometheus.batteryPowerPhase)
	fp.queryCommon(&wg, client, "_sum/EssActivePowerL3", phase3Labels, fp.prometheus.batteryPowerPhase)

	// grid
	fp.queryCommon(&wg, client, "_sum/GridActivePower", commonLabels, fp.prometheus.gridPower)
	fp.queryCommonIfNotZero(&wg, client, "_sum/GridBuyActiveEnergy", commonLabels, fp.prometheus.gridPowerBuyTotal)
	fp.queryCommonIfNotZero(&wg, client, "_sum/GridSellActiveEnergy", commonLabels, fp.prometheus.gridPowerSellTotal)
	fp.queryCommon(&wg, client, "_sum/GridActivePowerL1", phase1Labels, fp.prometheus.gridPowerPhase)
	fp.queryCommon(&wg, client, "_sum/GridActivePowerL2", phase2Labels, fp.prometheus.gridPowerPhase)
	fp.queryCommon(&wg, client, "_sum/GridActivePowerL3", phase3Labels, fp.prometheus.gridPowerPhase)

	// production
	fp.queryCommon(&wg, client, "_sum/ProductionActivePower", commonLabels, fp.prometheus.productionPower)
	fp.queryCommon(&wg, client, "_sum/ProductionAcActivePower", commonLabels, fp.prometheus.productionPowerAc)
	fp.queryCommon(&wg, client, "_sum/ProductionDcActualPower", commonLabels, fp.prometheus.productionPowerDc)
	fp.queryCommonIfNotZero(&wg, client, "_sum/ProductionActiveEnergy", commonLabels, fp.prometheus.productionPowerTotal)
	fp.queryCommonIfNotZero(&wg, client, "_sum/ProductionAcActiveEnergy", commonLabels, fp.prometheus.productionPowerAcTotal)
	fp.queryCommonIfNotZero(&wg, client, "_sum/ProductionDcActiveEnergy", commonLabels, fp.prometheus.productionPowerDcTotal)
	fp.queryCommon(&wg, client, "_sum/ProductionAcActivePowerL1", phase1Labels, fp.prometheus.productionPowerPhase)
	fp.queryCommon(&wg, client, "_sum/ProductionAcActivePowerL2", phase2Labels, fp.prometheus.productionPowerPhase)
	fp.queryCommon(&wg, client, "_sum/ProductionAcActivePowerL3", phase3Labels, fp.prometheus.productionPowerPhase)

	// consumption
	fp.queryCommon(&wg, client, "_sum/ConsumptionActivePower", commonLabels, fp.prometheus.consumptionPower)
	fp.queryCommonIfNotZero(&wg, client, "_sum/ConsumptionActiveEnergy", commonLabels, fp.prometheus.consumptionPowerTotal)
	fp.queryCommon(&wg, client, "_sum/ConsumptionActivePowerL1", phase1Labels, fp.prometheus.consumptionPowerPhase)
	fp.queryCommon(&wg, client, "_sum/ConsumptionActivePowerL2", phase2Labels, fp.prometheus.consumptionPowerPhase)
	fp.queryCommon(&wg, client, "_sum/ConsumptionActivePowerL3", phase3Labels, fp.prometheus.consumptionPowerPhase)

	wg.Wait()

	fp.logger.Debugf(`finished probe in %v`, time.Since(startTime).String())
}

func (fp *FeneconProber) queryCommon(wg *sizedwaitgroup.SizedWaitGroup, client *resty.Client, url string, labels prometheus.Labels, gaugeVec *prometheus.GaugeVec) {
	fp.queryCommonCallback(
		wg,
		client,
		url,
		func(result ResultCommon) {
			if result.Value != nil {
				gaugeVec.With(labels).Set(*result.Value)
			}
		},
	)
}

func (fp *FeneconProber) queryCommonIfNotZero(wg *sizedwaitgroup.SizedWaitGroup, client *resty.Client, url string, labels prometheus.Labels, gaugeVec *prometheus.GaugeVec) {
	fp.queryCommonCallback(
		wg,
		client,
		url,
		func(result ResultCommon) {
			if result.Value != nil && *result.Value > 0 {
				gaugeVec.With(labels).Set(*result.Value)
			}
		},
	)
}

func (fp *FeneconProber) queryCommonCallback(wg *sizedwaitgroup.SizedWaitGroup, client *resty.Client, url string, callback func(result ResultCommon)) {
	wg.Add()
	go func() {
		defer wg.Done()
		startTime := time.Now()
		fp.logger.Debugf(`start query %v`, url)

		result := ResultCommon{}
		_, err := client.R().SetContext(fp.ctx).SetResult(&result).Get(url)
		if err == nil {
			fp.logger.Debugf(`finished query %v in %v`, url, time.Since(startTime).String())
			callback(result)
		} else {
			fp.logger.Errorf(`failed query %v in %v: %v`, url, time.Since(startTime).String(), err)
		}
	}()
}
