package jsonrepair

import (
	"encoding/json"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Parse parses and repairs the JSON text
func (p *Parser) Parse() (string, error) {
	// Parse optional markdown code block at the start
	p.parseMarkdownCodeBlock([]string{"```", "[```", "{```"})

	// Parse the main value
	processed := p.parseValue()
	if !processed {
		return "", p.throwUnexpectedEnd()
	}

	// Parse optional markdown code block at the end (BEFORE checking for NDJSON)
	p.parseMarkdownCodeBlock([]string{"```", "```]", "```}"})

	// Handle trailing comma
	p.parseWhitespaceAndSkipComments(true)
	processedComma := p.parseCharacter(',')
	if processedComma {
		p.parseWhitespaceAndSkipComments(true)
	}

	// Check for newline delimited JSON
	if p.i < len(p.text) && isStartOfValue(p.text, p.i) && endsWithCommaOrNewline(p.output.String()) {
		if !processedComma {
			// Repair missing comma
			output := p.output.String()
			p.output.Reset()
			p.output.WriteString(insertBeforeLastWhitespace(output, ","))
		}
		p.parseNewlineDelimitedJSON()
	} else if processedComma {
		// Remove trailing comma
		output := p.output.String()
		p.output.Reset()
		p.output.WriteString(stripLastOccurrence(output, ",", false))
	}

	// Repair redundant end quotes
	for p.i < len(p.text) {
		r, _ := getCharAt(p.text, p.i)
		if r == '}' || r == ']' {
			p.i++
			p.parseWhitespaceAndSkipComments(true)
		} else {
			break
		}
	}

	// Check if we've reached the end
	if p.i >= len(p.text) {
		return p.output.String(), nil
	}

	return "", p.throwUnexpectedCharacter()
}

// parseValue parses any JSON value
func (p *Parser) parseValue() bool {
	p.parseWhitespaceAndSkipComments(true)
	processed := p.parseObject() ||
		p.parseArray() ||
		p.parseString(false, -1) ||
		p.parseNumber() ||
		p.parseKeywords() ||
		p.parseUnquotedString(false) ||
		p.parseRegex()
	p.parseWhitespaceAndSkipComments(true)
	return processed
}

// parseWhitespaceAndSkipComments parses whitespace and skips comments
func (p *Parser) parseWhitespaceAndSkipComments(skipNewline bool) bool {
	start := p.i
	changed := p.parseWhitespace(skipNewline)
	for {
		changed = p.parseComment()
		if changed {
			changed = p.parseWhitespace(skipNewline)
		}
		if !changed {
			break
		}
	}
	return p.i > start
}

// parseWhitespace parses whitespace characters
func (p *Parser) parseWhitespace(skipNewline bool) bool {
	var whitespace strings.Builder

	for p.i < len(p.text) {
		if skipNewline && isWhitespace(p.text, p.i) {
			r, size := utf8.DecodeRuneInString(p.text[p.i:])
			whitespace.WriteRune(r)
			p.i += size
		} else if !skipNewline && isWhitespaceExceptNewline(p.text, p.i) {
			r, size := utf8.DecodeRuneInString(p.text[p.i:])
			whitespace.WriteRune(r)
			p.i += size
		} else if isSpecialWhitespace(p.text, p.i) {
			// Repair special whitespace
			whitespace.WriteRune(' ')
			r, size := utf8.DecodeRuneInString(p.text[p.i:])
			_ = r
			p.i += size
		} else {
			break
		}
	}

	if whitespace.Len() > 0 {
		p.output.WriteString(whitespace.String())
		return true
	}
	return false
}

// parseComment parses and skips comments
func (p *Parser) parseComment() bool {
	// Block comment /* ... */
	if p.i < len(p.text)-1 && p.text[p.i] == '/' && p.text[p.i+1] == '*' {
		for p.i < len(p.text) && !atEndOfBlockComment(p.text, p.i) {
			p.i++
		}
		p.i += 2
		return true
	}

	// Line comment // ...
	if p.i < len(p.text)-1 && p.text[p.i] == '/' && p.text[p.i+1] == '/' {
		for p.i < len(p.text) && p.text[p.i] != '\n' {
			p.i++
		}
		return true
	}

	return false
}

// parseMarkdownCodeBlock parses and skips markdown code blocks
func (p *Parser) parseMarkdownCodeBlock(blocks []string) bool {
	if p.skipMarkdownCodeBlock(blocks) {
		// Check for optional language specifier
		if p.i < len(p.text) {
			r, _ := getCharAt(p.text, p.i)
			if isFunctionNameCharStart(r) {
				for p.i < len(p.text) {
					r, _ := getCharAt(p.text, p.i)
					if isFunctionNameChar(r) {
						p.i++
					} else {
						break
					}
				}
			}
		}
		// Skip whitespace and comments after the code block
		p.parseWhitespaceAndSkipComments(true)
		return true
	}
	return false
}

// skipMarkdownCodeBlock skips markdown code block markers
func (p *Parser) skipMarkdownCodeBlock(blocks []string) bool {
	p.parseWhitespace(true)

	for _, block := range blocks {
		end := p.i + len(block)
		if end <= len(p.text) && p.text[p.i:end] == block {
			p.i = end
			return true
		}
	}
	return false
}

// parseCharacter parses a specific character
func (p *Parser) parseCharacter(char rune) bool {
	if p.i < len(p.text) {
		r, size := utf8.DecodeRuneInString(p.text[p.i:])
		if r == char {
			p.output.WriteRune(r)
			p.i += size
			return true
		}
	}
	return false
}

// skipCharacter skips a specific character without outputting it
func (p *Parser) skipCharacter(char rune) bool {
	if p.i < len(p.text) {
		r, size := utf8.DecodeRuneInString(p.text[p.i:])
		if r == char {
			p.i += size
			return true
		}
	}
	return false
}

// skipEscapeCharacter skips an escape character
func (p *Parser) skipEscapeCharacter() bool {
	return p.skipCharacter('\\')
}

// skipEllipsis skips ellipsis like "[1,2,3,...]"
func (p *Parser) skipEllipsis() bool {
	p.parseWhitespaceAndSkipComments(true)

	if p.i+2 < len(p.text) && p.text[p.i] == '.' && p.text[p.i+1] == '.' && p.text[p.i+2] == '.' {
		p.i += 3
		p.parseWhitespaceAndSkipComments(true)
		p.skipCharacter(',')
		return true
	}
	return false
}

// parseObject parses a JSON object
func (p *Parser) parseObject() bool {
	if p.i >= len(p.text) {
		return false
	}

	r, size := utf8.DecodeRuneInString(p.text[p.i:])
	if r != '{' {
		return false
	}

	p.output.WriteRune('{')
	p.i += size
	p.parseWhitespaceAndSkipComments(true)

	// Skip leading comma
	if p.skipCharacter(',') {
		p.parseWhitespaceAndSkipComments(true)
	}

	initial := true
	for p.i < len(p.text) {
		r, _ := getCharAt(p.text, p.i)
		if r == '}' {
			break
		}

		var processedComma bool
		if !initial {
			processedComma = p.parseCharacter(',')
			if !processedComma {
				// Repair missing comma
				output := p.output.String()
				p.output.Reset()
				p.output.WriteString(insertBeforeLastWhitespace(output, ","))
			}
			p.parseWhitespaceAndSkipComments(true)
		} else {
			processedComma = true
			initial = false
		}

		p.skipEllipsis()

		processedKey := p.parseString(false, -1) || p.parseUnquotedString(true)
		if !processedKey {
			r, _ := getCharAt(p.text, p.i)
			if r == '}' || r == '{' || r == ']' || r == '[' || p.i >= len(p.text) {
				// Repair trailing comma
				output := p.output.String()
				p.output.Reset()
				p.output.WriteString(stripLastOccurrence(output, ",", false))
			} else {
				return false
			}
			break
		}

		p.parseWhitespaceAndSkipComments(true)
		processedColon := p.parseCharacter(':')
		truncatedText := p.i >= len(p.text)

		if !processedColon {
			if p.i < len(p.text) && isStartOfValue(p.text, p.i) || truncatedText {
				// Repair missing colon
				output := p.output.String()
				p.output.Reset()
				p.output.WriteString(insertBeforeLastWhitespace(output, ":"))
			} else {
				return false
			}
		}

		processedValue := p.parseValue()
		if !processedValue {
			if processedColon || truncatedText {
				// Repair missing object value
				p.output.WriteString("null")
			} else {
				return false
			}
		}
	}

	if p.i < len(p.text) {
		r, size := utf8.DecodeRuneInString(p.text[p.i:])
		if r == '}' {
			p.output.WriteRune('}')
			p.i += size
		} else {
			// Repair missing end bracket
			output := p.output.String()
			p.output.Reset()
			p.output.WriteString(insertBeforeLastWhitespace(output, "}"))
		}
	} else {
		// Repair missing end bracket
		output := p.output.String()
		p.output.Reset()
		p.output.WriteString(insertBeforeLastWhitespace(output, "}"))
	}

	return true
}

// parseArray parses a JSON array
func (p *Parser) parseArray() bool {
	if p.i >= len(p.text) {
		return false
	}

	r, size := utf8.DecodeRuneInString(p.text[p.i:])
	if r != '[' {
		return false
	}

	p.output.WriteRune('[')
	p.i += size
	p.parseWhitespaceAndSkipComments(true)

	// Skip leading comma
	if p.skipCharacter(',') {
		p.parseWhitespaceAndSkipComments(true)
	}

	initial := true
	for p.i < len(p.text) {
		r, _ := getCharAt(p.text, p.i)
		if r == ']' {
			break
		}

		if !initial {
			processedComma := p.parseCharacter(',')
			if !processedComma {
				// Repair missing comma
				output := p.output.String()
				p.output.Reset()
				p.output.WriteString(insertBeforeLastWhitespace(output, ","))
			}
		} else {
			initial = false
		}

		p.skipEllipsis()

		processedValue := p.parseValue()
		if !processedValue {
			// Repair trailing comma
			output := p.output.String()
			p.output.Reset()
			p.output.WriteString(stripLastOccurrence(output, ",", false))
			break
		}
	}

	if p.i < len(p.text) {
		r, size := utf8.DecodeRuneInString(p.text[p.i:])
		if r == ']' {
			p.output.WriteRune(']')
			p.i += size
		} else {
			// Repair missing closing bracket
			output := p.output.String()
			p.output.Reset()
			p.output.WriteString(insertBeforeLastWhitespace(output, "]"))
		}
	} else {
		// Repair missing closing bracket
		output := p.output.String()
		p.output.Reset()
		p.output.WriteString(insertBeforeLastWhitespace(output, "]"))
	}

	return true
}

// parseNewlineDelimitedJSON repairs newline delimited JSON
func (p *Parser) parseNewlineDelimitedJSON() {
	// Note: The first value has already been parsed in Parse()
	// and a comma has already been added if needed
	// We just need to parse the remaining values
	initial := true
	processedValue := true

	for processedValue {
		if !initial {
			processedComma := p.parseCharacter(',')
			if !processedComma {
				// Repair: add missing comma
				output := p.output.String()
				p.output.Reset()
				p.output.WriteString(insertBeforeLastWhitespace(output, ","))
			}
		} else {
			initial = false
		}
		processedValue = p.parseValue()
	}

	// Remove trailing comma if any
	output := p.output.String()
	p.output.Reset()
	p.output.WriteString(stripLastOccurrence(output, ",", false))

	// Wrap in array brackets
	result := p.output.String()
	p.output.Reset()
	p.output.WriteString("[\n")
	p.output.WriteString(result)
	p.output.WriteString("\n]")
}

// parseString parses a JSON string (to be continued in next part due to complexity)
func (p *Parser) parseString(stopAtDelimiter bool, stopAtIndex int) bool {
	if p.i >= len(p.text) {
		return false
	}

	// Check for escaped string
	skipEscapeChars := false
	if p.text[p.i] == '\\' {
		p.i++
		skipEscapeChars = true
	}

	if p.i >= len(p.text) {
		return false
	}

	r, size := utf8.DecodeRuneInString(p.text[p.i:])
	if !isQuote(r) {
		if skipEscapeChars {
			p.i-- // Restore position
		}
		return false
	}

	// Determine end quote function
	var isEndQuote func(rune) bool
	if isDoubleQuote(r) {
		isEndQuote = isDoubleQuote
	} else if isSingleQuote(r) {
		isEndQuote = isSingleQuote
	} else if isSingleQuoteLike(r) {
		isEndQuote = isSingleQuoteLike
	} else {
		isEndQuote = isDoubleQuoteLike
	}

	iBefore := p.i
	oBefore := p.output.Len()

	p.output.WriteRune('"')
	p.i += size

	for {
		if p.i >= len(p.text) {
			// Missing end quote
			iPrev := p.prevNonWhitespaceIndex(p.i - 1)
			if iPrev >= 0 && iPrev < len(p.text) {
				prevR, _ := getCharAt(p.text, iPrev)
				if !stopAtDelimiter && isDelimiter(prevR) {
					// Retry parsing
					p.i = iBefore
					outputStr := p.output.String()
					p.output.Reset()
					p.output.WriteString(outputStr[:oBefore])
					return p.parseString(true, -1)
				}
			}

			// Repair missing quote
			currentOutput := p.output.String()
			p.output.Reset()
			p.output.WriteString(insertBeforeLastWhitespace(currentOutput, "\""))
			return true
		}

		if p.i == stopAtIndex {
			// Use stop index
			currentOutput := p.output.String()
			p.output.Reset()
			p.output.WriteString(insertBeforeLastWhitespace(currentOutput, "\""))
			return true
		}

		currentR, currentSize := utf8.DecodeRuneInString(p.text[p.i:])

		if isEndQuote(currentR) {
			// Potential end quote
			iQuote := p.i
			oQuote := p.output.Len()
			p.output.WriteRune('"')
			p.i += currentSize

			p.parseWhitespaceAndSkipComments(false)

			nextR, _ := getCharAt(p.text, p.i)
			if stopAtDelimiter || p.i >= len(p.text) ||
				isDelimiter(nextR) || isQuote(nextR) || isDigit(nextR) {
				// Valid end quote
				p.parseConcatenatedString()
				return true
			}

			iPrevChar := p.prevNonWhitespaceIndex(iQuote - 1)
			if iPrevChar >= 0 && iPrevChar < len(p.text) {
				prevChar, _ := getCharAt(p.text, iPrevChar)
				if prevChar == ',' {
					// Comma before quote - retry
					p.i = iBefore
					outputStr := p.output.String()
					p.output.Reset()
					p.output.WriteString(outputStr[:oBefore])
					return p.parseString(false, iPrevChar)
				}

				if isDelimiter(prevChar) {
					// Delimiter before quote - retry
					p.i = iBefore
					outputStr := p.output.String()
					p.output.Reset()
					p.output.WriteString(outputStr[:oBefore])
					return p.parseString(true, -1)
				}
			}

			// Not a real end quote, continue
			outputStr := p.output.String()
			p.output.Reset()
			p.output.WriteString(outputStr[:oQuote+1])
			p.i = iQuote + currentSize

			// Repair unescaped quote - insert backslash at oQuote position
			currentOutput := p.output.String()
			p.output.Reset()
			p.output.WriteString(currentOutput[:oQuote])
			p.output.WriteString("\\")
			p.output.WriteString(currentOutput[oQuote:])

		} else if stopAtDelimiter && isUnquotedStringDelimiter(currentR) {
			// Stop at delimiter
			if p.i > 0 && p.text[p.i-1] == ':' && matchesUrlStart(p.text, iBefore+1, p.i+2) {
				// Handle URL - write directly to output
				for p.i < len(p.text) {
					r, size := utf8.DecodeRuneInString(p.text[p.i:])
					if matchesUrlChar(r) {
						p.output.WriteRune(r)
						p.i += size
					} else {
						break
					}
				}
			}

			// Repair missing quote
			currentOutput := p.output.String()
			p.output.Reset()
			p.output.WriteString(insertBeforeLastWhitespace(currentOutput, "\""))
			p.parseConcatenatedString()
			return true

		} else if currentR == '\\' {
			// Handle escape sequences
			if p.i+1 < len(p.text) {
				nextChar, nextSize := utf8.DecodeRuneInString(p.text[p.i+1:])

				// Check for truncated unicode: \\uXX or \\uXXX (less than 4 hex digits at end of text)
				if nextChar == '\\' && p.i+2 < len(p.text) && p.text[p.i+2] == 'u' {
					// Potential truncated unicode escape
					j := 3 // \, \, u already counted
					for j < 7 && p.i+j < len(p.text) && isHex(rune(p.text[p.i+j])) {
						j++
					}
					// If we're at end of text and have less than 6 chars total (\\uXXXX), it's truncated
					if p.i+j >= len(p.text) && j < 7 {
						// Truncated unicode - jump to end to trigger missing quote repair
						p.i = len(p.text)
						continue
					}
				}

				if _, ok := escapeCharacters[nextChar]; ok {
					p.output.WriteRune(currentR)
					p.output.WriteRune(nextChar)
					p.i += currentSize + nextSize
				} else if nextChar == 'u' {
					// Unicode escape
					j := 2
					for j < 6 && p.i+j < len(p.text) && isHex(rune(p.text[p.i+j])) {
						j++
					}
					if j == 6 {
						p.output.WriteString(p.text[p.i : p.i+6])
						p.i += 6
					} else if p.i+j >= len(p.text) {
						// Truncated unicode - skip these characters and treat as end of string
						// Jump to end to trigger missing quote repair
						p.i = len(p.text)
					} else {
						return false
					}
				} else {
					// Invalid escape - remove backslash
					p.output.WriteRune(nextChar)
					p.i += currentSize + nextSize
				}
			} else {
				p.i += currentSize
			}
		} else {
			// Regular character
			if currentR == '"' && (p.i == 0 || p.text[p.i-1] != '\\') {
				// Unescaped double quote
				p.output.WriteString("\\\"")
				p.i += currentSize
			} else if isControlCharacter(currentR) {
				// Control character
				if escaped, ok := controlCharacters[currentR]; ok {
					p.output.WriteString(escaped)
				}
				p.i += currentSize
			} else {
				if !isValidStringCharacter(currentR) {
					return false
				}
				p.output.WriteRune(currentR)
				p.i += currentSize
			}
		}

		if skipEscapeChars {
			p.skipEscapeCharacter()
		}
	}
}

// parseConcatenatedString repairs concatenated strings like "hello" + "world"
func (p *Parser) parseConcatenatedString() bool {
	processed := false

	p.parseWhitespaceAndSkipComments(true)
	for p.i < len(p.text) && p.text[p.i] == '+' {
		processed = true
		p.i++
		p.parseWhitespaceAndSkipComments(true)

		// Remove end quote of first string
		output := p.output.String()
		p.output.Reset()
		p.output.WriteString(stripLastOccurrence(output, "\"", true))

		start := p.output.Len()
		parsedStr := p.parseString(false, -1)
		if parsedStr {
			// Remove start quote of second string
			output := p.output.String()
			p.output.Reset()
			p.output.WriteString(removeAtIndex(output, start, 1))
		} else {
			// Remove the + because it's not followed by a string
			output := p.output.String()
			p.output.Reset()
			p.output.WriteString(insertBeforeLastWhitespace(output, "\""))
		}
	}

	return processed
}

// parseNumber parses a JSON number
func (p *Parser) parseNumber() bool {
	start := p.i

	if p.i < len(p.text) && p.text[p.i] == '-' {
		p.i++
		if p.atEndOfNumber() {
			p.repairNumberEndingWithNumericSymbol(start)
			return true
		}
		if p.i >= len(p.text) || !isDigit(rune(p.text[p.i])) {
			p.i = start
			return false
		}
	}

	// Integer part
	for p.i < len(p.text) && isDigit(rune(p.text[p.i])) {
		p.i++
	}

	// Decimal part
	if p.i < len(p.text) && p.text[p.i] == '.' {
		p.i++
		if p.atEndOfNumber() {
			p.repairNumberEndingWithNumericSymbol(start)
			return true
		}
		if p.i >= len(p.text) || !isDigit(rune(p.text[p.i])) {
			p.i = start
			return false
		}
		for p.i < len(p.text) && isDigit(rune(p.text[p.i])) {
			p.i++
		}
	}

	// Exponent part
	if p.i < len(p.text) && (p.text[p.i] == 'e' || p.text[p.i] == 'E') {
		p.i++
		if p.i < len(p.text) && (p.text[p.i] == '-' || p.text[p.i] == '+') {
			p.i++
		}
		if p.atEndOfNumber() {
			p.repairNumberEndingWithNumericSymbol(start)
			return true
		}
		if p.i >= len(p.text) || !isDigit(rune(p.text[p.i])) {
			p.i = start
			return false
		}
		for p.i < len(p.text) && isDigit(rune(p.text[p.i])) {
			p.i++
		}
	}

	if !p.atEndOfNumber() {
		p.i = start
		return false
	}

	if p.i > start {
		num := p.text[start:p.i]
		// Check for leading zeros
		if len(num) > 1 && num[0] == '0' && num[1] >= '0' && num[1] <= '9' {
			// Has invalid leading zero - quote it
			p.output.WriteString("\"")
			p.output.WriteString(num)
			p.output.WriteString("\"")
		} else {
			p.output.WriteString(num)
		}
		return true
	}

	return false
}

// parseKeywords parses JSON keywords (true, false, null) and Python variants
func (p *Parser) parseKeywords() bool {
	return p.parseKeyword("true", "true") ||
		p.parseKeyword("false", "false") ||
		p.parseKeyword("null", "null") ||
		p.parseKeyword("True", "true") ||
		p.parseKeyword("False", "false") ||
		p.parseKeyword("None", "null")
}

// parseKeyword parses a specific keyword
func (p *Parser) parseKeyword(name, value string) bool {
	end := p.i + len(name)
	if end <= len(p.text) && p.text[p.i:end] == name {
		p.output.WriteString(value)
		p.i = end
		return true
	}
	return false
}

// parseUnquotedString parses an unquoted string and adds quotes
func (p *Parser) parseUnquotedString(isKey bool) bool {
	start := p.i

	if p.i < len(p.text) {
		r, _ := getCharAt(p.text, p.i)
		if isFunctionNameCharStart(r) {
			for p.i < len(p.text) {
				r, size := utf8.DecodeRuneInString(p.text[p.i:])
				if isFunctionNameChar(r) {
					p.i += size
				} else {
					break
				}
			}

			// Check for function call
			j := p.i
			for j < len(p.text) && isWhitespace(p.text, j) {
				j++
			}

			if j < len(p.text) && p.text[j] == '(' {
				// Function call like NumberLong(2) or Timestamp(1234, 1) or callback({})
				p.i = j + 1

				// Parse the first value
				p.parseValue()

				// Skip any additional arguments (e.g., Timestamp(1234, 1))
				for p.i < len(p.text) && p.text[p.i] == ',' {
					p.i++ // skip comma
					// Save output BEFORE parsing whitespace to avoid trailing spaces
					savedOutput := p.output.String()
					p.parseWhitespaceAndSkipComments(true)
					// Skip this value - we only keep the first one
					p.parseValue()
					// Restore output to discard this value and any whitespace before it
					p.output.Reset()
					p.output.WriteString(savedOutput)
				}

				if p.i < len(p.text) && p.text[p.i] == ')' {
					p.i++
					if p.i < len(p.text) && p.text[p.i] == ';' {
						p.i++
					}
				}
				return true
			}
		}
	}

	// Parse unquoted string
	for p.i < len(p.text) {
		r, size := utf8.DecodeRuneInString(p.text[p.i:])
		if isUnquotedStringDelimiter(r) || isQuote(r) || (isKey && r == ':') {
			break
		}
		p.i += size
	}

	// Check for URL
	if p.i > 0 && p.text[p.i-1] == ':' && matchesUrlStart(p.text, start, p.i+2) {
		for p.i < len(p.text) {
			r, size := utf8.DecodeRuneInString(p.text[p.i:])
			if matchesUrlChar(r) {
				p.i += size
			} else {
				break
			}
		}
	}

	if p.i > start {
		// Remove trailing whitespace
		for p.i > start && isWhitespace(p.text, p.i-1) {
			p.i--
		}

		symbol := p.text[start:p.i]
		if symbol == "undefined" {
			p.output.WriteString("null")
		} else {
			// Quote the string
			jsonStr, _ := json.Marshal(symbol)
			p.output.WriteString(string(jsonStr))
		}

		// Skip end quote if present
		if p.i < len(p.text) && p.text[p.i] == '"' {
			p.i++
		}

		return true
	}

	return false
}

// parseRegex parses a regex literal and converts it to a string
func (p *Parser) parseRegex() bool {
	if p.i >= len(p.text) || p.text[p.i] != '/' {
		return false
	}

	start := p.i
	p.i++

	for p.i < len(p.text) && (p.text[p.i] != '/' || (p.i > 0 && p.text[p.i-1] == '\\')) {
		p.i++
	}

	if p.i < len(p.text) {
		p.i++ // Skip closing /
	}

	p.output.WriteString("\"")
	p.output.WriteString(p.text[start:p.i])
	p.output.WriteString("\"")
	return true
}

// Helper methods

func (p *Parser) prevNonWhitespaceIndex(start int) int {
	prev := start
	for prev > 0 && isWhitespace(p.text, prev) {
		prev--
	}
	return prev
}

func (p *Parser) atEndOfNumber() bool {
	return p.i >= len(p.text) || isDelimiter(rune(p.text[p.i])) || isWhitespace(p.text, p.i)
}

func (p *Parser) repairNumberEndingWithNumericSymbol(start int) {
	p.output.WriteString(p.text[start:p.i])
	p.output.WriteString("0")
}

// Error methods

func (p *Parser) throwUnexpectedCharacter() error {
	char := ""
	if p.i < len(p.text) {
		char = strconv.QuoteRune(rune(p.text[p.i]))
	}
	return NewJSONRepairError("Unexpected character "+char, p.i)
}

func (p *Parser) throwUnexpectedEnd() error {
	return NewJSONRepairError("Unexpected end of json string", len(p.text))
}
