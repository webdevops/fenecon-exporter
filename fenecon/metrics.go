package fenecon

import (
	"github.com/prometheus/client_golang/prometheus"
)

type (
	feneconMetrics struct {
		info   *prometheus.GaugeVec
		status *prometheus.GaugeVec

		batteryCharge                *prometheus.GaugeVec
		batteryCapacity              *prometheus.GaugeVec
		batteryPower                 *prometheus.GaugeVec
		batteryPowerPhase            *prometheus.GaugeVec
		batteryPowerChargeTotal      *prometheus.GaugeVec
		batteryPowerDischargeTotal   *prometheus.GaugeVec
		batteryPowerDcChargeTotal    *prometheus.GaugeVec
		batteryPowerDcDischargeTotal *prometheus.GaugeVec

		gridMode           *prometheus.GaugeVec
		gridPower          *prometheus.GaugeVec
		gridPowerPhase     *prometheus.GaugeVec
		gridPowerBuyTotal  *prometheus.GaugeVec
		gridPowerSellTotal *prometheus.GaugeVec

		productionPower        *prometheus.GaugeVec
		productionPowerPhase   *prometheus.GaugeVec
		productionPowerAc      *prometheus.GaugeVec
		productionPowerDc      *prometheus.GaugeVec
		productionPowerTotal   *prometheus.GaugeVec
		productionPowerAcTotal *prometheus.GaugeVec
		productionPowerDcTotal *prometheus.GaugeVec

		consumptionPower      *prometheus.GaugeVec
		consumptionPowerPhase *prometheus.GaugeVec
		consumptionPowerTotal *prometheus.GaugeVec
	}
)

func (fp *FeneconProber) initMetrics() {
	commonLabels := []string{"target"}
	phaseLabels := []string{"target", "phase"}

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
	// Battery

	fp.prometheus.batteryCharge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_battery_charge_percent",
			Help: "Fenecon battery charge in percent (EssSoc)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.batteryCharge)

	fp.prometheus.batteryCapacity = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_battery_capacity",
			Help: "Fenecon battery capacity in Watthours (EssCapacity)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.batteryCapacity)

	fp.prometheus.batteryPower = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_battery_power",
			Help: "Fenecon battery power load in Watts (EssActivePower)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.batteryPower)

	fp.prometheus.batteryPowerPhase = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_battery_power_phase",
			Help: "Fenecon battery power load in Watts (EssActivePowerLx)",
		},
		phaseLabels,
	)
	fp.registry.MustRegister(fp.prometheus.batteryPowerPhase)

	fp.prometheus.batteryPowerChargeTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_battery_power_charge_total",
			Help: "Fenecon battery power charge in Wattshours (EssActiveChargeEnergy)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.batteryPowerChargeTotal)

	fp.prometheus.batteryPowerDischargeTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_battery_power_discharge_total",
			Help: "Fenecon battery power discharge in Wattshours (EssActiveDischargeEnergy)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.batteryPowerDischargeTotal)

	fp.prometheus.batteryPowerDcChargeTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_battery_power_dc_charge_total",
			Help: "Fenecon battery power dc charge in Wattshours (EssDcChargeEnergy)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.batteryPowerDcChargeTotal)

	fp.prometheus.batteryPowerDcDischargeTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_battery_power_dc_discharge_total",
			Help: "Fenecon battery power dc discharge in Wattshours (EssDcDischargeEnergy)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.batteryPowerDcDischargeTotal)

	// ##########################################
	// Grid

	fp.prometheus.gridMode = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_grid_mode",
			Help: "Fenecon grid mode (0=undefined, 1=On-Grid, 2=Off-Grid; GridActivePower)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.gridMode)

	fp.prometheus.gridPower = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_grid_power",
			Help: "Fenecon grid power load in Watts (GridActivePower)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.gridPower)

	fp.prometheus.gridPowerPhase = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_grid_power_phase",
			Help: "Fenecon grid power load in Watts (GridActivePowerLx)",
		},
		phaseLabels,
	)
	fp.registry.MustRegister(fp.prometheus.gridPowerPhase)

	fp.prometheus.gridPowerBuyTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_grid_power_buy_total",
			Help: "Fenecon grid power buy in Wattshours (GridBuyActiveEnergy)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.gridPowerBuyTotal)

	fp.prometheus.gridPowerSellTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_grid_power_sell_total",
			Help: "Fenecon grid power sell in Wattshours (GridSellActiveEnergy)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.gridPowerSellTotal)

	// ##########################################
	// Production

	fp.prometheus.productionPower = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_production_power",
			Help: "Fenecon production power load in Watts (ProductionActivePower)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.productionPower)

	fp.prometheus.productionPowerPhase = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_production_power_phase",
			Help: "Fenecon production power load in Watts (ProductionAcActivePowerLx)",
		},
		phaseLabels,
	)
	fp.registry.MustRegister(fp.prometheus.productionPowerPhase)

	fp.prometheus.productionPowerAc = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_production_power_ac",
			Help: "Fenecon production power load in Watts (ProductionAcActivePower)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.productionPowerAc)

	fp.prometheus.productionPowerDc = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_production_power_dc",
			Help: "Fenecon production power load in Watts (ProductionDcActualPower)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.productionPowerDc)

	fp.prometheus.productionPowerTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_production_power_total",
			Help: "Fenecon production power load in Watthours (ProductionActiveEnergy)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.productionPowerTotal)

	fp.prometheus.productionPowerAcTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_production_power_ac_total",
			Help: "Fenecon production power load in Watthours (ProductionAcActiveEnergy)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.productionPowerAcTotal)

	fp.prometheus.productionPowerDcTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_production_power_dc_total",
			Help: "Fenecon production power load in Watthours (ProductionDcActiveEnergy)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.productionPowerDcTotal)

	// ##########################################
	// Consumer

	fp.prometheus.consumptionPower = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_consumption_power",
			Help: "Fenecon consumption power load in Watts (ConsumptionActivePower)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.consumptionPower)

	fp.prometheus.consumptionPowerPhase = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_consumption_power_phase",
			Help: "Fenecon consumption power load in Watts (ConsumptionActivePowerLX)",
		},
		phaseLabels,
	)
	fp.registry.MustRegister(fp.prometheus.consumptionPowerPhase)

	fp.prometheus.consumptionPowerTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_consumption_power_total",
			Help: "Fenecon consumption power load in Watts (ConsumptionActiveEnergy)",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.consumptionPowerTotal)

}
