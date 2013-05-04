package tsdb

import (
	"net"
	"net/http"
)

type client struct {
	http.Client
	host string
}

func NewClient(host string, port string) (*client) {
	client := new(client)
	client.host = net.JoinHostPort(host, port)
	return client
}
