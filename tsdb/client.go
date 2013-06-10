package tsdb

import (
	"net/http"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"strconv"
)

// Server represents an OpenTSDB server
type Server struct {
	Host string
	Port uint
}

// TSDB represents an OpenTSDB database serviced by one or more
// Servers.
type TSDB struct {
	Servers []Server
}

// Query takes a TSDB Request and returns the resulting query Response.
func (t *TSDB) Query(req Request) (*Response, error) {
	// TODO: Handle multiple Servers
	host := t.Servers[0].Host+":"+strconv.Itoa(int(t.Servers[0].Port))
	APIURL := "http://"+host+"/api/query"
	reqJSON, err := json.Marshal(req)
	if err != nil { return &Response{}, err }

	reqReader := bytes.NewReader(reqJSON)
	respHTTP, err := http.Post(APIURL, "application/json", reqReader)
	if err != nil { panic(err) }

	respJSON, err := ioutil.ReadAll(respHTTP.Body)
	if err != nil { return &Response{}, err }

	resp := new(Response)
	err = json.Unmarshal(respJSON, &resp)
	if err != nil { return &Response{}, err }

	return resp, nil
}
