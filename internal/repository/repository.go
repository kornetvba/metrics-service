package repository

import metrics "github.com/kornetvba/metrics-service/internal/model"

type Storage interface {
	UpdateCounter(string, int64) (int64, error)
	SetGauge(string, float64) (float64, error)
	GetMetric(string, string) (interface{}, error)
	GetAllMetrics() (map[string]int64, map[string]float64)
	AppendMetrics([]metrics.Metric) error
}

//INSERT INTO metrics (type_metric, name_metric,value_metric)
//VALUES('helo', 'frw', 10)
//ON CONFLICT (name_metric)
//DO UPDATE
//SET value_metric = metrics.value_metric + 10;
