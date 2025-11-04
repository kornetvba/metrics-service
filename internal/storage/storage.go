package storage

import (
	"errors"
)

type MemStorage struct {
	Counter map[string]int64
	Gauge   map[string]float64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Counter: make(map[string]int64),
		Gauge:   make(map[string]float64),
	}
}

type AdderMetric interface {
	IncrementCounter(string, int64) int64
	SetGauge(string, float64) float64
}

func (ms *MemStorage) UpdateCounter(name string, i int64) {
	if _, ok := ms.Counter[name]; !ok {
		ms.Counter[name] = i
		return
	}
	ms.Counter[name] += i

}

func (ms *MemStorage) SetGauge(name string, i float64) {
	ms.Gauge[name] = i

}

func (ms *MemStorage) GetMetric(typeMc, nameMc string) (interface{}, error) {
	switch typeMc {
	case "counter":
		v, ok := ms.Counter[nameMc]
		if !ok {
			return nil, errors.New("metric value is not found")
		}

		return v, nil

	case "gauge":
		v, ok := ms.Gauge[nameMc]
		if !ok {
			return nil, errors.New("metric value is not found")
		}

		return v, nil

	}
	return nil, errors.New("metric type is not found")
}

var MemStorageGlobal = NewMemStorage()

// NOTE: Не усложняем пример, вводя иерархическую вложенность структур.
// Органичиваясь плоской моделью.
// Delta и Value объявлены через указатели,
// что бы отличать значение "0", от не заданного значения
// и соответственно не кодировать в структуру.
//type Metrics struct {
//	ID    string   `json:"id"`
//	MType string   `json:"type"`
//	Delta *int64   `json:"delta,omitempty"`
//	Value *float64 `json:"value,omitempty"`
//	Hash  string   `json:"hash,omitempty"`
//}
