package collector

import "github.com/prometheus/client_golang/prometheus"

const (
	gitCommitLabel = "git_commit"
	gitTagLabel    = "git_tag"
	buildDateLabel = "build_date"
	compilerLabel  = "compiler"
)

type NewBuildInfoCollectorOpts struct {
	BuildGitCommit string
	BuildGitTag    string
	BuildDate      string
	BuildCompiler  string
}

// NewBuildInfoCollector is a collector with build information.
//
// Reference: https://github.com/prometheus/client_golang/blob/v1.1.0/prometheus/go_collector.go#L368
func NewBuildInfoCollector(opts *NewBuildInfoCollectorOpts) prometheus.Collector {
	c := &selfCollector{prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"build_info",
			"Build information about the app.",
			nil, prometheus.Labels{
				gitCommitLabel: opts.BuildGitCommit,
				gitTagLabel:    opts.BuildGitTag,
				buildDateLabel: opts.BuildDate,
				compilerLabel:  opts.BuildCompiler,
			},
		),
		prometheus.GaugeValue, 1)}
	c.init(c.self)

	return c
}

// selfCollector implements Collector for a single Metric so that the Metric
// collects itself. Add it as an anonymous field to a struct that implements
// Metric, and call init with the Metric itself as an argument.
//
// Reference: https://github.com/prometheus/client_golang/blob/v1.1.0/prometheus/collector.go#L98
type selfCollector struct {
	self prometheus.Metric
}

// init provides the selfCollector with a reference to the metric it is supposed
// to collect. It is usually called within the factory function to create a
// metric. See example.
//
// Reference: https://github.com/prometheus/client_golang/blob/v1.1.0/prometheus/collector.go#L105
func (c *selfCollector) init(self prometheus.Metric) {
	c.self = self
}

// Describe implements Collector.
func (c *selfCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.self.Desc()
}

// Collect implements Collector.
func (c *selfCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- c.self
}
