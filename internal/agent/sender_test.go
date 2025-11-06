package agent

import (
	"fmt"
	"github.com/kornetvba/metrics-service/internal/config"
	"github.com/kornetvba/metrics-service/internal/config/server"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUrlRequest(t *testing.T) {

	tableTest := []struct {
		name   string
		nameMc string
		valMc  interface{}
		result string
		addr   *config.NetAddr
	}{
		{
			name:   "test1",
			nameMc: "testCounter",
			valMc:  int64(5),
			result: fmt.Sprintf("http://localhost:8080/update/%s/%s/%v", "counter", "testCounter", 5),
			addr:   server.AddrServer,
		},
		{
			name:   "test2",
			nameMc: "testGauge",
			valMc:  float64(5.4),
			result: fmt.Sprintf("http://localhost:8080/update/%s/%s/%v", "gauge", "testGauge", 5.4),
			addr:   server.AddrServer,
		},
		{
			name:   "test3",
			nameMc: "testGauge",
			valMc:  float64(5.0),
			result: fmt.Sprintf("http://localhost:8080/update/%s/%s/%v", "gauge", "testGauge", 5.0),
			addr:   server.AddrServer,
		},
		{
			name:   "test4",
			nameMc: "testError",
			valMc:  "5.0",
			result: "",
			addr:   server.AddrServer,
		},
		{
			name:   "test4",
			nameMc: "testflt",
			valMc:  55.5,
			result: fmt.Sprintf("http://localhost:8080/update/%s/%s/%v", "gauge", "testflt", 55.5),
			addr:   server.AddrServer,
		},
	}
	for _, tt := range tableTest {
		t.Run(tt.name, func(t *testing.T) {

			res, _ := URLRequest(tt.addr, tt.nameMc, tt.valMc)

			assert.Equal(t, tt.result, res)

		})
	}
}
