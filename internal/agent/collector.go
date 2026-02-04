package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kornetvba/metrics-service/internal/config/agent"
	metrics "github.com/kornetvba/metrics-service/internal/model"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

var globalMetrics = &metrics.Metrics{}

func UpdateRuntimeMetrics(metricsRM *metrics.Metrics) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	metricsRM.Alloc = float64(m.Alloc)
	metricsRM.BuckHashSys = float64(m.BuckHashSys)
	metricsRM.Frees = float64(m.Frees)
	metricsRM.GCCPUFraction = m.GCCPUFraction
	metricsRM.GCSys = float64(m.GCSys)
	metricsRM.HeapAlloc = float64(m.HeapAlloc)
	metricsRM.HeapIdle = float64(m.HeapIdle)
	metricsRM.HeapInuse = float64(m.HeapInuse)
	metricsRM.HeapObjects = float64(m.HeapObjects)
	metricsRM.HeapReleased = float64(m.HeapReleased)
	metricsRM.HeapSys = float64(m.HeapSys)
	metricsRM.LastGC = float64(m.LastGC)
	metricsRM.Lookups = float64(m.Lookups)
	metricsRM.MCacheInuse = float64(m.MCacheInuse)
	metricsRM.MCacheSys = float64(m.MCacheSys)
	metricsRM.MSpanInuse = float64(m.MSpanInuse)
	metricsRM.MSpanSys = float64(m.MSpanSys)
	metricsRM.Mallocs = float64(m.Mallocs)
	metricsRM.NextGC = float64(m.NextGC)
	metricsRM.NumForcedGC = float64(m.NumForcedGC)
	metricsRM.NumGC = float64(m.NumGC)
	metricsRM.OtherSys = float64(m.OtherSys)
	metricsRM.PauseTotalNs = float64(m.PauseTotalNs)
	metricsRM.StackInuse = float64(m.StackInuse)
	metricsRM.StackSys = float64(m.StackSys)
	metricsRM.Sys = float64(m.Sys)
	metricsRM.TotalAlloc = float64(m.TotalAlloc)

	metricsRM.RandomValue = rand.Float64()

}

func CollectMetrics(timeDelay time.Duration) {

	for {
		UpdateRuntimeMetrics(globalMetrics)
		metric := metrics.Metric{}
		metric.MType = "counter"
		delta := int64(1)
		metric.Delta = &delta
		metric.ID = "PollCount"
		data, err := json.Marshal(metric)
		if err != nil {
			log.Print(err)
			continue
		}
		dataCompr, err := CompressData(&data)
		if err != nil {
			log.Print(err)
			continue
		}

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/update/", agent.AddrAgent.String()), bytes.NewBuffer(dataCompr))
		if err != nil {
			log.Print(err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Encoding", "gzip")
		_, err = http.DefaultClient.Do(req)
		if err != nil {
			log.Print(err)
			time.Sleep(3 * time.Second)
			continue
		}

		time.Sleep(timeDelay)
	}

}
