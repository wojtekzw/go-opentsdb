// http://opentsdb.net/docs/build/html/api_http/serializers/json.html#api-put

package tsdb

import (
	"time"
)

type dataPoint struct {
	metric    string
	timestamp time.Time
	value     int64
	tags      map[string]string
}

func NewDataPoint() *dataPoint {
	return new(dataPoint)
}
