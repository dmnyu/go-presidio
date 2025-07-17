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

type Anonymizer struct {
	AnonymizerType string `json:"type"`
	NewValue       string `json:"new_value"`
	Mask           string `json:"mask,omitempty"`
}

type AnonymizationRequest struct {
	Text            string                `json:"text"`
	AnalyzerResults AnalysisResults       `json:"analyzer_results"`
	Anonymizers     map[string]Anonymizer `json:"anonymizers"`
}

type AnonymizerAndLabel struct {
	Label      string     `json:"label"`
	Anonymizer Anonymizer `json:"anonymizer"`
}

func (ar *AnonymizationRequest) AddAnonymizer(al AnonymizerAndLabel) {
	ar.Anonymizers[al.Label] = al.Anonymizer
}

func (ar *AnonymizationRequest) AddAnonymizers(anonymizers map[string]Anonymizer) {
	for k, y := range anonymizers {
		ar.Anonymizers[k] = y
	}
}

func NewSimpleAnonymizer(value *string) Anonymizer {
	anonymizer := Anonymizer{
		AnonymizerType: "replace",
		NewValue:       "",
	}

	if value != nil {
		anonymizer.NewValue = *value
	} else {
		anonymizer.NewValue = "<REDACTED>"
	}

	return anonymizer
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
		return nil, fmt.Errorf("failed to anonymizer text: %s", resp.Status)
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
