package main

import (
	"fmt"
	jsonrepair "github.com/wokito/jsonrepair-go"
)

func main() {
	tests := []struct {
		input string
		name  string
	}{
		{"```json\n{\"name\": \"John\"}\n```", "markdown with language"},
		{"```\n{\"name\": \"John\"}\n```", "markdown without language"},
		{"{\"name\": \"John\"}", "plain JSON"},
	}

	for _, tt := range tests {
		result, err := jsonrepair.JSONRepair(tt.input)
		if err != nil {
			fmt.Printf("%s: ERROR - %v\n", tt.name, err)
		} else {
			fmt.Printf("%s:\n  Input:  %q\n  Result: %q\n\n", tt.name, tt.input, result)
		}
	}
}
