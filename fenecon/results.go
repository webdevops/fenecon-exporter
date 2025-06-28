package fenecon

import (
	"encoding/json"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type (
	ResultWildcard []ResultCommon

	ResultCommon struct {
		Address    string      `json:"address"`
		Type       string      `json:"type"`
		AccessMode string      `json:"accessMode"`
		Text       string      `json:"text"`
		Unit       string      `json:"unit"`
		Value      ResultValue `json:"value"`
	}

	ResultValue struct {
		ValueNumeric *float64
		ValueString  *string
	}
)

func (v *ResultValue) UnmarshalJSON(data []byte) error {
	var (
		valFloat  float64
		valString string
	)

	// as numeric
	if err := json.Unmarshal(data, &valFloat); err == nil {
		v.ValueNumeric = &valFloat
		return nil
	}

	// as string
	if err := json.Unmarshal(data, &valString); err == nil {
		v.ValueString = &valString
		return nil
	}

	return nil
}

func (r *ResultWildcard) AddressParts(part int) []string {
	list := map[string]string{}

	for _, row := range *r {
		parts := strings.Split(row.Address, "/")
		if len(parts) >= part {
			v := parts[part]
			list[v] = v
		}
	}

	// unique
	ret := []string{}
	for _, v := range list {
		ret = append(ret, v)
	}

	return ret
}

func (r *ResultWildcard) Address(val ...string) *ResultCommon {
	address := strings.Join(val, "/")
	for _, row := range *r {
		if strings.EqualFold(row.Address, address) {
			return &row
		}
	}

	return &ResultCommon{}
}

func (r *ResultCommon) SetGauge(labels prometheus.Labels, gaugeVec *prometheus.GaugeVec) {
	if r.Value.ValueNumeric != nil {
		gaugeVec.With(labels).Set(*r.Value.ValueNumeric)
	}
}

func (r *ResultCommon) SetGaugeIfNotZero(labels prometheus.Labels, gaugeVec *prometheus.GaugeVec) {
	if r.Value.ValueNumeric != nil && *r.Value.ValueNumeric > 0 {
		gaugeVec.With(labels).Set(*r.Value.ValueNumeric)
	}
}
