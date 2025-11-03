package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMetricPost(t *testing.T) {
	type data struct {
		typeMc string
		nameMc string
		valMc  interface{}
	}

	tableTests := []struct {
		name        string
		method      string
		status      int
		data        data
		contentType string
	}{
		{
			name:        "test1",
			method:      http.MethodPost,
			status:      200,
			contentType: "text/plain",
			data: data{
				typeMc: "counter",
				nameMc: "dsa",
				valMc:  50,
			},
		},
		{
			name:        "test2",
			method:      http.MethodPost,
			status:      200,
			contentType: "text/plain",
			data: data{
				typeMc: "gauge",
				nameMc: "dsa",
				valMc:  50.4,
			},
		},
		{
			name:        "test3",
			method:      http.MethodPost,
			status:      http.StatusBadRequest,
			contentType: "text/plain",
			data: data{
				typeMc: "gaug3e",
				nameMc: "dsa",
				valMc:  50.4,
			},
		},
	}

	testMx := chi.NewRouter()
	testMx.HandleFunc("/update/{type_metric}/{name_metric}/{value_metric}", MetricPost)

	for _, tt := range tableTests {

		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.data.typeMc == "counter" {
				req = httptest.NewRequest(tt.method, fmt.Sprintf("/update/%s/%s/%d", tt.data.typeMc, tt.data.nameMc, tt.data.valMc), nil)
			} else if tt.data.typeMc == "gauge" {
				req = httptest.NewRequest(tt.method, fmt.Sprintf("/update/%s/%s/%f", tt.data.typeMc, tt.data.nameMc, tt.data.valMc), nil)
			} else {
				req = httptest.NewRequest(tt.method, fmt.Sprintf("/update/%s/%s/%f", tt.data.typeMc, tt.data.nameMc, tt.data.valMc), nil)
			}

			w := httptest.NewRecorder()
			testMx.ServeHTTP(w, req)
			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, tt.status, res.StatusCode)
			if res.StatusCode == http.StatusOK {
				assert.Equal(t, tt.contentType, res.Header.Get("Content-Type"))
			}

		})
	}
}
