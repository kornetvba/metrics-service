package agent

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kornetvba/metrics-service/internal/config/agent"
	metrics "github.com/kornetvba/metrics-service/internal/model"
	"log"
	"net/http"
	"time"
)

func DecodeMetricBody(metricName string, metricValue interface{}) ([]byte, error) {
	metric := metrics.Metric{}
	metric.ID = metricName
	switch val := metricValue.(type) {
	case int64:
		metric.MType = "counter"
		metric.Delta = &val
	case float64:
		metric.MType = "gauge"
		metric.Value = &val
	default:
		return nil, errors.New("value is not validate")
	}
	data, err := json.Marshal(metric)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func ClientMetric(timeDelay time.Duration) error {
	log.Print("agent running!")
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	for {
		if globalMetrics == nil {
			time.Sleep(2 * time.Second)
			continue
		}
		gaugeMap := globalMetrics.ToMap()
		for k, v := range gaugeMap {
			body, err := DecodeMetricBody(k, v)
			if err != nil {
				continue
			}

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/update/", agent.AddrAgent.String()), bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			if err != nil {
				continue
			}

			res, err := client.Do(req)

			if err != nil {
				continue
			}
			res.Body.Close()

		}
		time.Sleep(timeDelay)
	}

}
