package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

func getPrometheusTimeout(r *http.Request, defaultTimeout float64) (timeout float64, err error) {
	// If a timeout is configured via the Prometheus header, add it to the request.
	if v := r.Header.Get("X-Prometheus-Scrape-Timeout-Seconds"); v != "" {
		timeout, err = strconv.ParseFloat(v, 64)
		if err != nil {
			return
		}
	}
	if timeout == 0 {
		timeout = defaultTimeout
	}

	return
}

func paramsGetRequired(params url.Values, name string) (value string, err error) {
	value = params.Get(name)
	if value == "" {
		err = fmt.Errorf("parameter \"%v\" is missing", name)
	}

	return
}
