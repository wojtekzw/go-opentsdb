package tsdb

import (
	"net"
	"net/http"
	"bytes"
	"io/ioutil"
	"encoding/json"
)

// connection represents a connection to an OpenTSDB server.
type connection struct {
	http.Client
	host string
}

// Dial returns a new connection type to an OpenTSDB host:port.
func Dial(host string, port string) (*connection, error) {
	connection := new(connection)
	connection.host = net.JoinHostPort(host, port)
	return connection, nil
}

// Query takes a TSDB Request and returns the resulting query Response.
func (c *connection) Query(req Request) (*Response, error) {
	APIURL := "http://"+c.host+"/api/query"
	reqJSON, err := json.Marshal(req)
	if err != nil { return &Response{}, err }

	reqReader := bytes.NewReader(reqJSON)
	respHTTP, err := c.Post(APIURL, "application/json", reqReader)
	if err != nil { panic(err) }

	respJSON, err := ioutil.ReadAll(respHTTP.Body)
	if err != nil { return &Response{}, err }

	resp := new(Response)
	err = json.Unmarshal(respJSON, &resp)
	if err != nil { return &Response{}, err }

	return resp, nil
}
