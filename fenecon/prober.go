package fenecon

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	resty "github.com/go-resty/resty/v2"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type (
	FeneconProber struct {
		ctx      context.Context
		client   *resty.Client
		logger   *zap.SugaredLogger
		registry *prometheus.Registry

		target string

		prometheus feneconMetrics
	}
)

func New(ctx context.Context, registry *prometheus.Registry, logger *zap.SugaredLogger) *FeneconProber {
	fp := FeneconProber{}
	fp.ctx = ctx
	fp.registry = registry
	fp.logger = logger
	fp.initResty()
	fp.initMetrics()

	return &fp
}

func (fp *FeneconProber) initResty() {
	fp.client = resty.New()

	fp.client.OnAfterResponse(func(client *resty.Client, response *resty.Response) error {
		switch statusCode := response.StatusCode(); statusCode {
		case 401:
			return errors.New(`fenecon requires authentication and/or credentials are invalid`)
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

func (fp *FeneconProber) SetTimeout(timeout time.Duration) {
	fp.client.SetTimeout(timeout)
}

func (fp *FeneconProber) SetHttpAuth(username, password string) {
	fp.client.SetDisableWarn(true)
	fp.client.SetBasicAuth(username, password)
}

func (fp *FeneconProber) Run(target string) {
	fp.target = target

	client := fp.client.SetBaseURL(
		fmt.Sprintf(`%s/rest/channel/`, strings.TrimRight(target, "/")),
	)

	commonLabels := prometheus.Labels{"target": target}
	phase1Labels := prometheus.Labels{"target": target, "phase": "1"}
	phase2Labels := prometheus.Labels{"target": target, "phase": "2"}
	phase3Labels := prometheus.Labels{"target": target, "phase": "3"}

	fp.prometheus.info.With(commonLabels).Set(1)

	wg := sync.WaitGroup{}

	// general
	wg.Add(1)
	go func() {
		defer wg.Done()
		fp.queryCommon(client, "_sum/State", commonLabels, fp.prometheus.status)
	}()

	// battery
	wg.Add(1)
	go func() {
		defer wg.Done()
		fp.queryCommon(client, "_sum/EssSoc", commonLabels, fp.prometheus.batteryCharge)
		fp.queryCommon(client, "_sum/EssActivePower", commonLabels, fp.prometheus.batteryPower)
		fp.queryCommon(client, "_sum/EssActiveChargeEnergy", commonLabels, fp.prometheus.batteryPowerChargeTotal)
		fp.queryCommon(client, "_sum/EssActiveDischargeEnergy", commonLabels, fp.prometheus.batteryPowerDischargeTotal)
		fp.queryCommon(client, "_sum/EssActivePowerL1", phase1Labels, fp.prometheus.batteryPowerPhase)
		fp.queryCommon(client, "_sum/EssActivePowerL2", phase2Labels, fp.prometheus.batteryPowerPhase)
		fp.queryCommon(client, "_sum/EssActivePowerL3", phase3Labels, fp.prometheus.batteryPowerPhase)
	}()

	// grid
	wg.Add(1)
	go func() {
		defer wg.Done()
		fp.queryCommon(client, "_sum/GridActivePower", commonLabels, fp.prometheus.gridPower)
		fp.queryCommon(client, "_sum/GridBuyActiveEnergy", commonLabels, fp.prometheus.gridPowerBuyTotal)
		fp.queryCommon(client, "_sum/GridSellActiveEnergy", commonLabels, fp.prometheus.gridPowerSellTotal)
		fp.queryCommon(client, "_sum/GridActivePowerL1", phase1Labels, fp.prometheus.gridPowerPhase)
		fp.queryCommon(client, "_sum/GridActivePowerL2", phase2Labels, fp.prometheus.gridPowerPhase)
		fp.queryCommon(client, "_sum/GridActivePowerL3", phase3Labels, fp.prometheus.gridPowerPhase)
	}()

	// production
	wg.Add(1)
	go func() {
		defer wg.Done()
		fp.queryCommon(client, "_sum/ProductionActivePower", commonLabels, fp.prometheus.productionPower)
		fp.queryCommon(client, "_sum/ProductionAcActivePower", commonLabels, fp.prometheus.productionPowerAc)
		fp.queryCommon(client, "_sum/ProductionDcActualPower", commonLabels, fp.prometheus.productionPowerDc)
		fp.queryCommon(client, "_sum/ProductionActiveEnergy", commonLabels, fp.prometheus.productionPowerTotal)
		fp.queryCommon(client, "_sum/ProductionAcActiveEnergy", commonLabels, fp.prometheus.productionPowerAcTotal)
		fp.queryCommon(client, "_sum/ProductionDcActiveEnergy", commonLabels, fp.prometheus.productionPowerDcTotal)
		fp.queryCommon(client, "_sum/ProductionAcActivePowerL1", phase1Labels, fp.prometheus.productionPowerPhase)
		fp.queryCommon(client, "_sum/ProductionAcActivePowerL2", phase2Labels, fp.prometheus.productionPowerPhase)
		fp.queryCommon(client, "_sum/ProductionAcActivePowerL3", phase3Labels, fp.prometheus.productionPowerPhase)
	}()

	// consumption
	wg.Add(1)
	go func() {
		defer wg.Done()
		fp.queryCommon(client, "_sum/ConsumptionActivePower", commonLabels, fp.prometheus.consumptionPower)
		fp.queryCommon(client, "_sum/ConsumptionActiveEnergy", commonLabels, fp.prometheus.consumptionPowerTotal)
		fp.queryCommon(client, "_sum/ConsumptionActivePowerL1", phase1Labels, fp.prometheus.consumptionPowerPhase)
		fp.queryCommon(client, "_sum/ConsumptionActivePowerL2", phase2Labels, fp.prometheus.consumptionPowerPhase)
		fp.queryCommon(client, "_sum/ConsumptionActivePowerL3", phase3Labels, fp.prometheus.consumptionPowerPhase)
	}()

	wg.Wait()
}

func (fp *FeneconProber) queryCommon(client *resty.Client, url string, labels prometheus.Labels, gaugeVec *prometheus.GaugeVec) {
	result := ResultCommon{}
	_, err := client.R().SetContext(fp.ctx).SetResult(&result).Get(url)
	if err != nil {
		fp.logger.Error(err)
	}

	gaugeVec.With(labels).Set(result.Value)
}
