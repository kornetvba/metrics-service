package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kornetvba/metrics-service/internal/agent"
	handlers "github.com/kornetvba/metrics-service/internal/handler"
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
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.GetAllMetricsHTML)
		r.Get("/value/{type_metric}/{name_metric}", handlers.MetricGet)

		r.Post("/update/{type_metric}/{name_metric}/{value_metric}", handlers.MetricPost)
	})

	return http.ListenAndServe(":8080", r)
}
