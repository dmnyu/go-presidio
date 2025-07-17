package go_presidio

import (
	"testing"
)

var text = "Sample text to analyze:\n  My name is John Doe,\n  my email is john@example.com,\n  my phone number is 123-456-7890,\n  and my SSN is 111-59-5959."

func TestGoPresidio(t *testing.T) {
	var (
		client               *PresidioClient
		analysis_results     *AnalysisResults
		anonymizationRequest *AnonymizationRequest
		anonymizationResult  *AnonymizationResult
	)
	t.Run("test init client", func(t *testing.T) {
		client = NewPresidioClient("http://localhost")
		if client == nil {
			t.Error("Failed to initialize Presidio client")
		}
	})
	t.Run("test analyzer", func(t *testing.T) {
		var err error
		analysis_results, err = client.AnalyzeText(&AnalysisRequest{text, "en"})
		if err != nil {
			panic(err)
		}

		if len(*analysis_results) != 5 {
			t.Errorf("Expected 5 analysis results, got %d", len(*analysis_results))
		}

		for i, result := range *analysis_results {
			t.Logf("Analysis Result %d: %s (%0.2f)", i+1, result.EntityType, result.Score)
		}
	})

	t.Run("test anonymization", func(t *testing.T) {
		anonymizationRequest = &AnonymizationRequest{
			Text:            text,
			AnalyzerResults: *analysis_results,
			Anonymizers:     make(map[string]Anonymizer),
		}

		anonymizationRequest.AddAnonymizer(AnonymizerAndLabel{Label: "DEFAULT", Anonymizer: NewDefaultAnonymizer()})

		var err error
		anonymizationResult, err = client.AnonymizeText(anonymizationRequest)
		if err != nil {
			panic(err)
		}

		if len(anonymizationResult.Items) != 4 {
			t.Errorf("Expected 4 anonymization results, got %d", len(anonymizationResult.Items))
		}

		for i, item := range anonymizationResult.Items {
			t.Logf("Anonymization Result %d: %s", i+1, item.EntityType)
		}

	})
}
