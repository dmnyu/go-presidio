package go_presidio

import (
	"encoding/json"
	"fmt"
)

func PrintFormattedJson(a any) error {
	b, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	return nil

}
