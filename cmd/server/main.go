package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/kornetvba/metrics-service/internal/config"
	handlers "github.com/kornetvba/metrics-service/internal/handler"
	"github.com/kornetvba/metrics-service/internal/storage"
	"log"
	"net/http"
)

func main() {
	config.ParseFlagServer()
	go func() {
		err := run() // сервер
		if err != nil {
			log.Fatal(err)
		}
	}()

	//go func() {
	//	agent.CollectMetrics(config.PollInterval) // сбор метрик
	//}()
	//
	//go func() {
	//	err := agent.ClientMetric(config.ReportInterval) // отправка метрик
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}()
	//
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

	return http.ListenAndServe(config.AddrServer.String(), r)
}
