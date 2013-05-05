// http://opentsdb.net/docs/build/html/api_http/serializers/json.html#api-query

package tsdb

import (
	"time"
	"encoding/json"
	"regexp"
	"fmt"
	"github.com/davecgh/go-spew/spew"
)

type request struct {
	Start   tsTime  `json:"start"`
	End     tsTime  `json:"end,omitempty"` // Optional
	Padding bool    `json:"padding,omitempty"` // Optional
	Queries []query `json:"queries"`
}

type tsTime struct {
	time.Time
}

func (t *tsTime) UnmarshalJSON(inJSON []byte) error {
	spew.Dump("---tsTime UnmarshalJSON---")
	var raw interface{}
	err := json.Unmarshal(inJSON, &raw)
	if err != nil { panic(err) }

	switch raw.(type) {
	case float64:
		t.Time = time.Unix(int64(raw.(float64)), 0)
	case string:
		err = t.Parse(raw.(string))
	}
	if err != nil { return &json.InvalidUnmarshalError{} }
	return err
}

func (t *tsTime) Parse(timeIn string) error {
	switch {
	case !IsValidTime(timeIn):
		return fmt.Errorf("Invalid Time Value")
	case IsAbsoluteTime(timeIn):
		return t.fromAbsoluteTime(timeIn)
	case IsRelativeTime(timeIn):
		return t.fromRelativeTime(timeIn)
	case IsUnixTime(timeIn):
		return t.fromUnixTime(timeIn)
	}
	panic(fmt.Errorf("Invalid Time Value (Uncaught)"))
}

func IsValidTime(timeIn string) bool {
	if IsAbsoluteTime(timeIn) || IsRelativeTime(timeIn) || IsUnixTime(timeIn) {
		return true
	}
	return false
}

func IsAbsoluteTime(timeIn string) bool {
	// yyyy/MM/dd-HH:mm:ss
	// yyyy/MM/dd HH:mm:ss
	// yyyy/MM/dd-HH:mm
	// yyyy/MM/dd HH:mm
	// yyyy/MM/dd
	pattern := `^\d{4}\/\d{1,2}\/\d{1,2}`
	match, err := regexp.MatchString(pattern, timeIn)
	if err != nil { panic(err) }
	return match
}

func IsRelativeTime(timeIn string) bool {
	// 1{s,m,h,d,w,n,y}-ago
	pattern := `^\d+[smhdwmny]\-ago`
	match, err := regexp.MatchString(pattern, timeIn)
	if err != nil { panic(err) }
	return match
}

func IsUnixTime(timeIn string) bool {
	// 10-digit integer
	// 13-digit optional millisecond precision
	pattern := `^\d{10}|^\d{13}`
	match, err := regexp.MatchString(pattern, timeIn)
	if err != nil { panic(err) }
	return match
}

func (t *tsTime) fromAbsoluteTime(timeIn string) error {
	return nil
}

func (t *tsTime) fromRelativeTime(timeIn string) error {
	return nil
}

func (t *tsTime) fromUnixTime(timeIn string) error {
	return nil
}

type response struct {
	Metric          string
	Tags            map[string]string
	Aggregated_tags []string
	Dps             map[int64]int64
}

type query struct {
	Aggregator string
	Metric     string
	Rate       bool
	Downsample string `json:"downsample,omitempty"`
	Tags       map[string]string
}

func NewEmptyRequest() *request {
	newReq        := new(request)
	// newReq.Queries = make([]query, 0)
	return newReq
}

func NewEmptyResponse() *response {
	return new(response)
}
