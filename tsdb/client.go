package tsdb

import (
	"net"
	"net/http"
	"bytes"
	"io/ioutil"
	"encoding/json"
)

type connection struct {
	http.Client
	host string
}

func NewConnection(host string, port string) (*connection, error) {
	connection := new(connection)
	connection.host = net.JoinHostPort(host, port)
	return connection, nil
}

func (c *connection) Query(req request) (*response, error) {
	APIURL := "http://"+c.host+"/api/query"
	reqJSON, err := json.Marshal(req)
	if err != nil { return &response{}, err }

	reqReader := bytes.NewReader(reqJSON)
	respHTTP, err := c.Post(APIURL, "application/json", reqReader)
	if err != nil { panic(err) }

	respJSON, err := ioutil.ReadAll(respHTTP.Body)
	if err != nil { return &response{}, err }

	resp := NewEmptyResponse()
	err = json.Unmarshal(respJSON, &resp)
	if err != nil { return &response{}, err }

	return resp, nil
}
