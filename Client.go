package go_presidio

import (
	"fmt"
	"io"
	"net/http"
)

type PresidioClient struct {
	URL    string
	client *http.Client
}

func NewPresidioClient(url string) *PresidioClient {
	return &PresidioClient{
		URL:    url,
		client: http.DefaultClient,
	}
}

func (c *PresidioClient) do(request *http.Request) (*http.Response, error) {

	resp, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *PresidioClient) POST(endpoint string, body io.Reader) (*http.Response, error) {

	request, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := c.do(request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to post to %s: %s", endpoint, resp.Status)
	}

	return resp, nil
}
