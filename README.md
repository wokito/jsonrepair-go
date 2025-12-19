# jsonrepair-go

[![Go Reference](https://pkg.go.dev/badge/github.com/wokito/jsonrepair-go.svg)](https://pkg.go.dev/github.com/wokito/jsonrepair-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/wokito/jsonrepair-go)](https://goreportcard.com/report/github.com/wokito/jsonrepair-go)
[![License: ISC](https://img.shields.io/badge/License-ISC-blue.svg)](https://opensource.org/licenses/ISC)

Repair invalid JSON documents.

This is a Go port of the original [jsonrepair](https://github.com/josdejong/jsonrepair) library by Jos de Jong.

Read the background article ["How to fix JSON and validate it with ease"](https://jsoneditoronline.org/indepth/parse/fix-json/)

The following issues can be fixed:

- Add missing quotes around keys
- Add missing escape characters
- Add missing commas
- Add missing closing brackets
- Repair truncated JSON
- Replace single quotes with double quotes
- Replace special quote characters like `"..."` with regular double quotes
- Replace special white space characters with regular spaces
- Replace Python constants `None`, `True`, and `False` with `null`, `true`, and `false`
- Strip trailing commas
- Strip comments like `/* ... */` and `// ...`
- Strip fenced code blocks like ` ```json ` and ` ``` `
- Strip ellipsis in arrays and objects like `[1, 2, 3, ...]`
- Strip JSONP notation like `callback({ ... })`
- Strip escape characters from an escaped string like `{\"stringified\": \"content\"}`
- Strip MongoDB data types like `NumberLong(2)` and `ISODate("2012-12-19T06:01:17.171Z")`
- Concatenate strings like `"long text" + "more text on next line"`
- Turn newline delimited JSON into a valid JSON array, for example:
    ```
    { "id": 1, "name": "John" }
    { "id": 2, "name": "Sarah" }
    ```

## Install

```bash
go get github.com/wokito/jsonrepair-go
```

## Use

```go
package main

import (
    "fmt"
    "log"

    jsonrepair "github.com/wokito/jsonrepair-go"
)

func main() {
    // The following is invalid JSON: it consists of JSON contents copied from
    // a JavaScript code base, where the keys are missing double quotes,
    // and strings are using single quotes:
    json := "{name: 'John'}"

    repaired, err := jsonrepair.JSONRepair(json)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(repaired) // {"name": "John"}
}
```

### MustJSONRepair

Use `MustJSONRepair` when you want to panic on error:

```go
// Will panic if repair fails
result := jsonrepair.MustJSONRepair("[1, 2, 3,]")
fmt.Println(result) // [1, 2, 3]
```

## API

### JSONRepair

```go
func JSONRepair(text string) (string, error)
```

Repairs a string containing an invalid JSON document. Returns the repaired JSON string, or an error when an issue is encountered which could not be solved.

### MustJSONRepair

```go
func MustJSONRepair(text string) string
```

Same as `JSONRepair`, but panics instead of returning an error.

## Examples

### Fix missing quotes on keys

```go
result, _ := jsonrepair.JSONRepair("{name: 'John'}")
// Output: {"name": "John"}
```

### Fix single quotes

```go
result, _ := jsonrepair.JSONRepair("{'name': 'Alice'}")
// Output: {"name": "Alice"}
```

### Fix trailing commas

```go
result, _ := jsonrepair.JSONRepair("[1, 2, 3,]")
// Output: [1, 2, 3]
```

### Strip comments

```go
result, _ := jsonrepair.JSONRepair(`{
    "name": "Bob", // comment
    "active": true /* another comment */
}`)
// Output: {"name": "Bob", "active": true}
```

### Repair truncated JSON

```go
result, _ := jsonrepair.JSONRepair(`{"message": "hello`)
// Output: {"message": "hello"}
```

### Fix Python constants

```go
result, _ := jsonrepair.JSONRepair(`{"enabled": True, "data": None}`)
// Output: {"enabled": true, "data": null}
```

### Concatenate strings

```go
result, _ := jsonrepair.JSONRepair(`"hello" + " world"`)
// Output: "hello world"
```

### Strip MongoDB data types

```go
result, _ := jsonrepair.JSONRepair(`{"_id": ObjectId("123"), "count": NumberLong("456")}`)
// Output: {"_id": "123", "count": "456"}
```

### Repair newline delimited JSON (NDJSON)

```go
result, _ := jsonrepair.JSONRepair(`{"id":1}
{"id":2}`)
// Output: [{"id":1},{"id":2}]
```

## Differences from TypeScript version

This is a Go port of [josdejong/jsonrepair](https://github.com/josdejong/jsonrepair). The functionality is mostly identical, with the following known differences:

| Feature | TypeScript | Go |
|---------|------------|-----|
| MongoDB Timestamp | ✅ Supported | ❌ Not supported |
| Triple quotes `'''` | ✅ Repaired | ❌ Returns error |
| Streaming API | ✅ Available | ❌ Not available |

## Testing

```bash
go test -v
```

The test suite covers 78 test cases, all passing. The tests are aligned with the original TypeScript test suite.

## License

Released under the [ISC license](https://opensource.org/licenses/ISC).

This project is a Go port of [jsonrepair](https://github.com/josdejong/jsonrepair), originally created by Jos de Jong.

## Acknowledgments

- [Jos de Jong](https://github.com/josdejong) - Author of the original jsonrepair TypeScript library
