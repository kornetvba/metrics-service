package handlers

import (
	"embed"
	"fmt"
	"github.com/go-chi/chi/v5"
	models "github.com/kornetvba/metrics-service/internal/storage"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type MetricHandler struct {
	storage models.Storage
}

func NewMetricHandler(storage models.Storage) *MetricHandler {
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
		valInt, err := strconv.Atoi(metricValue)
		if err != nil {

			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.storage.UpdateCounter(metricName, int64(valInt))

	case "gauge":
		valFloat, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {

			w.WriteHeader(http.StatusBadRequest)

			return
		}
		h.storage.SetGauge(metricName, valFloat)

	default:

		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("methic post successful %s %s %v", metricType, metricName, metricValue)
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
	w.WriteHeader(http.StatusOK)
	log.Printf("methic get successful %s %s %v", metricType, metricValue, valueStr)

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

}
