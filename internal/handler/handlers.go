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

func MetricPost(w http.ResponseWriter, r *http.Request) {
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
		models.MemStorageGlobal.SetCounter(nameMc, int64(valInt))

	case "gauge":
		valFloat, err := strconv.ParseFloat(valueMc, 64)
		if err != nil {

			w.WriteHeader(http.StatusBadRequest)

			return
		}
		models.MemStorageGlobal.SetGauge(nameMc, valFloat)

	default:

		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Print(fmt.Sprintf("methic add successful %s %s %v", typeMc, nameMc, valueMc))
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func MetricGet(w http.ResponseWriter, r *http.Request) {
	typeMc := chi.URLParam(r, "type_metric")
	nameMc := chi.URLParam(r, "name_metric")
	value, err := models.MemStorageGlobal.GetMetric(typeMc, nameMc)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	valueStr := fmt.Sprintf("%v", value)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	log.Print(fmt.Sprintf("methic add successful %s %s %v", typeMc, nameMc, valueStr))

	w.Write([]byte(valueStr))
}

//go:embed templates/*.html
var templateFS embed.FS

func GetAllMetricsHTML(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Counters map[string]int64
		Gauges   map[string]float64
	}{
		Counters: models.MemStorageGlobal.Counter,
		Gauges:   models.MemStorageGlobal.Gauge,
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
