package jsonrepair

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// Compiled regular expressions
var (
	startOfValueRegex   = regexp.MustCompile(`^[[\{\w-]$`)
	urlStartRegex       = regexp.MustCompile(`^(http|https|ftp|mailto|file|data|irc)://$`)
	urlCharRegex        = regexp.MustCompile(`^[A-Za-z0-9\-._~:/?#@!$&'()*+;=]$`)
	commaOrNewlineRegex = regexp.MustCompile(`[,\n][ \t\r]*$`)
)

// isHex checks if a character is a hexadecimal digit
func isHex(char rune) bool {
	return (char >= '0' && char <= '9') || (char >= 'A' && char <= 'F') || (char >= 'a' && char <= 'f')
}

// isDigit checks if a character is a digit
func isDigit(char rune) bool {
	return char >= '0' && char <= '9'
}

// isValidStringCharacter checks if a character is valid in a JSON string
func isValidStringCharacter(char rune) bool {
	// Valid range is between \u0020 and \u10ffff
	return char >= '\u0020'
}

// isDelimiter checks if a character is a delimiter
func isDelimiter(char rune) bool {
	return strings.ContainsRune(",:[]/{}()\n+", char)
}

// isFunctionNameCharStart checks if a character can start a function name
func isFunctionNameCharStart(char rune) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || char == '_' || char == '$'
}

// isFunctionNameChar checks if a character can be in a function name
func isFunctionNameChar(char rune) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
		char == '_' || char == '$' || (char >= '0' && char <= '9')
}

// isUnquotedStringDelimiter checks if a character is an unquoted string delimiter
func isUnquotedStringDelimiter(char rune) bool {
	return strings.ContainsRune(",[]/{}+\n", char)
}

// isStartOfValue checks if a character marks the start of a value
func isStartOfValue(text string, index int) bool {
	if index >= len(text) {
		return false
	}
	r, _ := utf8.DecodeRuneInString(text[index:])
	return isQuote(r) || startOfValueRegex.MatchString(string(r))
}

// isControlCharacter checks if a character is a control character
func isControlCharacter(char rune) bool {
	return char == '\n' || char == '\r' || char == '\t' || char == '\b' || char == '\f'
}

// isWhitespace checks if the character at the given index is a whitespace
func isWhitespace(text string, index int) bool {
	if index >= len(text) {
		return false
	}
	r, _ := utf8.DecodeRuneInString(text[index:])
	return r == ' ' || r == '\n' || r == '\t' || r == '\r'
}

// isWhitespaceExceptNewline checks if the character is whitespace but not newline
func isWhitespaceExceptNewline(text string, index int) bool {
	if index >= len(text) {
		return false
	}
	r, _ := utf8.DecodeRuneInString(text[index:])
	return r == ' ' || r == '\t' || r == '\r'
}

// isSpecialWhitespace checks if the character is a special unicode whitespace
func isSpecialWhitespace(text string, index int) bool {
	if index >= len(text) {
		return false
	}
	r, _ := utf8.DecodeRuneInString(text[index:])
	return r == codeNonBreakingSpace ||
		(r >= codeEnQuad && r <= codeHairSpace) ||
		r == codeNarrowNoBreakSpace ||
		r == codeMediumMathematicalSpace ||
		r == codeIdeographicSpace
}

// isQuote checks if a character is any type of quote
func isQuote(char rune) bool {
	return isDoubleQuoteLike(char) || isSingleQuoteLike(char)
}

// isDoubleQuoteLike checks if a character is a double quote or similar
func isDoubleQuoteLike(char rune) bool {
	return char == '"' || char == '\u201c' || char == '\u201d'
}

// isDoubleQuote checks if a character is exactly a double quote
func isDoubleQuote(char rune) bool {
	return char == '"'
}

// isSingleQuoteLike checks if a character is a single quote or similar
func isSingleQuoteLike(char rune) bool {
	return char == '\'' || char == '\u2018' || char == '\u2019' || char == '\u0060' || char == '\u00b4'
}

// isSingleQuote checks if a character is exactly a single quote
func isSingleQuote(char rune) bool {
	return char == '\''
}

// stripLastOccurrence removes the last occurrence of a substring
func stripLastOccurrence(text, textToStrip string, stripRemainingText bool) string {
	index := strings.LastIndex(text, textToStrip)
	if index == -1 {
		return text
	}
	if stripRemainingText {
		return text[:index]
	}
	// Note: Original TypeScript uses index + 1, not index + textToStrip.length
	// This works correctly for single character strings
	return text[:index] + text[index+1:]
}

// insertBeforeLastWhitespace inserts text before the last whitespace characters
func insertBeforeLastWhitespace(text, textToInsert string) string {
	index := len(text)

	// If no trailing whitespace, append at the end
	if index == 0 || !isWhitespace(text, index-1) {
		return text + textToInsert
	}

	// Find the start of trailing whitespace
	for index > 0 && isWhitespace(text, index-1) {
		// Move back by the size of the character
		_, size := utf8.DecodeLastRuneInString(text[:index])
		index -= size
	}

	return text[:index] + textToInsert + text[index:]
}

// removeAtIndex removes count characters starting at the given index
func removeAtIndex(text string, start, count int) string {
	if start >= len(text) {
		return text
	}
	end := start + count
	if end > len(text) {
		end = len(text)
	}
	return text[:start] + text[end:]
}

// endsWithCommaOrNewline checks if text ends with a comma or newline (with optional whitespace)
func endsWithCommaOrNewline(text string) bool {
	return commaOrNewlineRegex.MatchString(text)
}

// atEndOfBlockComment checks if we're at the end of a block comment
func atEndOfBlockComment(text string, i int) bool {
	return i < len(text)-1 && text[i] == '*' && text[i+1] == '/'
}

// getCharAt safely gets the rune at the given index
func getCharAt(text string, index int) (rune, bool) {
	if index >= len(text) {
		return 0, false
	}
	r, _ := utf8.DecodeRuneInString(text[index:])
	return r, true
}

// matchesUrlStart checks if the text at index matches a URL start pattern
func matchesUrlStart(text string, start, end int) bool {
	if end > len(text) {
		end = len(text)
	}
	if start >= end {
		return false
	}
	return urlStartRegex.MatchString(text[start:end])
}

// matchesUrlChar checks if a character is valid in a URL
func matchesUrlChar(char rune) bool {
	return urlCharRegex.MatchString(string(char))
}
