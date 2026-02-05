package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/kornetvba/metrics-service/internal/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type reqData struct {
	typeMc string
	nameMc string
	valMc  interface{}
}

func testGenerateURL(data *reqData, method string) (req *http.Request) {
	if method == http.MethodPost {
		if data.typeMc == "counter" {
			req = httptest.NewRequest(method, fmt.Sprintf("/update/%s/%s/%d", data.typeMc, data.nameMc, data.valMc), nil)
			return
		} else if data.typeMc == "gauge" {
			req = httptest.NewRequest(method, fmt.Sprintf("/update/%s/%s/%f", data.typeMc, data.nameMc, data.valMc), nil)
			return
		} else {
			req = httptest.NewRequest(method, fmt.Sprintf("/update/%s/%s/%f", data.typeMc, data.nameMc, data.valMc), nil)
			return
		}
	}
	if method == http.MethodGet {
		if data.typeMc == "counter" {
			req = httptest.NewRequest(method, fmt.Sprintf("/value/%s/%s", data.typeMc, data.nameMc), nil)
			return
		} else if data.typeMc == "gauge" {
			req = httptest.NewRequest(method, fmt.Sprintf("/value/%s/%s", data.typeMc, data.nameMc), nil)
			return
		} else {
			req = httptest.NewRequest(method, fmt.Sprintf("/value/%s/%s", data.typeMc, data.nameMc), nil)
			return
		}
	}
	return
}

func TestMetricPost(t *testing.T) {

	tableTests := []struct {
		name        string
		method      string
		status      int
		data        reqData
		contentType string
	}{
		{
			name:        "test1",
			method:      http.MethodPost,
			status:      200,
			contentType: "text/plain",
			data: reqData{
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
			data: reqData{
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
			data: reqData{
				typeMc: "gaug3e",
				nameMc: "dsa",
				valMc:  50.4,
			},
		},
	}
	testMetric := NewMetricHandler(memory.NewMemStorage())
	testMx := chi.NewRouter()
	testMx.HandleFunc("/update/{type_metric}/{name_metric}/{value_metric}", testMetric.MetricPost)

	for _, tt := range tableTests {

		t.Run(tt.name, func(t *testing.T) {
			req := testGenerateURL(&tt.data, http.MethodPost)

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

func TestMetricGet(t *testing.T) {

	tableTest := []struct {
		name           string
		Data           reqData
		statusCode     int
		incrementIndex int
	}{
		{
			name: "test1",
			Data: reqData{
				typeMc: "counter",
				nameMc: "danil",
				valMc:  50,
			},
		},
		{
			name: "test2",
			Data: reqData{
				typeMc: "gauge",
				nameMc: "nikita",
				valMc:  50.5,
			},
		},
		{
			name: "test2",
			Data: reqData{
				typeMc: "counter",
				nameMc: "nikita",
				valMc:  100,
			},
			incrementIndex: 3,
		},
	}
	testMetric := NewMetricHandler(memory.NewMemStorage())
	rTest := chi.NewRouter()
	rTest.Post("/update/{type_metric}/{name_metric}/{value_metric}", testMetric.MetricPost)
	rTest.Get("/value/{type_metric}/{name_metric}", testMetric.MetricGet)
	for _, tt := range tableTest {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.incrementIndex != 0 {
				for i := 1; i <= tt.incrementIndex; i++ {
					req = testGenerateURL(&tt.Data, http.MethodPost)
					w := httptest.NewRecorder()
					rTest.ServeHTTP(w, req)
				}
			} else {
				req = testGenerateURL(&tt.Data, http.MethodPost)
				w := httptest.NewRecorder()
				rTest.ServeHTTP(w, req)
			}

			req = testGenerateURL(&tt.Data, http.MethodGet)
			w := httptest.NewRecorder()
			rTest.ServeHTTP(w, req)
			resp := w.Result()
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err)

				if tt.Data.typeMc == "counter" {
					resVal, err := strconv.Atoi(string(body))
					require.NoError(t, err)
					if tt.incrementIndex != 0 {
						expectedVal := fmt.Sprintf("%v", tt.Data.valMc)
						val, err := strconv.Atoi(expectedVal)
						require.NoError(t, err)
						assert.Equal(t, val*tt.incrementIndex, resVal)
					} else {
						assert.Equal(t, tt.Data.valMc, resVal)
					}

				} else if tt.Data.valMc == "gauge" {
					resVal, err := strconv.ParseFloat(string(body), 64)
					require.NoError(t, err)

					assert.Equal(t, tt.Data.valMc, resVal)
				}
			} else {
				assert.Equal(t, http.StatusNotFound, resp.StatusCode)
			}

		})
	}

}

//func TestMetricGetMock(t *testing.T) {
//
//	tableTests := []struct {
//		name       string
//		statusCode int
//		par        struct {
//			metricType  string
//			metricValue string
//		}
//		resp interface{}
//	}{
//		{
//			name:       "testMock1",
//			statusCode: http.StatusOK,
//			par: struct {
//				metricType  string
//				metricValue string
//			}{metricType: "hello", metricValue: "test"},
//			resp: "hello test too",
//		},
//	}
//
//	for _, tt := range tableTests {
//		t.Run(tt.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//
//			m := mock_storage.NewMockStorage(ctrl)
//			m.EXPECT().GetMetric(tt.par.metricType, tt.par.metricValue).Return(tt.resp, nil)
//			app := NewMetricHandler(m)
//			w := httptest.NewRecorder()
//			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s/%s", tt.par.metricType, tt.par.metricValue), nil)
//			require.NoError(t, err)
//
//			router := chi.NewRouter()
//			router.Get("/{type_metric}/{name_metric}", app.MetricGet)
//
//			router.ServeHTTP(w, req)
//
//			resp := w.Result()
//			defer resp.Body.Close()
//
//			require.Equal(t, tt.statusCode, resp.StatusCode)
//
//		})
//	}
//
//}
