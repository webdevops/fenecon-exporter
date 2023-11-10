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

	fp.newGaugeVec(&fp.prometheus.info, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_info",
			Help: "Fenecon info",
		},
		commonLabels,
	))

	// ##########################################
	// Status

	fp.newGaugeVec(&fp.prometheus.status, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_status",
			Help: "Fenecon status (0=ok, 1=info, 2=warning, 3=error; State)",
		},
		commonLabels,
	))

	// ##########################################
	// Meter

	fp.newGaugeVec(&fp.prometheus.meter.frequency, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_meter_frequency",
			Help: "Fenecon meter frequenc in Hz (Frequency)",
		},
		commonLabels,
	))

	fp.newGaugeVec(&fp.prometheus.meter.voltage, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_meter_voltage",
			Help: "Fenecon meter voltage in Volt (Voltage)",
		},
		commonLabels,
	))

	fp.newGaugeVec(&fp.prometheus.meter.voltagePhase, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_meter_voltage_phase",
			Help: "Fenecon meter voltage in Volt (VoltagePhase)",
		},
		phaseLabels,
	))

	fp.newGaugeVec(&fp.prometheus.meter.power, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_meter_power",
			Help: "Fenecon meter power in Watts (ActivePower)",
		},
		commonLabels,
	))

	fp.newGaugeVec(&fp.prometheus.meter.powerPhase, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_meter_power_phase",
			Help: "Fenecon meter power in Watts (ActivePowerLx)",
		},
		phaseLabels,
	))

	fp.newGaugeVec(&fp.prometheus.meter.reactivePower, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_meter_reactive_power",
			Help: "Fenecon meter reactive  power in Watts (ReactivePower)",
		},
		commonLabels,
	))

	fp.newGaugeVec(&fp.prometheus.meter.reactivePowerPhase, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_meter_reactive_power_phase",
			Help: "Fenecon meter reactive power in Watts (ReactivePowerLx)",
		},
		phaseLabels,
	))

	fp.newGaugeVec(&fp.prometheus.meter.current, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_meter_current",
			Help: "Fenecon meter current in mA (Current)",
		},
		commonLabels,
	))

	fp.newGaugeVec(&fp.prometheus.meter.currentPhase, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_meter_current_phase",
			Help: "Fenecon meter current in mA (CurrentL1x)",
		},
		phaseLabels,
	))

	fp.newGaugeVec(&fp.prometheus.meter.minActivePower, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_meter_min_active_power",
			Help: "Fenecon meter min active power in Watts (MinActivePower)",
		},
		commonLabels,
	))

	fp.newGaugeVec(&fp.prometheus.meter.maxActivePower, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_meter_max_active_power",
			Help: "Fenecon meter max active power in Watts (MaxActivePower)",
		},
		commonLabels,
	))

	fp.newGaugeVec(&fp.prometheus.meter.powerProductionTotal, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_meter_power_production_total",
			Help: "Fenecon meter power production total in Watthours (ActiveProductionEnergy)",
		},
		commonLabels,
	))

	fp.newGaugeVec(&fp.prometheus.meter.powerConsumptionTotal, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_meter_power_consumption_total",
			Help: "Fenecon meter power consumption total in Watthours (ActiveConsumptionEnergy)",
		},
		commonLabels,
	))

	// ##########################################
	// Battery

	fp.newGaugeVec(&fp.prometheus.battery.charge, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_battery_charge_percent",
			Help: "Fenecon battery charge in percent (EssSoc)",
		},
		commonLabels,
	))

	fp.newGaugeVec(&fp.prometheus.battery.capacity, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_battery_capacity",
			Help: "Fenecon battery capacity in Watthours (EssCapacity)",
		},
		commonLabels,
	))

	fp.newGaugeVec(&fp.prometheus.battery.power, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_battery_power",
			Help: "Fenecon battery power load in Watts (EssActivePower)",
		},
		commonLabels,
	))

	fp.newGaugeVec(&fp.prometheus.battery.powerPhase, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_battery_power_phase",
			Help: "Fenecon battery power load in Watts (EssActivePowerLx)",
		},
		phaseLabels,
	))

	fp.newGaugeVec(&fp.prometheus.battery.powerChargeTotal, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_battery_power_charge_total",
			Help: "Fenecon battery power charge in Wattshours (EssActiveChargeEnergy)",
		},
		commonLabels,
	))

	fp.newGaugeVec(&fp.prometheus.battery.powerDischargeTotal, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_battery_power_discharge_total",
			Help: "Fenecon battery power discharge in Wattshours (EssActiveDischargeEnergy)",
		},
		commonLabels,
	))

	fp.newGaugeVec(&fp.prometheus.battery.powerDcChargeTotal, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_battery_power_dc_charge_total",
			Help: "Fenecon battery power dc charge in Wattshours (EssDcChargeEnergy)",
		},
		commonLabels,
	))

	fp.newGaugeVec(&fp.prometheus.battery.powerDcDischargeTotal, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_battery_power_dc_discharge_total",
			Help: "Fenecon battery power dc discharge in Wattshours (EssDcDischargeEnergy)",
		},
		commonLabels,
	))

	fp.newGaugeVec(&fp.prometheus.battery.allowedChargePower, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_battery_allowed_charge_power",
			Help: "Fenecon battery allowed scharge power Watts (AllowedChargePower)",
		},
		commonLabels,
	))

	fp.newGaugeVec(&fp.prometheus.battery.allowedDischargePower, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_battery_allowed_discharge_power",
			Help: "Fenecon battery allowed discharge power Watts (AllowedDischargePower)",
		},
		commonLabels,
	))

	// ##########################################
	// Grid

	fp.newGaugeVec(&fp.prometheus.grid.mode, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_grid_mode",
			Help: "Fenecon grid mode (0=undefined, 1=On-Grid, 2=Off-Grid; GridActivePower)",
		},
		commonLabels,
	))

	fp.newGaugeVec(&fp.prometheus.grid.power, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_grid_power",
			Help: "Fenecon grid power load in Watts (GridActivePower)",
		},
		commonLabels,
	))

	fp.newGaugeVec(&fp.prometheus.grid.powerPhase, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_grid_power_phase",
			Help: "Fenecon grid power load in Watts (GridActivePowerLx)",
		},
		phaseLabels,
	))

	fp.newGaugeVec(&fp.prometheus.grid.powerBuyTotal, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_grid_power_buy_total",
			Help: "Fenecon grid power buy in Wattshours (GridBuyActiveEnergy)",
		},
		commonLabels,
	))

	fp.newGaugeVec(&fp.prometheus.grid.powerSellTotal, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_grid_power_sell_total",
			Help: "Fenecon grid power sell in Wattshours (GridSellActiveEnergy)",
		},
		commonLabels,
	))

	// ##########################################
	// Production

	fp.newGaugeVec(&fp.prometheus.production.power, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_production_power",
			Help: "Fenecon production power load in Watts (ProductionActivePower)",
		},
		commonLabels,
	))

	fp.newGaugeVec(&fp.prometheus.production.powerPhase, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_production_power_phase",
			Help: "Fenecon production power load in Watts (ProductionAcActivePowerLx)",
		},
		phaseLabels,
	))

	fp.newGaugeVec(&fp.prometheus.production.powerAc, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_production_power_ac",
			Help: "Fenecon production power load in Watts (ProductionAcActivePower)",
		},
		commonLabels,
	))

	fp.newGaugeVec(&fp.prometheus.production.powerDc, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_production_power_dc",
			Help: "Fenecon production power load in Watts (ProductionDcActualPower)",
		},
		commonLabels,
	))

	fp.newGaugeVec(&fp.prometheus.production.powerTotal, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_production_power_total",
			Help: "Fenecon production power load in Watthours (ProductionActiveEnergy)",
		},
		commonLabels,
	))

	fp.newGaugeVec(&fp.prometheus.production.powerAcTotal, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_production_power_ac_total",
			Help: "Fenecon production power load in Watthours (ProductionAcActiveEnergy)",
		},
		commonLabels,
	))

	fp.newGaugeVec(&fp.prometheus.production.powerDcTotal, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_production_power_dc_total",
			Help: "Fenecon production power load in Watthours (ProductionDcActiveEnergy)",
		},
		commonLabels,
	))

	fp.newGaugeVec(&fp.prometheus.production.maxActualPower, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_production_max_actual_power",
			Help: "Fenecon production max acutal power Watts (MaxActualPower)",
		},
		commonLabels,
	))

	// ##########################################
	// Consumer

	fp.newGaugeVec(&fp.prometheus.consumption.power, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_consumption_power",
			Help: "Fenecon consumption power load in Watts (ConsumptionActivePower)",
		},
		commonLabels,
	))

	fp.newGaugeVec(&fp.prometheus.consumption.powerPhase, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_consumption_power_phase",
			Help: "Fenecon consumption power load in Watts (ConsumptionActivePowerLX)",
		},
		phaseLabels,
	))

	fp.newGaugeVec(&fp.prometheus.consumption.powerTotal, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_consumption_power_total",
			Help: "Fenecon consumption power load in Watts (ConsumptionActiveEnergy)",
		},
		commonLabels,
	))
}

func (fp *FeneconProber) newGaugeVec(dest **prometheus.GaugeVec, def *prometheus.GaugeVec) {
	(*dest) = def
	fp.registry.MustRegister(def)
}
