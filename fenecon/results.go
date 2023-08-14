package fenecon

type ResultCommon struct {
	Address    string  `json:"address"`
	Type       string  `json:"type"`
	AccessMode string  `json:"accessMode"`
	Text       string  `json:"text"`
	Unit       string  `json:"unit"`
	Value      float64 `json:"value"`
}
