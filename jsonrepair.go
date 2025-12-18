// Package jsonrepair provides functionality to repair invalid JSON documents.
// This is a Go port of the original jsonrepair project by Jos de Jong.
// Original project: https://github.com/josdejong/jsonrepair
// Licensed under the ISC License
package jsonrepair

// JSONRepair repairs a string containing an invalid JSON document.
// It converts JavaScript notation into JSON notation and fixes various issues.
//
// Example:
//
//	input := "{name: 'John'}"
//	repaired, err := JSONRepair(input)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(repaired) // {"name": "John"}
//
// The following issues can be fixed:
//   - Add missing quotes around keys
//   - Add missing escape characters
//   - Add missing commas
//   - Add missing closing brackets
//   - Repair truncated JSON
//   - Replace single quotes with double quotes
//   - Replace special quote characters
//   - Replace special whitespace characters
//   - Replace Python constants (None, True, False)
//   - Strip trailing commas
//   - Strip comments (/* */ and //)
//   - Strip fenced code blocks
//   - Strip ellipsis in arrays and objects
//   - Strip JSONP notation
//   - Strip MongoDB data types
//   - Concatenate strings
//   - Turn newline delimited JSON into a valid JSON array
func JSONRepair(text string) (string, error) {
	parser := NewParser(text)
	return parser.Parse()
}

// MustJSONRepair repairs a string containing an invalid JSON document.
// It panics if the JSON cannot be repaired.
func MustJSONRepair(text string) string {
	result, err := JSONRepair(text)
	if err != nil {
		panic(err)
	}
	return result
}
