//go:build windows

package promtailconvert_test

import (
	"testing"

	"github.com/grafana/alloy/internal/converter/internal/promtailconvert"
	"github.com/grafana/alloy/internal/converter/internal/test_common"
	_ "github.com/grafana/alloy/internal/static/metrics/instance" // Imported to override default values via the init function.
)

func TestConvert(t *testing.T) {
	test_common.TestDirectory(t, "testdata", ".yaml", true, []string{}, map[string]struct{}{}, promtailconvert.Convert)
	test_common.TestDirectory(t, "testdata_windows", ".yaml", true, []string{}, map[string]struct{}{}, promtailconvert.Convert)
}
