package agent

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

func URLRequest(nameMc string, valMc interface{}) (string, error) {

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
		if globalMetrics == nil {
			time.Sleep(2 * time.Second)
			continue
		}
		gaugeMap := globalMetrics.ToMap()
		for k, v := range gaugeMap {
			url, err := URLRequest(k, v)

			if err != nil {
				continue
			}

			req, err := http.NewRequest(http.MethodPost, url, nil)

			if err != nil {
				continue
			}

			res, err := client.Do(req)

			if err != nil {
				continue
			}
			res.Body.Close()

		}
		log.Print("Metrics post ", time.Now())
		time.Sleep(timeDelay * time.Second)
	}

}
