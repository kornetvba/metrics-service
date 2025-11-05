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

func (mh *MetricHandler) MetricPost(w http.ResponseWriter, r *http.Request) {
	typeMc := chi.URLParam(r, "type_metric")
	nameMc := chi.URLParam(r, "name_metric")
	valueMc := chi.URLParam(r, "value_metric")

	switch typeMc {
	case "counter":
		valInt, err := strconv.Atoi(valueMc)
		if err != nil {

			w.WriteHeader(http.StatusBadRequest)
			return
		}
		mh.storage.UpdateCounter(nameMc, int64(valInt))

	case "gauge":
		valFloat, err := strconv.ParseFloat(valueMc, 64)
		if err != nil {

			w.WriteHeader(http.StatusBadRequest)

			return
		}
		mh.storage.SetGauge(nameMc, valFloat)

	default:

		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("methic post successful %s %s %v", typeMc, nameMc, valueMc)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func (mh *MetricHandler) MetricGet(w http.ResponseWriter, r *http.Request) {
	typeMc := chi.URLParam(r, "type_metric")
	nameMc := chi.URLParam(r, "name_metric")
	value, err := mh.storage.GetMetric(typeMc, nameMc)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	valueStr := fmt.Sprintf("%v", value)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	log.Printf("methic get successful %s %s %v", typeMc, nameMc, valueStr)

	w.Write([]byte(valueStr))
}

//go:embed templates/*.html
var templateFS embed.FS

func (mh *MetricHandler) GetAllMetricsHTML(w http.ResponseWriter, _ *http.Request) {
	counter, gauge := mh.storage.GetAllMetrics()
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
