package util

import (
	"io"
	"net/http"
	neturl "net/url"
)

type HttpClient struct {
	DefaultHeaders map[string]string
}

func (c *HttpClient) Fetch(url string, fn func(io.Reader) error) error {
	return c.FetchWithHeader(url, make(map[string]string), fn)
}

func (c *HttpClient) FetchWithHeader(url string, headers map[string]string, fn func(io.Reader) error) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	u, err := neturl.Parse(url)
	if err != nil {
		return err
	}

	if u.Host == "twitter.com" {
		c.addHeaders(req, c.DefaultHeaders)
	}
	c.addHeaders(req, headers)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return fn(res.Body)
}

func (c *HttpClient) addHeaders(req *http.Request, headers map[string]string) {
	for k, v := range headers {
		req.Header.Set(k, v)
	}
}
