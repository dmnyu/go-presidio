package go_presidio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var analysisEndpoint = "analyze"

type AnalysisResults []AnalysisResult

type AnalysisResult struct {
	AnalysisExplanation string              `json:"analysis_explanation"`
	End                 int                 `json:"end"`
	EntityType          string              `json:"entity_type"`
	Score               float64             `json:"score"`
	Start               int                 `json:"start"`
	RecognitionMetadata RecognitionMetadata `json:"recognition_metadata"`
}

type RecognitionMetadata struct {
	RecognitionIdentifier string `json:"recognition_identifier"`
	RecognitionName       string `json:"recognition_name"`
}

type AnalysisRequest struct {
	Text           string  `json:"text"`
	Language       string  `json:"language"`
	ScoreThreshold float32 `json:"score_threshold,omitempty"`
}

func (ar AnalysisResult) String() string {
	return fmt.Sprintf("  Entity Type: %s\n  Score: %.2f\n  Start: %d\n  End: %d\n  Explanation: %s\n  Recognition Identifier: %s\n  Recognition Name: %s",
		ar.EntityType, ar.Score, ar.Start, ar.End, ar.AnalysisExplanation,
		ar.RecognitionMetadata.RecognitionIdentifier, ar.RecognitionMetadata.RecognitionName)
}

func (c *PresidioClient) AnalyzeText(ar *AnalysisRequest) (*AnalysisResults, error) {

	jsonData, err := json.Marshal(ar)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal analysis request: %v", err)
	}

	endpoint := fmt.Sprintf("%s:%d/%s", c.URL, c.AnalyzerPort, analysisEndpoint)
	resp, err := c.POST(endpoint, bytes.NewBuffer(jsonData))
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
	results := AnalysisResults{}

	if err := json.Unmarshal(reader, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return &results, nil
}

func (c *PresidioClient) AnalyzerHealth() (*string, error) {
	endpoint := fmt.Sprintf("%s:%d/health", c.URL, c.AnalyzerPort)
	resp, err := http.Get(endpoint)
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

func (c *PresidioClient) GetAnalyzerSupportedEntities() (*[]string, error) {
	endpoint := fmt.Sprintf("%s:%d/supportedentities?language=en", c.URL, c.AnalyzerPort)
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get supported entities: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get supported entities: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var entities []string
	if err := json.Unmarshal(body, &entities); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return &entities, nil
}

func (c *PresidioClient) GetAnalyzerRecognizers() (*[]string, error) {
	endpoint := fmt.Sprintf("%s:%d/recognizers?language=en", c.URL, c.AnalyzerPort)
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get recognizers: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get recognizers: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var recognizers []string
	if err := json.Unmarshal(body, &recognizers); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return &recognizers, nil
}
