package testutils

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
)

// PrometheusEnv represents a Prometheus testing environment.
type PrometheusEnv struct {
	Server *httptest.Server
}

// SetupPrometheus prepares a new Prometheus testing environment.
func SetupPrometheus(collector prometheus.Collector) (*PrometheusEnv, error) {
	if err := prometheus.Register(collector); err != nil {
		return nil, err
	}

	return &PrometheusEnv{httptest.NewServer(promhttp.Handler())}, nil
}

// TearDown releases the Prometheus testing environment.
func (p *PrometheusEnv) TearDown(collector prometheus.Collector) {
	prometheus.Unregister(collector)
	p.Server.Close()
}

// HandlePrometheusMetricOpts contains options for the HandlePrometheusMetric.
type HandlePrometheusMetricOpts struct {
	Env      *PrometheusEnv
	Expected []*regexp.Regexp
}

// HandlePrometheusMetric makes a request against configured Prometheus and
// compares received metrics against expected set.
func HandlePrometheusMetric(t *testing.T, opts *HandlePrometheusMetricOpts) {
	resp, err := http.Get(opts.Env.Server.URL)
	assert.NoError(t, err)
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	for _, re := range opts.Expected {
		assert.Regexp(t, re, string(buf))
	}
}
