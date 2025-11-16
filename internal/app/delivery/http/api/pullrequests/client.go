package pullrequests

import (
	"net/http"
)

type Client struct {
	c       *http.Client
	baseUrl string
}

func NewClient(c *http.Client, baseUrl string) Client {
	return Client{c: c, baseUrl: baseUrl}
}
