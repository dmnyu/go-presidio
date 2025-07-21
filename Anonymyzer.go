package go_presidio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Anonymizer struct {
	AnonymizerType string `json:"type"`
	NewValue       string `json:"new_value"`
	MaskingChar    string `json:"masking_char,omitempty"`
	CharsToMask    int    `json:"chars_to_mask,omitempty"`
	FromEnd        bool   `json:"from_end,omitempty"`
	Key            string `json:"key,omitempty"`
}

type AnonymizerAndLabel struct {
	Label      string     `json:"label"`
	Anonymizer Anonymizer `json:"anonymizer"`
}

type AnonymizationResult struct {
	Text  string `json:"text"`
	Items []Item `json:"items"`
}

type AnonymizationRequest struct {
	Text            string                `json:"text"`
	AnalyzerResults AnalysisResults       `json:"analyzer_results"`
	Anonymizers     map[string]Anonymizer `json:"anonymizers,omitempty"`
}

type Item struct {
	Start      int    `json:"start"`
	End        int    `json:"end"`
	EntityType string `json:"entity_type"`
	Text       string `json:"text"`
	Operator   string `json:"operator"`
}

func (ar *AnonymizationRequest) AddAnonymizer(al AnonymizerAndLabel) {
	ar.Anonymizers[al.Label] = al.Anonymizer
}

func (ar *AnonymizationRequest) AddAnonymizers(anonymizers map[string]Anonymizer) {
	for k, y := range anonymizers {
		ar.Anonymizers[k] = y
	}
}

func NewDefaultAnonymizer() Anonymizer {
	return NewReplaceAnonymizer(nil)
}

func NewReplaceAnonymizer(value *string) Anonymizer {
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

func NewMaskAnonymizer(maskingChar string, charsToMask int, fromEnd bool) Anonymizer {
	anonymizer := Anonymizer{
		AnonymizerType: "mask",
		MaskingChar:    "*",
		CharsToMask:    charsToMask,
		FromEnd:        fromEnd,
	}

	return anonymizer
}

func NewHashAnonymizer() Anonymizer {
	return Anonymizer{
		AnonymizerType: "hash",
	}
}

func (c *PresidioClient) AnonymizeText(ar *AnonymizationRequest) (*AnonymizationResult, error) {
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

	results := AnonymizationResult{}
	if err := json.Unmarshal(reader, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return &results, nil
}

func (c *PresidioClient) DeAnonymizeText() (*string, error) {
	return nil, nil
}

func (c *PresidioClient) AnonymizerHealth() (*string, error) {
	endpoint := fmt.Sprintf("%s:%d/health", c.URL, c.AnonymizerPort)
	resp, err := c.GET(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to check health: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("health check failed: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read health response: %v", err)
	}

	bodyString := string(body)
	var out *string = &bodyString
	return out, nil
}

func (c *PresidioClient) GetAnonymizers() (*[]string, error) {
	endpoint := fmt.Sprintf("%s:%d/anonymizers", c.URL, c.AnonymizerPort)
	resp, err := c.GET(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get anonymizers: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get anonymizers: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read anonymizers response: %v", err)
	}

	var anonymizers []string
	if err := json.Unmarshal(body, &anonymizers); err != nil {
		return nil, fmt.Errorf("failed to unmarshal anonymizers response: %v", err)
	}

	return &anonymizers, nil

}
