// http://opentsdb.net/docs/build/html/api_http/serializers/json.html#api-query

package tsdb

import (
	"time"
)

type request struct {
	start   time.Time
	end     time.Time `json:"end,omitempty"` // Optional
	padding bool      `json:"padding,omitempty"` // Optional
	queries []query
}

type response struct {
	metric          string
	tags            map[string]string
	aggregated_tags []string
	dps             map[time.Time]int64
}

type query struct {
	aggregator string
	metric     string
	rate       bool
	downsample string
	tags       map[string]string
}

func NewEmptyRequest() *request {
	newReq        := new(request)
	newReq.queries = make([]query, 0)
	return newReq
}

func NewEmptyResponse() *response {
	return new(response)
}
