package fenecon

import (
	"github.com/prometheus/client_golang/prometheus"
)

type (
	feneconMetrics struct {
		info   *prometheus.GaugeVec
		status *prometheus.GaugeVec

		meter struct {
			frequency             *prometheus.GaugeVec
			voltage               *prometheus.GaugeVec
			voltagePhase          *prometheus.GaugeVec
			power                 *prometheus.GaugeVec
			powerPhase            *prometheus.GaugeVec
			reactivePower         *prometheus.GaugeVec
			reactivePowerPhase    *prometheus.GaugeVec
			current               *prometheus.GaugeVec
			currentPhase          *prometheus.GaugeVec
			minActivePower        *prometheus.GaugeVec
			maxActivePower        *prometheus.GaugeVec
			powerProductionTotal  *prometheus.GaugeVec
			powerConsumptionTotal *prometheus.GaugeVec
		}

		battery struct {
			charge                *prometheus.GaugeVec
			capacity              *prometheus.GaugeVec
			power                 *prometheus.GaugeVec
			powerPhase            *prometheus.GaugeVec
			powerChargeTotal      *prometheus.GaugeVec
			powerDischargeTotal   *prometheus.GaugeVec
			powerDcChargeTotal    *prometheus.GaugeVec
			powerDcDischargeTotal *prometheus.GaugeVec
			allowedChargePower    *prometheus.GaugeVec
			allowedDischargePower *prometheus.GaugeVec
		}

		grid struct {
			mode           *prometheus.GaugeVec
			power          *prometheus.GaugeVec
			powerPhase     *prometheus.GaugeVec
			powerBuyTotal  *prometheus.GaugeVec
			powerSellTotal *prometheus.GaugeVec
		}

		production struct {
			power          *prometheus.GaugeVec
			powerPhase     *prometheus.GaugeVec
			powerAc        *prometheus.GaugeVec
			powerDc        *prometheus.GaugeVec
			powerTotal     *prometheus.GaugeVec
			powerAcTotal   *prometheus.GaugeVec
			powerDcTotal   *prometheus.GaugeVec
			maxActualPower *prometheus.GaugeVec
		}

		consumption struct {
			power      *prometheus.GaugeVec
			powerPhase *prometheus.GaugeVec
			powerTotal *prometheus.GaugeVec
		}
	}
)

func (fp *FeneconProber) initMetrics() {
	commonLabels := []string{"target", "module"}
	phaseLabels := []string{"target", "module", "phase"}

	// ##########################################
	// Info

	fp.prometheus.info = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_info",
			Help: "Fenecon info",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.info)

	// ##########################################
	// Status

	fp.prometheus.status = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_status",
			Help: "Fenecon status (0=ok, 1=info, 2=warning, 3=error; State)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.status)

	// ##########################################
	// Meter

	fp.prometheus.meter.frequency = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_meter_frequency",
			Help: "Fenecon meter frequenc in Hz (Frequency)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.meter.frequency)

	fp.prometheus.meter.voltage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_meter_voltage",
			Help: "Fenecon meter voltage in Volt (Voltage)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.meter.voltage)

	fp.prometheus.meter.voltagePhase = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_meter_voltage_phase",
			Help: "Fenecon meter voltage in Volt (VoltagePhase)",
		},
		phaseLabels,
	)
	fp.registry.MustRegister(fp.prometheus.meter.voltagePhase)

	fp.prometheus.meter.power = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_meter_power",
			Help: "Fenecon meter power in Watts (ActivePower)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.meter.power)

	fp.prometheus.meter.powerPhase = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_meter_power_phase",
			Help: "Fenecon meter power in Watts (ActivePowerLx)",
		},
		phaseLabels,
	)
	fp.registry.MustRegister(fp.prometheus.meter.powerPhase)

	fp.prometheus.meter.reactivePower = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_meter_reactive_power",
			Help: "Fenecon meter reactive  power in Watts (ReactivePower)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.meter.reactivePower)

	fp.prometheus.meter.reactivePowerPhase = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_meter_reactive_power_phase",
			Help: "Fenecon meter reactive power in Watts (ReactivePowerLx)",
		},
		phaseLabels,
	)
	fp.registry.MustRegister(fp.prometheus.meter.reactivePowerPhase)

	fp.prometheus.meter.current = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_meter_current",
			Help: "Fenecon meter current in mA (Current)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.meter.current)

	fp.prometheus.meter.currentPhase = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_meter_current_phase",
			Help: "Fenecon meter current in mA (CurrentL1x)",
		},
		phaseLabels,
	)
	fp.registry.MustRegister(fp.prometheus.meter.currentPhase)

	fp.prometheus.meter.minActivePower = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_meter_min_active_power",
			Help: "Fenecon meter min active power in Watts (MinActivePower)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.meter.minActivePower)

	fp.prometheus.meter.maxActivePower = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_meter_max_active_power",
			Help: "Fenecon meter max active power in Watts (MaxActivePower)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.meter.maxActivePower)

	fp.prometheus.meter.powerProductionTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_meter_power_production_total",
			Help: "Fenecon meter power production total in Watthours (ActiveProductionEnergy)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.meter.powerProductionTotal)

	fp.prometheus.meter.powerConsumptionTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_meter_power_consumption_total",
			Help: "Fenecon meter power consumption total in Watthours (ActiveConsumptionEnergy)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.meter.powerConsumptionTotal)

	// ##########################################
	// Battery

	fp.prometheus.battery.charge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_battery_charge_percent",
			Help: "Fenecon battery charge in percent (EssSoc)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.battery.charge)

	fp.prometheus.battery.capacity = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_battery_capacity",
			Help: "Fenecon battery capacity in Watthours (EssCapacity)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.battery.capacity)

	fp.prometheus.battery.power = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_battery_power",
			Help: "Fenecon battery power load in Watts (EssActivePower)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.battery.power)

	fp.prometheus.battery.powerPhase = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_battery_power_phase",
			Help: "Fenecon battery power load in Watts (EssActivePowerLx)",
		},
		phaseLabels,
	)
	fp.registry.MustRegister(fp.prometheus.battery.powerPhase)

	fp.prometheus.battery.powerChargeTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_battery_power_charge_total",
			Help: "Fenecon battery power charge in Wattshours (EssActiveChargeEnergy)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.battery.powerChargeTotal)

	fp.prometheus.battery.powerDischargeTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_battery_power_discharge_total",
			Help: "Fenecon battery power discharge in Wattshours (EssActiveDischargeEnergy)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.battery.powerDischargeTotal)

	fp.prometheus.battery.powerDcChargeTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_battery_power_dc_charge_total",
			Help: "Fenecon battery power dc charge in Wattshours (EssDcChargeEnergy)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.battery.powerDcChargeTotal)

	fp.prometheus.battery.powerDcDischargeTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_battery_power_dc_discharge_total",
			Help: "Fenecon battery power dc discharge in Wattshours (EssDcDischargeEnergy)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.battery.powerDcDischargeTotal)

	fp.prometheus.battery.allowedChargePower = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_battery_allowed_charge_power",
			Help: "Fenecon battery allowed scharge power Watts (AllowedChargePower)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.battery.allowedChargePower)

	fp.prometheus.battery.allowedDischargePower = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_battery_allowed_discharge_power",
			Help: "Fenecon battery allowed discharge power Watts (AllowedDischargePower)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.battery.allowedDischargePower)

	// ##########################################
	// Grid

	fp.prometheus.grid.mode = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_grid_mode",
			Help: "Fenecon grid mode (0=undefined, 1=On-Grid, 2=Off-Grid; GridActivePower)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.grid.mode)

	fp.prometheus.grid.power = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_grid_power",
			Help: "Fenecon grid power load in Watts (GridActivePower)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.grid.power)

	fp.prometheus.grid.powerPhase = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_grid_power_phase",
			Help: "Fenecon grid power load in Watts (GridActivePowerLx)",
		},
		phaseLabels,
	)
	fp.registry.MustRegister(fp.prometheus.grid.powerPhase)

	fp.prometheus.grid.powerBuyTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_grid_power_buy_total",
			Help: "Fenecon grid power buy in Wattshours (GridBuyActiveEnergy)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.grid.powerBuyTotal)

	fp.prometheus.grid.powerSellTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_grid_power_sell_total",
			Help: "Fenecon grid power sell in Wattshours (GridSellActiveEnergy)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.grid.powerSellTotal)

	// ##########################################
	// Production

	fp.prometheus.production.power = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_production_power",
			Help: "Fenecon production power load in Watts (ProductionActivePower)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.production.power)

	fp.prometheus.production.powerPhase = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_production_power_phase",
			Help: "Fenecon production power load in Watts (ProductionAcActivePowerLx)",
		},
		phaseLabels,
	)
	fp.registry.MustRegister(fp.prometheus.production.powerPhase)

	fp.prometheus.production.powerAc = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_production_power_ac",
			Help: "Fenecon production power load in Watts (ProductionAcActivePower)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.production.powerAc)

	fp.prometheus.production.powerDc = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_production_power_dc",
			Help: "Fenecon production power load in Watts (ProductionDcActualPower)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.production.powerDc)

	fp.prometheus.production.powerTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_production_power_total",
			Help: "Fenecon production power load in Watthours (ProductionActiveEnergy)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.production.powerTotal)

	fp.prometheus.production.powerAcTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_production_power_ac_total",
			Help: "Fenecon production power load in Watthours (ProductionAcActiveEnergy)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.production.powerAcTotal)

	fp.prometheus.production.powerDcTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_production_power_dc_total",
			Help: "Fenecon production power load in Watthours (ProductionDcActiveEnergy)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.production.powerDcTotal)

	fp.prometheus.production.maxActualPower = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_production_max_actual_power",
			Help: "Fenecon production max acutal power Watts (MaxActualPower)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.production.maxActualPower)

	// ##########################################
	// Consumer

	fp.prometheus.consumption.power = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_consumption_power",
			Help: "Fenecon consumption power load in Watts (ConsumptionActivePower)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.consumption.power)

	fp.prometheus.consumption.powerPhase = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_consumption_power_phase",
			Help: "Fenecon consumption power load in Watts (ConsumptionActivePowerLX)",
		},
		phaseLabels,
	)
	fp.registry.MustRegister(fp.prometheus.consumption.powerPhase)

	fp.prometheus.consumption.powerTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_consumption_power_total",
			Help: "Fenecon consumption power load in Watts (ConsumptionActiveEnergy)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.consumption.powerTotal)

}
