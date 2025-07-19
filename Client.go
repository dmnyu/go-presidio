package go_presidio

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type PresidioClient struct {
	URL            string
	client         *http.Client
	AnalyzerPort   int
	AnonymizerPort int
}

type ClientConfig struct {
	Host           string `json:"host"`
	AnalyzerPort   int    `json:"analyzer_port"`
	AnonymizerPort int    `json:"anonymizer_port"`
}

func NewPresidioClient(configPath string) (*PresidioClient, error) {

	if _, err := os.Stat(configPath); err != nil {
		fmt.Println("CONFIGPATH:", configPath)
		return nil, err
	}

	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := ClientConfig{}
	if err := json.Unmarshal(configBytes, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	fmt.Println(config)

	return &PresidioClient{
		URL:            config.Host,
		client:         http.DefaultClient,
		AnalyzerPort:   config.AnalyzerPort,
		AnonymizerPort: config.AnonymizerPort,
	}, nil
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

func (c *PresidioClient) GET(endpoint string) (*http.Response, error) {
	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get from %s: %s", endpoint, resp.Status)
	}

	return resp, nil
}
