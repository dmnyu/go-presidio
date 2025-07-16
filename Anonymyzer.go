package go_presidio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AnonymizationResponse struct {
	Text  string `json:"text"`
	Items []Item `json:"items"`
}

type Item struct {
	Start      int    `json:"start"`
	End        int    `json:"end"`
	EntityType string `json:"entity_type"`
	Text       string `json:"text"`
	Operator   string `json:"operator"`
}

type AnonymizationRequest struct {
	Text            string          `json:"text"`
	AnalyzerResults AnalysisResults `json:"analyzer_results"`
	Anonymizers     struct {
		Default DefaultAnonymizer `json:"DEFAULT"`
	} `json:"anonymizers"`
}

type DefaultAnonymizer struct {
	AnonymizerType string `json:"type"`
	NewValue       string `json:"new_value"`
}

func (c *PresidioClient) AnonymizeText(ar *AnonymizationRequest) (*AnonymizationResponse, error) {
	jsonData, err := json.Marshal(ar)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal analysis request: %v", err)
	}

	resp, err := c.POST(c.URL+":5001/anonymize", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to analyze text: %s", resp.Status)
	}

	reader, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	results := AnonymizationResponse{}
	if err := json.Unmarshal(reader, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return &results, nil
}
