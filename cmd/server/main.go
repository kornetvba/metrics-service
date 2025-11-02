package main

import (
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

	// Главная горутина ждет бесконечно
	select {}
}

func run() error {
	log.Print("serv is running")
	mx := http.NewServeMux()
	//mx.HandleFunc("/update/gauge/", handlers.GaugePost)
	mx.HandleFunc("/update/{type_metric}/{name_metric}/{value_metric}", handlers.MetricPost)

	return http.ListenAndServe(":8080", mx)
}
