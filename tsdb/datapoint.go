// http://opentsdb.net/docs/build/html/api_http/serializers/json.html#api-put

package tsdb

import (
	"time"
)

// DataPoint represents a single data point good for recording in OpenTSDB.
// See: http://opentsdb.net/docs/build/html/api_http/serializers/json.html#api-put
type DataPoint struct {
	metric    string
	timestamp time.Time
	value     int64
	tags      map[string]string
}
