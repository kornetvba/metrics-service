package agent

import (
	"bytes"
	"context"
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

func ClientMetric(ctx context.Context, timeDelay time.Duration) error {
	log.Print("agent running!")
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	ticker := time.NewTicker(timeDelay)

	for {
		select {
		case <-ctx.Done():
			log.Println("Отправка метрик завершилась")
			return ctx.Err()
		case <-ticker.C:
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
				bodyCompress, err := CompressData(&body)
				if err != nil {
					return err
				}
				req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/update/", agent.AddrAgent.String()), bytes.NewBuffer(bodyCompress))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Content-Encoding", "gzip")

				if err != nil {
					continue
				}

				_, err = client.Do(req)

				if err != nil {
					continue
				}

			}
		}
	}

}
