package handlers

import (
	models "github.com/kornetvba/metrics-service/internal/model"
	"net/http"
	"strconv"
)

func MetricPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	typeMc := r.PathValue("type_metric")
	nameMc := r.PathValue("name_metric")
	valueMc := r.PathValue("value_metric")

	switch typeMc {
	case "counter":
		valInt, err := strconv.Atoi(valueMc)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		models.MemStorageGlobal.IncrementCounter(nameMc, int64(valInt))
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
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

}
