package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/kornetvba/metrics-service/internal/agent"
	handlers "github.com/kornetvba/metrics-service/internal/handler"
	"github.com/kornetvba/metrics-service/internal/storage"
	"log"
	"net/http"
)

func main() {
	go func() {
		err := run() // сервер
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		agent.CollectMetrics(2) // сбор метрик
	}()

	go func() {
		err := agent.ClientMetric(10) // отправка метрик
		if err != nil {
			log.Fatal(err)
		}
	}()

	select {}
}

func run() error {
	log.Print("serv is running")
	handler := handlers.NewMetricHandler(storage.NewMemStorage())
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/", handler.GetAllMetricsHTML)
		r.Get("/value/{type_metric}/{name_metric}", handler.MetricGet)

		r.Post("/update/{type_metric}/{name_metric}/{value_metric}", handler.MetricPost)
	})

	return http.ListenAndServe(":8080", r)
}
