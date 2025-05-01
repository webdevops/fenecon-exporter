package fenecon

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/remeh/sizedwaitgroup"
	"go.uber.org/zap"
	resty "resty.dev/v3"
)

type (
	FeneconProber struct {
		ctx      context.Context
		client   *resty.Client
		logger   *zap.SugaredLogger
		registry *prometheus.Registry

		parallelRequests int

		target FeneconProberTarget

		prometheus feneconMetrics
	}

	FeneconProberTarget struct {
		Target  string
		Meter   int
		Ess     int
		Charger int
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
	fp.client.SetRetryCount(3)
	fp.client.SetRetryWaitTime(2 * time.Second)
	fp.client.SetRetryMaxWaitTime(5 * time.Second)
	fp.client.EnableRetryDefaultConditions()
	fp.client.AddResponseMiddleware(func(client *resty.Client, response *resty.Response) error {
		switch statusCode := response.StatusCode(); statusCode {
		case 401:
			return errors.New(`fenecon requires authentication or credentials are invalid`)
		case 404:
			// ignore non existing endpoints
			fp.logger.Debugf(`got status 404 for url "%v"`, response.Request.URL)
			return nil
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

func (fp *FeneconProber) SetRetry(retry int, waitTime, maxWaitTime time.Duration) {
	fp.client.SetRetryCount(retry)
	fp.client.SetRetryWaitTime(waitTime)
	fp.client.SetRetryMaxWaitTime(maxWaitTime)
}

func (fp *FeneconProber) SetTimeout(timeout time.Duration) {
	fp.client.SetTimeout(timeout)
}

func (fp *FeneconProber) SetHttpAuth(username, password string) {
	fp.client.SetDisableWarn(true)
	fp.client.SetBasicAuth(username, password)
}

func (fp *FeneconProber) Run(target FeneconProberTarget) {
	fp.target = target
	fp.logger.With(zap.String("target", target.Target))

	startTime := time.Now()
	fp.logger.Debugf(`start probe`)

	client := fp.client.SetBaseURL(
		fmt.Sprintf(`%s/rest/channel/`, strings.TrimRight(target.Target, "/")),
	)

	commonLabels := prometheus.Labels{"target": target.Target, "module": "_sum"}
	phase1Labels := prometheus.Labels{"target": target.Target, "module": "_sum", "phase": "1"}
	phase2Labels := prometheus.Labels{"target": target.Target, "module": "_sum", "phase": "2"}
	phase3Labels := prometheus.Labels{"target": target.Target, "module": "_sum", "phase": "3"}

	fp.prometheus.info.With(commonLabels).Set(1)

	wg := sizedwaitgroup.New(fp.parallelRequests)

	// ------------------------------------------------------------------------
	// SUM

	// general
	fp.queryCommon(&wg, client, "_sum/State", commonLabels, fp.prometheus.status)

	// battery
	fp.queryCommon(&wg, client, "_sum/EssSoc", commonLabels, fp.prometheus.battery.charge)
	fp.queryCommon(&wg, client, "_sum/EssCapacity", commonLabels, fp.prometheus.battery.capacity)
	fp.queryCommon(&wg, client, "_sum/EssActivePower", commonLabels, fp.prometheus.battery.power)
	fp.queryCommonIfNotZero(&wg, client, "_sum/EssActiveChargeEnergy", commonLabels, fp.prometheus.battery.powerChargeTotal)
	fp.queryCommonIfNotZero(&wg, client, "_sum/EssActiveDischargeEnergy", commonLabels, fp.prometheus.battery.powerDischargeTotal)
	fp.queryCommonIfNotZero(&wg, client, "_sum/EssDcChargeEnergy", commonLabels, fp.prometheus.battery.powerDcChargeTotal)
	fp.queryCommonIfNotZero(&wg, client, "_sum/EssDcDischargeEnergy", commonLabels, fp.prometheus.battery.powerDcDischargeTotal)
	fp.queryCommon(&wg, client, "_sum/EssActivePowerL1", phase1Labels, fp.prometheus.battery.powerPhase)
	fp.queryCommon(&wg, client, "_sum/EssActivePowerL2", phase2Labels, fp.prometheus.battery.powerPhase)
	fp.queryCommon(&wg, client, "_sum/EssActivePowerL3", phase3Labels, fp.prometheus.battery.powerPhase)

	// grid
	fp.queryCommon(&wg, client, "_sum/GridMode", commonLabels, fp.prometheus.grid.mode)
	fp.queryCommon(&wg, client, "_sum/GridActivePower", commonLabels, fp.prometheus.grid.power)
	fp.queryCommonIfNotZero(&wg, client, "_sum/GridBuyActiveEnergy", commonLabels, fp.prometheus.grid.powerBuyTotal)
	fp.queryCommonIfNotZero(&wg, client, "_sum/GridSellActiveEnergy", commonLabels, fp.prometheus.grid.powerSellTotal)
	fp.queryCommon(&wg, client, "_sum/GridActivePowerL1", phase1Labels, fp.prometheus.grid.powerPhase)
	fp.queryCommon(&wg, client, "_sum/GridActivePowerL2", phase2Labels, fp.prometheus.grid.powerPhase)
	fp.queryCommon(&wg, client, "_sum/GridActivePowerL3", phase3Labels, fp.prometheus.grid.powerPhase)

	// production
	fp.queryCommon(&wg, client, "_sum/ProductionActivePower", commonLabels, fp.prometheus.production.power)
	fp.queryCommon(&wg, client, "_sum/ProductionAcActivePower", commonLabels, fp.prometheus.production.powerAc)
	fp.queryCommon(&wg, client, "_sum/ProductionDcActualPower", commonLabels, fp.prometheus.production.powerDc)
	fp.queryCommonIfNotZero(&wg, client, "_sum/ProductionActiveEnergy", commonLabels, fp.prometheus.production.powerTotal)
	fp.queryCommonIfNotZero(&wg, client, "_sum/ProductionAcActiveEnergy", commonLabels, fp.prometheus.production.powerAcTotal)
	fp.queryCommonIfNotZero(&wg, client, "_sum/ProductionDcActiveEnergy", commonLabels, fp.prometheus.production.powerDcTotal)
	fp.queryCommon(&wg, client, "_sum/ProductionAcActivePowerL1", phase1Labels, fp.prometheus.production.powerPhase)
	fp.queryCommon(&wg, client, "_sum/ProductionAcActivePowerL2", phase2Labels, fp.prometheus.production.powerPhase)
	fp.queryCommon(&wg, client, "_sum/ProductionAcActivePowerL3", phase3Labels, fp.prometheus.production.powerPhase)

	// consumption
	fp.queryCommon(&wg, client, "_sum/ConsumptionActivePower", commonLabels, fp.prometheus.consumption.power)
	fp.queryCommonIfNotZero(&wg, client, "_sum/ConsumptionActiveEnergy", commonLabels, fp.prometheus.consumption.powerTotal)
	fp.queryCommon(&wg, client, "_sum/ConsumptionActivePowerL1", phase1Labels, fp.prometheus.consumption.powerPhase)
	fp.queryCommon(&wg, client, "_sum/ConsumptionActivePowerL2", phase2Labels, fp.prometheus.consumption.powerPhase)
	fp.queryCommon(&wg, client, "_sum/ConsumptionActivePowerL3", phase3Labels, fp.prometheus.consumption.powerPhase)

	// ------------------------------------------------------------------------
	// Ess (Batteries)
	for i := 0; i < target.Ess; i++ {
		module := fmt.Sprintf("ess%v", i)
		batteryLabels := prometheus.Labels{"target": target.Target, "module": module}

		// battery
		fp.queryCommon(&wg, client, fmt.Sprintf("%s/State", module), batteryLabels, fp.prometheus.status)
		fp.queryCommon(&wg, client, fmt.Sprintf("%s/GridMode", module), batteryLabels, fp.prometheus.grid.mode)
		fp.queryCommon(&wg, client, fmt.Sprintf("%s/Soc", module), batteryLabels, fp.prometheus.battery.charge)
		fp.queryCommon(&wg, client, fmt.Sprintf("%s/Capacity", module), batteryLabels, fp.prometheus.battery.capacity)
		fp.queryCommon(&wg, client, fmt.Sprintf("%s/ActivePower", module), batteryLabels, fp.prometheus.battery.power)
		fp.queryCommonIfNotZero(&wg, client, fmt.Sprintf("%s/ActiveChargeEnergy", module), batteryLabels, fp.prometheus.battery.powerChargeTotal)
		fp.queryCommonIfNotZero(&wg, client, fmt.Sprintf("%s/ActiveDischargeEnergy", module), batteryLabels, fp.prometheus.battery.powerDischargeTotal)
		fp.queryCommon(&wg, client, fmt.Sprintf("%s/AllowedChargePower", module), batteryLabels, fp.prometheus.battery.allowedChargePower)
		fp.queryCommon(&wg, client, fmt.Sprintf("%s/AllowedDischargePower", module), batteryLabels, fp.prometheus.battery.allowedDischargePower)
	}

	// ------------------------------------------------------------------------
	// Charger (eg. panels)
	for i := 0; i < target.Charger; i++ {
		module := fmt.Sprintf("charger%v", i)
		chargerLabels := prometheus.Labels{"target": target.Target, "module": module}

		fp.queryCommon(&wg, client, fmt.Sprintf("%s/State", module), chargerLabels, fp.prometheus.status)
		fp.queryCommon(&wg, client, fmt.Sprintf("%s/ActualPower", module), chargerLabels, fp.prometheus.production.power)
		fp.queryCommonIfNotZero(&wg, client, fmt.Sprintf("%s/ActualEnergy", module), chargerLabels, fp.prometheus.production.powerTotal)
		fp.queryCommon(&wg, client, fmt.Sprintf("%s/MaxActualPower", module), chargerLabels, fp.prometheus.production.maxActualPower)
	}

	// ------------------------------------------------------------------------
	// Meter (eg. panels)
	for i := 0; i < target.Meter; i++ {
		module := fmt.Sprintf("meter%v", i)
		meterLabels := prometheus.Labels{"target": target.Target, "module": module}
		meterPhase1Labels := prometheus.Labels{"target": target.Target, "module": module, "phase": "1"}
		meterPhase2Labels := prometheus.Labels{"target": target.Target, "module": module, "phase": "2"}
		meterPhase3Labels := prometheus.Labels{"target": target.Target, "module": module, "phase": "3"}

		fp.queryCommon(&wg, client, fmt.Sprintf("%s/State", module), meterLabels, fp.prometheus.status)
		fp.queryCommon(&wg, client, fmt.Sprintf("%s/Frequency", module), meterLabels, fp.prometheus.meter.frequency)
		fp.queryCommon(&wg, client, fmt.Sprintf("%s/Voltage", module), meterLabels, fp.prometheus.meter.voltage)
		fp.queryCommon(&wg, client, fmt.Sprintf("%s/VoltageL1", module), meterPhase1Labels, fp.prometheus.meter.voltagePhase)
		fp.queryCommon(&wg, client, fmt.Sprintf("%s/VoltageL2", module), meterPhase2Labels, fp.prometheus.meter.voltagePhase)
		fp.queryCommon(&wg, client, fmt.Sprintf("%s/VoltageL3", module), meterPhase3Labels, fp.prometheus.meter.voltagePhase)
		fp.queryCommon(&wg, client, fmt.Sprintf("%s/ActivePower", module), meterLabels, fp.prometheus.meter.power)
		fp.queryCommon(&wg, client, fmt.Sprintf("%s/ActivePowerL1", module), meterPhase1Labels, fp.prometheus.meter.powerPhase)
		fp.queryCommon(&wg, client, fmt.Sprintf("%s/ActivePowerL2", module), meterPhase2Labels, fp.prometheus.meter.powerPhase)
		fp.queryCommon(&wg, client, fmt.Sprintf("%s/ActivePowerL3", module), meterPhase3Labels, fp.prometheus.meter.powerPhase)
		fp.queryCommon(&wg, client, fmt.Sprintf("%s/ReactivePower", module), meterLabels, fp.prometheus.meter.reactivePower)
		fp.queryCommon(&wg, client, fmt.Sprintf("%s/ReactivePowerL1", module), meterPhase1Labels, fp.prometheus.meter.reactivePowerPhase)
		fp.queryCommon(&wg, client, fmt.Sprintf("%s/ReactivePowerL2", module), meterPhase2Labels, fp.prometheus.meter.reactivePowerPhase)
		fp.queryCommon(&wg, client, fmt.Sprintf("%s/ReactivePowerL3", module), meterPhase3Labels, fp.prometheus.meter.reactivePowerPhase)
		fp.queryCommon(&wg, client, fmt.Sprintf("%s/Current", module), meterLabels, fp.prometheus.meter.current)
		fp.queryCommon(&wg, client, fmt.Sprintf("%s/CurrentL1", module), meterPhase1Labels, fp.prometheus.meter.currentPhase)
		fp.queryCommon(&wg, client, fmt.Sprintf("%s/CurrentL2", module), meterPhase2Labels, fp.prometheus.meter.currentPhase)
		fp.queryCommon(&wg, client, fmt.Sprintf("%s/CurrentL3", module), meterPhase3Labels, fp.prometheus.meter.currentPhase)
		fp.queryCommon(&wg, client, fmt.Sprintf("%s/MinActivePower", module), meterLabels, fp.prometheus.meter.minActivePower)
		fp.queryCommon(&wg, client, fmt.Sprintf("%s/MaxActivePower", module), meterLabels, fp.prometheus.meter.maxActivePower)
		fp.queryCommon(&wg, client, fmt.Sprintf("%s/ActiveProductionEnergy", module), meterLabels, fp.prometheus.meter.powerProductionTotal)
		fp.queryCommon(&wg, client, fmt.Sprintf("%s/ActiveConsumptionEnergy", module), meterLabels, fp.prometheus.meter.powerConsumptionTotal)
	}

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
