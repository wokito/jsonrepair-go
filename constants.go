package jsonrepair

// Control characters mapping
var controlCharacters = map[rune]string{
	'\b': "\\b",
	'\f': "\\f",
	'\n': "\\n",
	'\r': "\\r",
	'\t': "\\t",
}

// Escape characters mapping
var escapeCharacters = map[rune]rune{
	'"':  '"',
	'\\': '\\',
	'/':  '/',
	'b':  '\b',
	'f':  '\f',
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
	// Note: \u is handled separately in parseString()
}

// Unicode code points for special whitespace characters
const (
	codeSpace                   = 0x20   // " "
	codeNewline                 = 0x0A   // "\n"
	codeTab                     = 0x09   // "\t"
	codeReturn                  = 0x0D   // "\r"
	codeNonBreakingSpace        = 0xA0   // non-breaking space
	codeEnQuad                  = 0x2000 // en quad
	codeHairSpace               = 0x200A // hair space
	codeNarrowNoBreakSpace      = 0x202F // narrow no-break space
	codeMediumMathematicalSpace = 0x205F // medium mathematical space
	codeIdeographicSpace        = 0x3000 // ideographic space
)
