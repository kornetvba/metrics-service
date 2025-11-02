package agent

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUrlRequest(t *testing.T) {

	tableTest := []struct {
		name   string
		nameMc string
		valMc  interface{}
		result string
	}{
		{
			name:   "test1",
			nameMc: "testCounter",
			valMc:  int64(5),
			result: fmt.Sprintf("http://localhost:8080/update/%s/%s/%v", "counter", "testCounter", 5),
		},
		{
			name:   "test2",
			nameMc: "testGauge",
			valMc:  float64(5.4),
			result: fmt.Sprintf("http://localhost:8080/update/%s/%s/%v", "gauge", "testGauge", 5.4),
		},
		{
			name:   "test3",
			nameMc: "testGauge",
			valMc:  float64(5.0),
			result: fmt.Sprintf("http://localhost:8080/update/%s/%s/%v", "gauge", "testGauge", 5.0),
		},
		{
			name:   "test4",
			nameMc: "testError",
			valMc:  "5.0",
			result: "",
		},
		{
			name:   "test4",
			nameMc: "testflt",
			valMc:  55.5,
			result: fmt.Sprintf("http://localhost:8080/update/%s/%s/%v", "gauge", "testflt", 55.5),
		},
	}
	for _, tt := range tableTest {
		t.Run(tt.name, func(t *testing.T) {

			res, _ := UrlRequest(tt.nameMc, tt.valMc)

			assert.Equal(t, tt.result, res)

		})
	}
}
