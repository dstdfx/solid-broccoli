package collector

import (
	"regexp"
	"testing"

	"github.com/dstdfx/solid-broccoli/internal/pkg/testutils"
	"github.com/stretchr/testify/assert"
)

func TestBuildInfoCollector(t *testing.T) {
	// Prepare expected Prometheus metrics.
	expected := []*regexp.Regexp{
		regexp.MustCompile(`build_info{build_date="20190909",compiler="go1.15",git_commit="7ede33ee9cbc22290d904ced379c3541002d816e",git_tag="v0.1.0"} 1`),
	}

	// Register NewInstanceCollector and run Prometheus server.
	collector := NewBuildInfoCollector(&NewBuildInfoCollectorOpts{
		BuildGitCommit: "7ede33ee9cbc22290d904ced379c3541002d816e",
		BuildGitTag:    "v0.1.0",
		BuildDate:      "20190909",
		BuildCompiler:  "go1.15",
	})
	prometheusEnv, err := testutils.SetupPrometheus(collector)
	assert.NoError(t, err)
	defer prometheusEnv.TearDown(collector)

	// Retrieve data and compare it with the expected metrics.
	testutils.HandlePrometheusMetric(t, &testutils.HandlePrometheusMetricOpts{
		Env:      prometheusEnv,
		Expected: expected,
	})
}
