package fenecon

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
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
		Target string
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

	// ------------------------------------------------------------------------
	// SUM
	result, err := fp.queryWildcard(client, "_sum/.*")
	if err == nil {
		// general
		result.Address("_sum/State").SetGauge(commonLabels, fp.prometheus.status)

		// battery
		result.Address("_sum/EssSoc").SetGauge(commonLabels, fp.prometheus.battery.charge)
		result.Address("_sum/EssCapacity").SetGauge(commonLabels, fp.prometheus.battery.capacity)
		result.Address("_sum/EssActivePower").SetGauge(commonLabels, fp.prometheus.battery.power)
		result.Address("_sum/EssActiveChargeEnergy").SetGaugeIfNotZero(commonLabels, fp.prometheus.battery.powerChargeTotal)
		result.Address("_sum/EssActiveDischargeEnergy").SetGaugeIfNotZero(commonLabels, fp.prometheus.battery.powerDischargeTotal)
		result.Address("_sum/EssDcChargeEnergy").SetGaugeIfNotZero(commonLabels, fp.prometheus.battery.powerDcChargeTotal)
		result.Address("_sum/EssDcDischargeEnergy").SetGaugeIfNotZero(commonLabels, fp.prometheus.battery.powerDcDischargeTotal)
		result.Address("_sum/EssActivePowerL1").SetGauge(phase1Labels, fp.prometheus.battery.powerPhase)
		result.Address("_sum/EssActivePowerL2").SetGauge(phase2Labels, fp.prometheus.battery.powerPhase)
		result.Address("_sum/EssActivePowerL3").SetGauge(phase3Labels, fp.prometheus.battery.powerPhase)

		// grid
		result.Address("_sum/GridMode").SetGauge(commonLabels, fp.prometheus.grid.mode)
		result.Address("_sum/GridActivePower").SetGauge(commonLabels, fp.prometheus.grid.power)
		result.Address("_sum/GridBuyActiveEnergy").SetGaugeIfNotZero(commonLabels, fp.prometheus.grid.powerBuyTotal)
		result.Address("_sum/GridSellActiveEnergy").SetGaugeIfNotZero(commonLabels, fp.prometheus.grid.powerSellTotal)
		result.Address("_sum/GridActivePowerL1").SetGauge(phase1Labels, fp.prometheus.grid.powerPhase)
		result.Address("_sum/GridActivePowerL2").SetGauge(phase2Labels, fp.prometheus.grid.powerPhase)
		result.Address("_sum/GridActivePowerL3").SetGauge(phase3Labels, fp.prometheus.grid.powerPhase)

		// grid
		result.Address("_sum/GridMode").SetGauge(commonLabels, fp.prometheus.grid.mode)
		result.Address("_sum/GridActivePower").SetGauge(commonLabels, fp.prometheus.grid.power)
		result.Address("_sum/GridBuyActiveEnergy").SetGaugeIfNotZero(commonLabels, fp.prometheus.grid.powerBuyTotal)
		result.Address("_sum/GridSellActiveEnergy").SetGaugeIfNotZero(commonLabels, fp.prometheus.grid.powerSellTotal)
		result.Address("_sum/GridActivePowerL1").SetGauge(phase1Labels, fp.prometheus.grid.powerPhase)
		result.Address("_sum/GridActivePowerL2").SetGauge(phase2Labels, fp.prometheus.grid.powerPhase)
		result.Address("_sum/GridActivePowerL3").SetGauge(phase3Labels, fp.prometheus.grid.powerPhase)

		// production
		result.Address("_sum/ProductionActivePower").SetGauge(commonLabels, fp.prometheus.production.power)
		result.Address("_sum/ProductionAcActivePower").SetGauge(commonLabels, fp.prometheus.production.powerAc)
		result.Address("_sum/ProductionDcActualPower").SetGauge(commonLabels, fp.prometheus.production.powerDc)
		result.Address("_sum/ProductionActiveEnergy").SetGaugeIfNotZero(commonLabels, fp.prometheus.production.powerTotal)
		result.Address("_sum/ProductionAcActiveEnergy").SetGaugeIfNotZero(commonLabels, fp.prometheus.production.powerAcTotal)
		result.Address("_sum/ProductionDcActiveEnergy").SetGaugeIfNotZero(commonLabels, fp.prometheus.production.powerDcTotal)
		result.Address("_sum/ProductionAcActivePowerL1").SetGauge(phase1Labels, fp.prometheus.production.powerPhase)
		result.Address("_sum/ProductionAcActivePowerL2").SetGauge(phase2Labels, fp.prometheus.production.powerPhase)
		result.Address("_sum/ProductionAcActivePowerL3").SetGauge(phase3Labels, fp.prometheus.production.powerPhase)

		// consumption
		result.Address("_sum/ConsumptionActivePower").SetGauge(commonLabels, fp.prometheus.consumption.power)
		result.Address("_sum/ConsumptionActiveEnergy").SetGaugeIfNotZero(commonLabels, fp.prometheus.consumption.powerTotal)
		result.Address("_sum/ConsumptionActivePowerL1").SetGauge(phase1Labels, fp.prometheus.consumption.powerPhase)
		result.Address("_sum/ConsumptionActivePowerL2").SetGauge(phase2Labels, fp.prometheus.consumption.powerPhase)
		result.Address("_sum/ConsumptionActivePowerL3").SetGauge(phase3Labels, fp.prometheus.consumption.powerPhase)
	}

	// ------------------------------------------------------------------------
	// Ess (Batteries)
	result, err = fp.queryWildcard(client, "ess.*/.*")
	if err == nil {
		for _, module := range result.AddressParts(0) {
			batteryLabels := prometheus.Labels{"target": target.Target, "module": module}

			result.Address(module, "State").SetGauge(batteryLabels, fp.prometheus.status)
			result.Address(module, "GridMode").SetGauge(batteryLabels, fp.prometheus.grid.mode)
			result.Address(module, "Soc").SetGauge(batteryLabels, fp.prometheus.battery.charge)
			result.Address(module, "Capacity").SetGauge(batteryLabels, fp.prometheus.battery.capacity)
			result.Address(module, "ActivePower").SetGauge(batteryLabels, fp.prometheus.battery.power)
			result.Address(module, "ActiveChargeEnergy").SetGaugeIfNotZero(batteryLabels, fp.prometheus.battery.powerChargeTotal)
			result.Address(module, "ActiveDischargeEnergy").SetGaugeIfNotZero(batteryLabels, fp.prometheus.battery.powerDischargeTotal)
			result.Address(module, "AllowedChargePower").SetGauge(batteryLabels, fp.prometheus.battery.allowedChargePower)
			result.Address(module, "AllowedDischargePower").SetGauge(batteryLabels, fp.prometheus.battery.allowedDischargePower)
		}
	}

	// ------------------------------------------------------------------------
	// Charger (eg. panels)
	result, err = fp.queryWildcard(client, "charger.*/.*")
	if err == nil {
		for _, module := range result.AddressParts(0) {
			chargerLabels := prometheus.Labels{"target": target.Target, "module": module}

			result.Address(module, "State").SetGauge(chargerLabels, fp.prometheus.status)
			result.Address(module, "ActualPower").SetGauge(chargerLabels, fp.prometheus.production.power)
			result.Address(module, "ActualEnergy").SetGaugeIfNotZero(chargerLabels, fp.prometheus.production.powerTotal)
			result.Address(module, "MaxActualPower").SetGauge(chargerLabels, fp.prometheus.production.maxActualPower)
		}
	}

	// ------------------------------------------------------------------------
	// Meter (eg. panels)
	result, err = fp.queryWildcard(client, "meter.*/.*")
	if err == nil {
		for _, module := range result.AddressParts(0) {
			meterLabels := prometheus.Labels{"target": target.Target, "module": module}
			meterPhase1Labels := prometheus.Labels{"target": target.Target, "module": module, "phase": "1"}
			meterPhase2Labels := prometheus.Labels{"target": target.Target, "module": module, "phase": "2"}
			meterPhase3Labels := prometheus.Labels{"target": target.Target, "module": module, "phase": "3"}

			result.Address(module, "State").SetGauge(meterLabels, fp.prometheus.status)
			result.Address(module, "Frequency").SetGauge(meterLabels, fp.prometheus.meter.frequency)
			result.Address(module, "Voltage").SetGauge(meterLabels, fp.prometheus.meter.voltage)
			result.Address(module, "VoltageL1").SetGauge(meterPhase1Labels, fp.prometheus.meter.voltagePhase)
			result.Address(module, "VoltageL2").SetGauge(meterPhase2Labels, fp.prometheus.meter.voltagePhase)
			result.Address(module, "VoltageL3").SetGauge(meterPhase3Labels, fp.prometheus.meter.voltagePhase)
			result.Address(module, "ActivePower").SetGauge(meterLabels, fp.prometheus.meter.power)
			result.Address(module, "ActivePowerL1").SetGauge(meterPhase1Labels, fp.prometheus.meter.powerPhase)
			result.Address(module, "ActivePowerL2").SetGauge(meterPhase2Labels, fp.prometheus.meter.powerPhase)
			result.Address(module, "ActivePowerL3").SetGauge(meterPhase3Labels, fp.prometheus.meter.powerPhase)
			result.Address(module, "ReactivePower").SetGauge(meterLabels, fp.prometheus.meter.reactivePower)
			result.Address(module, "ReactivePowerL1").SetGauge(meterPhase1Labels, fp.prometheus.meter.reactivePowerPhase)
			result.Address(module, "ReactivePowerL2").SetGauge(meterPhase2Labels, fp.prometheus.meter.reactivePowerPhase)
			result.Address(module, "ReactivePowerL3").SetGauge(meterPhase3Labels, fp.prometheus.meter.reactivePowerPhase)
			result.Address(module, "Current").SetGauge(meterLabels, fp.prometheus.meter.current)
			result.Address(module, "CurrentL1").SetGauge(meterPhase1Labels, fp.prometheus.meter.currentPhase)
			result.Address(module, "CurrentL2").SetGauge(meterPhase2Labels, fp.prometheus.meter.currentPhase)
			result.Address(module, "CurrentL3").SetGauge(meterPhase3Labels, fp.prometheus.meter.currentPhase)
			result.Address(module, "MinActivePower").SetGauge(meterLabels, fp.prometheus.meter.minActivePower)
			result.Address(module, "MaxActivePower").SetGauge(meterLabels, fp.prometheus.meter.maxActivePower)
			result.Address(module, "ActiveProductionEnergy").SetGauge(meterLabels, fp.prometheus.meter.powerProductionTotal)
			result.Address(module, "ActiveConsumptionEnergy").SetGauge(meterLabels, fp.prometheus.meter.powerConsumptionTotal)
		}
	}

	fp.logger.Debugf(`finished probe in %v`, time.Since(startTime).String())
}

func (fp *FeneconProber) queryWildcard(client *resty.Client, url string) (*ResultWildcard, error) {
	result := ResultWildcard{}

	startTime := time.Now()
	fp.logger.Debugf(`start query %v`, url)

	_, err := client.R().SetContext(fp.ctx).SetResult(&result).Get(url)
	if err == nil {
		fp.logger.Debugf(`finished query %v in %v`, url, time.Since(startTime).String())
	} else {
		fp.logger.Errorf(`failed query %v in %v: %v`, url, time.Since(startTime).String(), err)
	}

	return &result, err
}
