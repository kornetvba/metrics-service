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

func DecodeMetricBody(metricName string, metricValue interface{}) (*metrics.Metric, error) {
	metric := &metrics.Metric{}
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

	return metric, nil
}

func ClientMetric(timeDelay time.Duration) error {
	log.Print("agent running!")
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	for {
		metrics := []metrics.Metric{}
		gaugeMap := globalMetrics.ToMap()
		for k, v := range gaugeMap {
			body, err := DecodeMetricBody(k, v)
			if err != nil {
				log.Print(err)
				time.Sleep(2 * time.Second)
				continue
			}
			metrics = append(metrics, *body)
		}
		data, err := json.Marshal(metrics)
		if err != nil {
			log.Print(err)
			time.Sleep(2 * time.Second)
			continue
		}
		dataCompress, err := CompressData(&data)
		if err != nil {
			log.Print(err)
			time.Sleep(2 * time.Second)
			continue
		}

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/updates/", agent.AddrAgent.String()), bytes.NewBuffer(dataCompress))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Encoding", "gzip")

		resp, err := client.Do(req)
		if err != nil {
			log.Print(err)
			time.Sleep(2 * time.Second)
			continue
		}
		resp.Body.Close()

		time.Sleep(timeDelay)
	}
}
