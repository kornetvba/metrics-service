package handlers

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/kornetvba/metrics-service/internal/storage/psql"
	"io"

	metrics "github.com/kornetvba/metrics-service/internal/model"
	"github.com/kornetvba/metrics-service/internal/repository"

	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

type MetricHandler struct {
	storage repository.Storage
}

func NewMetricHandler(storage repository.Storage) *MetricHandler {
	return &MetricHandler{
		storage: storage,
	}
}

func (h *MetricHandler) MetricPost(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type_metric")
	metricName := chi.URLParam(r, "name_metric")
	metricValue := chi.URLParam(r, "value_metric")

	switch metricType {
	case "counter":
		valInt, err := strconv.ParseInt(string(metricValue), 10, 64)
		if err != nil {

			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = h.storage.UpdateCounter(metricName, int64(valInt))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

	case "gauge":
		valFloat, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {

			w.WriteHeader(http.StatusBadRequest)

			return
		}
		_, err = h.storage.SetGauge(metricName, valFloat)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

	default:

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func (h *MetricHandler) MetricGet(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type_metric")
	metricValue := chi.URLParam(r, "name_metric")
	value, err := h.storage.GetMetric(metricType, metricValue)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	valueStr := fmt.Sprintf("%v", value)
	w.Header().Set("Content-Type", "text/plain")

	w.Write([]byte(valueStr))
}

//go:embed templates/*.html
var templateFS embed.FS

func (h *MetricHandler) GetAllMetricsHTML(w http.ResponseWriter, _ *http.Request) {
	counter, gauge := h.storage.GetAllMetrics()
	data := struct {
		Counters map[string]int64
		Gauges   map[string]float64
	}{
		Counters: counter,
		Gauges:   gauge,
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl := template.Must(template.ParseFS(templateFS, "templates/*.html"))
	//создаем html-шаблон
	err := tmpl.ExecuteTemplate(w, "metrics.html", data)
	if err != nil {
		log.Print(err)
		return
	}
	//w.WriteHeader(http.StatusOK)

}

func (h *MetricHandler) MetricPostJSON(w http.ResponseWriter, r *http.Request) {
	metric := metrics.Metric{}
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if metric.Delta != nil && metric.MType == "counter" {
		updateDelta, err := h.storage.UpdateCounter(metric.ID, *metric.Delta)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		metric.Delta = &updateDelta
	} else if metric.Value != nil && metric.MType == "gauge" {
		setValue, err := h.storage.SetGauge(metric.ID, *metric.Value)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		metric.Value = &setValue
	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	dataResp, err := json.Marshal(metric)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dataResp)
}

func (h *MetricHandler) MetricGetJSON(w http.ResponseWriter, r *http.Request) {
	metric := metrics.Metric{}

	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	val, err := h.storage.GetMetric(metric.MType, metric.ID)
	if err != nil {

		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch valueType := val.(type) {
	case int64:
		delta := valueType
		metric.Delta = &delta
	case float64:
		value := valueType
		metric.Value = &value
	}
	data, err := json.Marshal(metric)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *MetricHandler) PingHandler(w http.ResponseWriter, _ *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := psql.DB.PingContext(ctx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *MetricHandler) MetricsPostJSON(w http.ResponseWriter, r *http.Request) {
	metrics := make([]metrics.Metric, 0)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &metrics)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(metrics) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.storage.AppendMetrics(metrics)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}
