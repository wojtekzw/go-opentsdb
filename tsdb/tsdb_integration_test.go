// +build integration

// Run integration tests with:
// go test -tags=integration -host=127.0.0.01 -port=4242

package tsdb

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var host = flag.String("host", "127.0.0.1", "host")
var port = flag.Uint("port", 4242, "port")

func TestPut(t *testing.T) {
	assert := assert.New(t)

	db := &TSDB{
		[]Server{
			Server{
				Host: *host,
				Port: *port,
			},
		},
	}

	tags := make(map[string]string)
	tags["host"] = "web01"
	tags["dc"] = "lga"

	dataPoints := []DataPoint{
		DataPoint{
			Timestamp: &Time{time.Now().Add(-30 * time.Second), "Unix", ""},
			Metric:    &Metric{"sys.cpu.nice"},
			Value:     &Value{float64: 80},
			Tags:      &Tags{tags},
		},
		DataPoint{
			Timestamp: &Time{time.Now().Add(-20 * time.Second), "Unix", ""},
			Metric:    &Metric{"sys.cpu.nice"},
			Value:     &Value{float64: 90},
			Tags:      &Tags{tags},
		},
		DataPoint{
			Timestamp: &Time{time.Now().Add(-10 * time.Second), "Unix", ""},
			Metric:    &Metric{"sys.cpu.nice"},
			Value:     &Value{float64: 100},
			Tags:      &Tags{tags},
		},
	}

	response, err := db.Put(dataPoints)
	assert.Nil(err)
	assert.NotNil(response)
}
