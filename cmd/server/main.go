package main

import (
	handlers "github.com/kornetvba/metrics-service/internal/handler"
	"log"
	"net/http"
)

func main() {

	err := run()

	if err != nil {
		log.Fatal(err)
	}

}

func run() error {
	log.Print("serv is running")
	mx := http.NewServeMux()
	//mx.HandleFunc("/update/gauge/", handlers.GaugePost)
	mx.HandleFunc("/update/{type_metric}/{name_metric}/{value_metric}", handlers.MetricPost)

	return http.ListenAndServe(":8080", mx)
}
