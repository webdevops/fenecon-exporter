package fenecon

import (
	"github.com/prometheus/client_golang/prometheus"
)

type (
	feneconMetrics struct {
		info   *prometheus.GaugeVec
		status *prometheus.GaugeVec

		batteryCharge              *prometheus.GaugeVec
		batteryPower               *prometheus.GaugeVec
		batteryPowerPhase          *prometheus.GaugeVec
		batteryPowerChargeTotal    *prometheus.GaugeVec
		batteryPowerDischargeTotal *prometheus.GaugeVec

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
			Help: "Fenecon status",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.status)

	// ##########################################
	// Battery

	fp.prometheus.batteryCharge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_battery_charge_percent",
			Help: "Fenecon battery charge in percent",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.batteryCharge)

	fp.prometheus.batteryPower = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_battery_power",
			Help: "Fenecon battery power load in Watts",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.batteryPower)

	fp.prometheus.batteryPowerPhase = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_battery_power_phase",
			Help: "Fenecon battery power load in Watts",
		},
		phaseLabels,
	)
	fp.registry.MustRegister(fp.prometheus.batteryPowerPhase)

	fp.prometheus.batteryPowerChargeTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_battery_power_charge_total",
			Help: "Fenecon battery power charge in Wattshours",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.batteryPowerChargeTotal)

	fp.prometheus.batteryPowerDischargeTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_battery_power_discharge_total",
			Help: "Fenecon battery power discharge in Wattshours",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.batteryPowerDischargeTotal)

	// ##########################################
	// Grid

	fp.prometheus.gridPower = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_grid_power",
			Help: "Fenecon grid power load in Watts",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.gridPower)

	fp.prometheus.gridPowerPhase = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_grid_power_phase",
			Help: "Fenecon grid power load in Watts",
		},
		phaseLabels,
	)
	fp.registry.MustRegister(fp.prometheus.gridPowerPhase)

	fp.prometheus.gridPowerBuyTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_grid_power_buy_total",
			Help: "Fenecon grid power buy in Wattshours",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.gridPowerBuyTotal)

	fp.prometheus.gridPowerSellTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_grid_power_sell_total",
			Help: "Fenecon grid power sell in Wattshours",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.gridPowerSellTotal)

	// ##########################################
	// Production

	fp.prometheus.productionPower = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_production_power",
			Help: "Fenecon production power load in Watts",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.productionPower)

	fp.prometheus.productionPowerPhase = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_production_power_phase",
			Help: "Fenecon production power load in Watts",
		},
		phaseLabels,
	)
	fp.registry.MustRegister(fp.prometheus.productionPowerPhase)

	fp.prometheus.productionPowerAc = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_production_power_ac",
			Help: "Fenecon production power load in Watts",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.productionPowerAc)

	fp.prometheus.productionPowerDc = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_production_power_dc",
			Help: "Fenecon production power load in Watts",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.productionPowerDc)

	fp.prometheus.productionPowerTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_production_power_total",
			Help: "Fenecon production power load in Watthours",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.productionPowerTotal)

	fp.prometheus.productionPowerAcTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_production_power_ac_total",
			Help: "Fenecon production power load in Watthours",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.productionPowerAcTotal)

	fp.prometheus.productionPowerDcTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_production_power_dc_total",
			Help: "Fenecon production power load in Watthours",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.productionPowerDcTotal)

	// ##########################################
	// Consumer

	fp.prometheus.consumptionPower = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_consumption_power",
			Help: "Fenecon consumption power load in Watts",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.consumptionPower)

	fp.prometheus.consumptionPowerPhase = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_consumption_power_phase",
			Help: "Fenecon consumption power load in Watts",
		},
		phaseLabels,
	)
	fp.registry.MustRegister(fp.prometheus.consumptionPowerPhase)

	fp.prometheus.consumptionPowerTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fenecon_consumption_power_total",
			Help: "Fenecon consumption power load in Watts",
		},
		commonLabels,
	)
	fp.registry.MustRegister(fp.prometheus.consumptionPowerTotal)

}
