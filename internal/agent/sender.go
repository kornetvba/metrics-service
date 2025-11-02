package agent

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

func UrlRequest(nameMc string, valMc interface{}) (string, error) {

	switch valMc.(type) {
	case int64:
		return fmt.Sprintf("http://localhost:8080/update/%s/%s/%v", "counter", nameMc, valMc), nil
	case float64:
		return fmt.Sprintf("http://localhost:8080/update/%s/%s/%v", "gauge", nameMc, valMc), nil
	}
	return "", errors.New("type metric is not valid")

}

func ClientMetric(timeDelay time.Duration) error {
	log.Print("agent running!")
	client := http.Client{}

	for {

		gaugeMap := globalMetrics.ToMap()
		for k, v := range gaugeMap {
			url, err := UrlRequest(k, v)

			if err != nil {
				return err
			}

			req, err := http.NewRequest(http.MethodPost, url, nil)

			if err != nil {
				return err
			}

			_, err = client.Do(req)

			if err != nil {
				return err
			}

		}
		log.Print("Metrics post", time.Now())
		time.Sleep(timeDelay * time.Second)
	}

}
