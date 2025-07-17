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
		var err error
		client, err = NewPresidioClient("go-presidio-config.json")
		if err != nil {
			t.Error("Failed to initialize Presidio client")
		}

		if client.URL != "http://localhost" {
			t.Errorf("wanted http://localhost got %s", client.URL)
		}

		if client.AnalyzerPort != 5002 {
			t.Errorf("wanted 5002 got %d", client.AnalyzerPort)
		}
	})

	t.Run("test analyzer health", func(t *testing.T) {
		var err error
		health, err := client.AnalyzerHealth()
		if err != nil {
			t.Error("Failed to get analyzer health")
		}
		t.Log(*health)
	})

	t.Run("test analyzer", func(t *testing.T) {
		var err error
		analysis_results, err = client.AnalyzeText(&AnalysisRequest{Text: text, Language: "en"})
		if err != nil {
			panic(err)
		}

		if len(*analysis_results) != 5 {
			t.Errorf("Expected 5 analysis results, got %d", len(*analysis_results))
		}
	})

	t.Run("test analyzer scoring", func(t *testing.T) {
		var err error
		analysis_results, err = client.AnalyzeText(&AnalysisRequest{Text: text, Language: "en", ScoreThreshold: 0.75})
		if err != nil {
			panic(err)
		}

		for _, result := range *analysis_results {
			if result.Score < 0.75 {
				t.Errorf("Expected score >= 0.75, got %.2f for entity type %s", result.Score, result.EntityType)
			}
		}

	})

	t.Run("test analyzer supported entities", func(t *testing.T) {
		supportedEntities, err := client.GetAnalyzerSupportedEntities()
		if err != nil {
			t.Error("Failed to get supported entities")
		}

		if len(*supportedEntities) <= 1 {
			t.Errorf("Expected more than 1 supported entity, got %d", len(*supportedEntities))
		}
	})

	t.Run("test analyzer recognizers", func(t *testing.T) {
		recognizers, err := client.GetAnalyzerRecognizers()
		if err != nil {
			t.Error("Failed to get recognizers")
		}
		if len(*recognizers) <= 1 {
			t.Errorf("Expected more than 1 recognizer, got %d", len(*recognizers))
		}
	})

	t.Run("test anonymizer health", func(t *testing.T) {
		var err error
		health, err := client.AnonymizerHealth()
		if err != nil {
			t.Error("Failed to get anonymizer health")
		}
		t.Log(*health)
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

	})
}
