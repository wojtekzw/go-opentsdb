// http://opentsdb.net/docs/build/html/api_http/serializers/json.html#api-query

package tsdb

import (
	"time"
	"encoding/json"
	"regexp"
	"fmt"
	"strconv"
	// "github.com/davecgh/go-spew/spew"
)

// Request represents the information needed to query a TSDB for timeseries data.
type Request struct {
	Start   *TSTime  `json:"start"`
	End     *TSTime  `json:"end,omitempty"` // Optional
	Padding bool     `json:"padding,omitempty"` // Optional
	Queries []query  `json:"queries"`
}

/*
TSTime represents a timeseries time value.

Valid formats for TSTime are:
	Relative (see: http://opentsdb.net/docs/build/html/user_guide/query/index.html#relative)
	Unix     (see: http://opentsdb.net/docs/build/html/user_guide/query/index.html#absolute-unix-time)
	Absolute (see: http://opentsdb.net/docs/build/html/user_guide/query/index.html#absolute-formatted-time)
*/
type TSTime struct {
	time.Time
	format string
	string
}

// UnmarshalJSON implements json.Unmarshaler for consistant conversion from JSON.
func (t *TSTime) UnmarshalJSON(inJSON []byte) error {
	var raw interface{}
	err := json.Unmarshal(inJSON, &raw)
	if err != nil { panic(err) }

	switch raw.(type) {
	case float64:
		err = t.Parse(strconv.FormatInt(int64(raw.(float64)), 10))
	case string:
		err = t.Parse(raw.(string))
	}
	return err
}

// MarshalJSON implements json.Marshaler for consistant conversion to JSON.
func (t TSTime) MarshalJSON() ([]byte, error) {
	switch t.format {
	case "":         return nil, nil
	case "Unix":     return json.Marshal(t.Unix())
	case "Absolute": return json.Marshal(t.AbsoluteTime())
	case "Relative": return json.Marshal(t.string)
	}
	return json.Marshal(t.Unix())
}

// Parse takes a string, verifies that it is a valid TSTime, and if so sets t to that time.
// If the input string is not a valid TSTime then t is unchanged.
func (t *TSTime) Parse(timeIn string) error {
	switch {
	case !IsValidTime(timeIn):   return fmt.Errorf("Invalid Time Value")
	case IsAbsoluteTime(timeIn): return t.fromAbsoluteTime(timeIn)
	case IsRelativeTime(timeIn): return t.fromRelativeTime(timeIn)
	case IsUnixTime(timeIn):     return t.fromUnixTime(timeIn)
	}
	return fmt.Errorf("Invalid Time Value (Uncaught)")
}

// IsValidTime verifies that a string can be converted to a TSTime.
func IsValidTime(timeIn string) bool {
	if IsAbsoluteTime(timeIn) || IsRelativeTime(timeIn) || IsUnixTime(timeIn) {
		return true
	}
	return false
}

/*
IsAbsoluteTime verifies if a string is a valid Absolute format time.
Valid formats are:
	yyyy/MM/dd-HH:mm:ss
	yyyy/MM/dd HH:mm:ss
	yyyy/MM/dd-HH:mm
	yyyy/MM/dd HH:mm
	yyyy/MM/dd
*/
func IsAbsoluteTime(timeIn string) bool {
	pattern := `^\d{4}\/\d{1,2}\/\d{1,2}`
	match, err := regexp.MatchString(pattern, timeIn)
	if err != nil { panic(err) }
	return match
}
/*
IsAbsoluteTime verifies if a string is a valid Relative format time.

Valid formats are:
	[0-9]*{s,m,h,d,w,n,y}-ago
*/
func IsRelativeTime(timeIn string) bool {
	pattern := `^\d+[smhdwmny]\-ago`
	match, err := regexp.MatchString(pattern, timeIn)
	if err != nil { panic(err) }
	return match
}

/*
IsAbsoluteTime verifies if a string is a valid Unix format time.

Valid formats are:
	10-digit integer
	13-digit optional millisecond precision
*/
func IsUnixTime(timeIn string) bool {
	pattern := `^\d{10}|^\d{13}`
	match, err := regexp.MatchString(pattern, timeIn)
	if err != nil { panic(err) }
	return match
}

// fromAbsoluteTime parses the provided timeIn string and if possible
// assigns the time to TSTime t.
func (t *TSTime) fromAbsoluteTime(timeIn string) (err error) {
	t.format = "Absolute"
	t.string = timeIn
	t.Time, err = time.Parse("2006/01/02-15:04:05", timeIn)
	if err == nil { return }
	t.Time, err = time.Parse("2006/01/02 15:04:05", timeIn)
	if err == nil { return }
	t.Time, err = time.Parse("2006/01/02-15:04", timeIn)
	if err == nil { return }
	t.Time, err = time.Parse("2006/01/02 15:04", timeIn)
	if err == nil { return }
	t.Time, err = time.Parse("2006/01/02-15", timeIn)
	if err == nil { return }
	t.Time, err = time.Parse("2006/01/02 15", timeIn)
	if err == nil { return }
	t.Time, err = time.Parse("2006/01/02", timeIn)
	return
}

// AbsoluteTime returns the a string version of a TSTime in Absolute format.
func (t *TSTime) AbsoluteTime() (string) {
	switch {
	case t.Second() > 0: return t.Time.Format("2006/01/02-15:04:05")
	case t.Minute() > 0: return t.Time.Format("2006/01/02-15:04")
	case t.Hour()   > 0: return t.Time.Format("2006/01/02-15")
	}
	return t.Time.Format("2006/01/02")
}

// fromRelativeTime parses the provided timeIn string and if possible
// assigns the time to TSTime t.
func (t *TSTime) fromRelativeTime(timeIn string) error {
	t.format = "Relative"
	t.string = timeIn
	return nil
}

// RelativeTime returns the a string version of a TSTime in Relative format.
func (t *TSTime) RelativeTime() (string) {
	return t.Time.Format("2006/01/02-15:04:05")
}

// fromUnixTime parses the provided timeIn string and if possible
// assigns the time to TSTime t.
func (t *TSTime) fromUnixTime(timeIn string) (err error) {
	t.format = "Unix"
	t.string = timeIn
	var timeInInt64 int64
	timeInInt64, err = strconv.ParseInt(timeIn, 10, 64)
	if err != nil { return err }
	t.Time = time.Unix(timeInInt64, 0)
	return nil
}

/*
Response respresents a full, valid (non-error) Response from an OpenTSDB query 
made up of zero or more result types.

See: http://opentsdb.net/docs/build/html/api_http/serializers/json.html#Response
*/
type Response []result

// result is a single timeseries or aggregate Response from an OpenTSDB query.
// See: http://opentsdb.net/docs/build/html/api_http/serializers/json.html#Response
type result struct {
	Metric          string             `json:"metric"`
	Tags            map[string]string  `json:"tags"`
	AggregatedTags  []string           `json:"aggregateTags"`
	Dps             map[string]tsValue `json:"dps"`
}

// tsValue represents a timeseries datapoint (timestamp + value).
type tsValue float64

// UnmarshalJSON implements json.Unmarshaler for consistant conversion from JSON.
func (v *tsValue) UnmarshalJSON(inJSON []byte) error {
	var raw interface{}
	err := json.Unmarshal(inJSON, &raw)
	if err != nil { panic(err) }

	switch raw.(type) {
	case float64, int64:
		*v = tsValue(raw.(float64))
		return nil
	case string:
		i, err := strconv.ParseFloat(raw.(string), 64)
		if err != nil { return err }
		*v = tsValue(i)
		return nil
	}
	return &json.UnmarshalTypeError{}
}

// query represents the information needed for a single query to OpenTSDB sans
// time interval.
type query struct {
	Aggregator string `json:"aggregator"`
	Metric     string `json:"metric"`
	Rate       bool   `json:"rate"`
	Downsample string `json:"downsample,omitempty"`
	Tags       map[string]string `json:"tags"`
}

// NewEmptyRequest returns an empty Request type.
func NewEmptyRequest() *Request {
	newReq        := new(Request)
	// newReq.Queries = make([]query, 0)
	return newReq
}

// NewEmptyResponse returns an empty Response type.
func NewEmptyResponse() *Response {
	resp := make(Response, 0)
	return &resp
	// return new(Response)
}
