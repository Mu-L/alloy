//go:build (linux && arm64) || (linux && amd64)

package ebpf

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.uber.org/atomic"

	"github.com/grafana/alloy/internal/component/pyroscope"
	"github.com/grafana/alloy/internal/util"
	"github.com/grafana/pyroscope/ebpf/pprof"
	"github.com/grafana/pyroscope/ebpf/sd"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendProfilesConcurrently(t *testing.T) {
	tests := []struct {
		name                 string
		profilesCount        int
		delay                time.Duration
		collectionInterval   time.Duration
		expectedSendDuration time.Duration
		expectedSuccess      int
		expectedErrors       int
		expectedDrops        int
	}{
		{
			name:                 "send 64 profiles in 1 second",
			profilesCount:        64,
			delay:                500 * time.Millisecond,
			collectionInterval:   15 * time.Second,
			expectedSendDuration: time.Second, // 64 / 32 * 0.5s
			expectedSuccess:      64,
			expectedErrors:       0,
			expectedDrops:        0,
		},
		{
			name: "send 32 profiles in 500ms, " +
				"start sending another 32," +
				" fail with timeout, completely the rest 192",
			profilesCount:        256,
			delay:                500 * time.Millisecond,
			collectionInterval:   800 * time.Millisecond,
			expectedSendDuration: 800 * time.Millisecond,
			expectedSuccess:      32,
			expectedErrors:       32,
			expectedDrops:        192,
		},
	}
	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			reg := prometheus.NewRegistry()
			c := new(Component)
			c.metrics = newMetrics(reg)
			c.options.Logger = util.TestAlloyLogger(t)
			c.args.CollectInterval = td.collectionInterval
			successes := atomic.Uint32{}
			failures := atomic.Uint32{}
			c.appendable = pyroscope.NewFanout([]pyroscope.Appendable{
				pyroscope.AppendableFunc(func(ctx context.Context, labels labels.Labels, samples []*pyroscope.RawSample) error {
					after := time.After(td.delay)
					select {
					case <-after:
						successes.Add(1)
						return nil
					case <-ctx.Done():
						failures.Add(1)
						return ctx.Err()
					}
				}),
			}, "", reg)

			profiles := pprof.NewProfileBuilders(pprof.BuildersOptions{
				SampleRate: 97,
			})

			for i := 0; i < td.profilesCount; i++ {
				cid := fmt.Sprintf("cid_%d", i)
				pid := uint32(239 + i)
				target := sd.DiscoveryTarget(map[string]string{
					"service_name": fmt.Sprintf("service_%d", i),
				})
				profiles.AddSample(&pprof.ProfileSample{
					Target:      sd.NewTargetForTesting(cid, pid, target),
					Pid:         pid,
					SampleType:  pprof.SampleTypeCpu,
					Aggregation: pprof.SampleAggregated,
					Stack: []string{
						"func1", "func2",
					},
					Value: 42,
				})
			}

			t1 := time.Now()
			c.sendProfiles(t.Context(), profiles)
			duration := time.Since(t1)
			expectedDuration := td.expectedSendDuration
			diff := duration - expectedDuration
			if diff < 0 {
				diff = -diff
			}
			assert.Less(t, diff, 100*time.Millisecond)
			require.EqualValues(t, td.expectedErrors, failures.Load())
			require.EqualValues(t, td.expectedSuccess, successes.Load())
			require.EqualValues(t, float64(td.expectedDrops), gatherDrops(t, reg))
		})
	}
}

func gatherDrops(t *testing.T, reg *prometheus.Registry) float64 {
	gather, err := reg.Gather()
	require.NoError(t, err)

	for _, f := range gather {
		if *f.Name == "pyroscope_ebpf_pprofs_dropped_total" {
			require.Len(t, f.Metric, 1)
			c := f.Metric[0].GetCounter()
			return *c.Value
		}
	}
	require.Fail(t, "metric not found")
	return 0
}
