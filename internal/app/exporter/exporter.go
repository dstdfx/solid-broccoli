package exporter

import (
	"sync"

	"github.com/dstdfx/solid-broccoli/internal/app/exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
)

// APIExporter wraps all the API collectors and provides a single global
// exporter to extract metrics out of. It also ensures that the collection
// is done in a thread-safe manner and implements a prometheus.Collector
// interface in order to register it correctly.
type APIExporter struct {
	mu         sync.Mutex
	collectors []prometheus.Collector
}

type NewAPIExporterOpts struct {
	BuildGitCommit string
	BuildGitTag    string
	BuildDate      string
	BuildCompiler  string
}

// NewAPIExporter returns a reference to a new instance of APIExporter.
func NewAPIExporter(opts *NewAPIExporterOpts) *APIExporter {
	return &APIExporter{
		collectors: []prometheus.Collector{
			collector.NewBuildInfoCollector(&collector.NewBuildInfoCollectorOpts{
				BuildGitCommit: opts.BuildGitCommit,
				BuildGitTag:    opts.BuildGitTag,
				BuildDate:      opts.BuildDate,
				BuildCompiler:  opts.BuildCompiler,
			}),
		},
	}
}

// Describe sends all the descriptors of the collectors included to
// the provided channel.
func (v *APIExporter) Describe(ch chan<- *prometheus.Desc) {
	for _, c := range v.collectors {
		c.Describe(ch)
	}
}

// Collect sends the collected metrics from each of the collectors to the
// Prometheus.
// Collect could be called several times concurrently and thus its run is
// protected by a single mutex.
func (v *APIExporter) Collect(ch chan<- prometheus.Metric) {
	v.mu.Lock()
	defer v.mu.Unlock()

	for _, c := range v.collectors {
		c.Collect(ch)
	}
}
