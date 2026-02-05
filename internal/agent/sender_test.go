package agent

import (
	metrics "github.com/kornetvba/metrics-service/internal/model"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDecodeMetricBody(t *testing.T) {
	tableTests := []struct {
		name        string
		metricName  string
		metricValue interface{}
		data        metrics.Metric
	}{
		{
			name:        "test1",
			metricName:  "danil",
			metricValue: int64(10),
			data: metrics.Metric{
				Delta: func() *int64 { v := int64(10); return &v }(),
				ID:    "danil",
				MType: "counter",
			},
		},
		{
			name:        "test2",
			metricName:  "nikita",
			metricValue: float64(54.5),
			data: metrics.Metric{
				Value: func() *float64 { v := float64(54.5); return &v }(),
				ID:    "nikita",
				MType: "gauge",
			},
		},
	}
	for _, tt := range tableTests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := DecodeMetricBody(tt.metricName, tt.metricValue)
			require.NoError(t, err)

			require.Equal(t, tt.data, *resp)
		})
	}
}
