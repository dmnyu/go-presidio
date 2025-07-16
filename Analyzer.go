package go_presidio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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
	Text     string `json:"text"`
	Language string `json:"language"`
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

	resp, err := c.POST(c.URL+":5002/analyze", bytes.NewBuffer(jsonData))
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
