package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/webdevops/fenecon-exporter/fenecon"
)

const (
	DefaultTimeout = 30
)

func newFeneconProber(ctx context.Context, registry *prometheus.Registry, logger *zap.SugaredLogger) *fenecon.FeneconProber {
	sp := fenecon.New(ctx, registry, logger)
	sp.SetUserAgent(UserAgent + gitTag)
	if len(opts.Fenecon.Auth.Password) >= 1 {
		sp.SetHttpAuth(opts.Fenecon.Auth.Username, opts.Fenecon.Auth.Password)
	}

	return sp
}

func probeFenecon(w http.ResponseWriter, r *http.Request) {
	var (
		err            error
		timeoutSeconds float64
		target         string
	)

	// startTime := time.Now()
	contextLogger := buildContextLoggerFromRequest(r)
	registry := prometheus.NewRegistry()

	// If a timeout is configured via the Prometheus header, add it to the request.
	timeoutSeconds, err = getPrometheusTimeout(r, DefaultTimeout)
	if err != nil {
		contextLogger.Error(err)
		http.Error(w, fmt.Sprintf("failed to parse timeout from Prometheus header: %s", err), http.StatusBadRequest)
		return
	}

	if target, err = paramsGetRequired(r.URL.Query(), "target"); err != nil {
		contextLogger.Warnln(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds*float64(time.Second)))
	defer cancel()
	r = r.WithContext(ctx)

	prober := newFeneconProber(ctx, registry, contextLogger)
	prober.Run(target)

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

func buildContextLoggerFromRequest(r *http.Request) *zap.SugaredLogger {
	return logger.With(zap.String("requestPath", r.URL.Path))
}