// +build integration

// Run integration tests with:
// go test -tags=integration -host=127.0.0.01 -port=4242

package tsdb

import (
	"flag"
	"fmt"
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

	tags := &Tags{}
	tags.Set("host", "web01")
	tags.Set("dc", "lga")

	metric := &Metric{}
	metric.Set("sys.cpu.nice")

	value := &Value{}
	timestamp := &Time{}
	err := timestamp.Parse(fmt.Sprint(time.Now().Unix()))

	dataPoints := []DataPoint{
		DataPoint{
			Timestamp: timestamp,
			Metric:    metric,
			Value:     value,
			Tags:      tags,
		},
	}

	response, err := db.Put(dataPoints)
	assert.Nil(err)
	assert.NotNil(response)
}
