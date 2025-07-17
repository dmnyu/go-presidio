package main

import (
	"fmt"

	gp "github.com/dmnyu/go-presidio"
)

var text = "Sample text to analyze: My name is John Doe, my email is john@example.com"

func main() {

	client := gp.NewPresidioClient("http://localhost")
	ar := &gp.AnalysisRequest{
		Text:     text,
		Language: "en",
	}

	analyzer_results, err := client.AnalyzeText(ar)
	if err != nil {
		panic(err)
	}

	anonymizerRequest := &gp.AnonymizationRequest{
		Text:            text,
		AnalyzerResults: *analyzer_results,
		Anonymizers:     make(map[string]gp.Anonymizer),
	}

	anonymizerRequest.AddAnonymizer(gp.AnonymizerAndLabel{Label: "PERSON", Anonymizer: gp.NewSimpleAnonymizer(nil)})
	anonymization_result, err := client.AnonymizeText(anonymizerRequest)
	if err != nil {
		panic(err)
	}

	fmt.Println("Input Text:", text)
	fmt.Println("Output Text:", anonymization_result.Text)
}
