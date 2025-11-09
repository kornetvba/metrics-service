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
	switch metricValue.(type) {
	case int64:
		metric.MType = "counter"
		v, ok := metricValue.(int64)
		if !ok {
			return nil, errors.New("type counter access only int64")
		}
		metric.Delta = &v
	case float64:
		metric.MType = "gauge"
		v, ok := metricValue.(float64)
		if !ok {
			return nil, errors.New("type gauge access only float64")
		}
		metric.Value = &v
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
