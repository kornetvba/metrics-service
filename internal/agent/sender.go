package agent

import (
	"errors"
	"fmt"
	"github.com/kornetvba/metrics-service/internal/config"
	"log"
	"net/http"
	"time"
)

func URLRequest(addr *config.NetAddr, nameMc string, valMc interface{}) (string, error) {

	switch valMc.(type) {
	case int64:
		return fmt.Sprintf("http://%s/update/%s/%s/%v", addr.String(), "counter", nameMc, valMc), nil
	case float64:
		return fmt.Sprintf("http://%s/update/%s/%s/%v", addr.String(), "gauge", nameMc, valMc), nil
	}
	return "", errors.New("type metric is not valid")

}

func ClientMetric(timeDelay time.Duration) error {
	log.Print("agent running!")
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	for {
		if globalMetrics == nil {
			time.Sleep(2 * time.Second)
			continue
		}
		gaugeMap := globalMetrics.ToMap()
		for k, v := range gaugeMap {
			url, err := URLRequest(config.Addr, k, v)

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
		time.Sleep(timeDelay)
	}

}
