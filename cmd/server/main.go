package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/kornetvba/metrics-service/internal/config/logger"
	"github.com/kornetvba/metrics-service/internal/config/server"
	handlers "github.com/kornetvba/metrics-service/internal/handler"
	"github.com/kornetvba/metrics-service/internal/storage"
	"log"
	"net/http"
)

func main() {
	server.ParseFlagServer()
	err := logger.Initialization()

	if err != nil {
		log.Print(err)
	}

	err = run() // сервер
	if err != nil {
		log.Fatal(err)
	}

}

func run() error {
	log.Printf("serv is running %s", server.AddrServer.String())
	handler := handlers.NewMetricHandler(storage.NewMemStorage())
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Use(logger.LogMiddlewareGet)
		r.Get("/", handler.GetAllMetricsHTML)
		r.Get("/value/{type_metric}/{name_metric}", handler.MetricGet)
	})

	r.Route("/update", func(r chi.Router) {
		r.Use(logger.LogMiddlewarePost)
		r.Post("/{type_metric}/{name_metric}/{value_metric}", handler.MetricPost)

	})

	return http.ListenAndServe(server.AddrServer.String(), r)
}
