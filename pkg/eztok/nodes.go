package eztok

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// TokenType definitions returned by some common Node.ParseToken calls.
const (
	TokenTypeIdentifier TokenType = "identifier"
	TokenTypeInteger    TokenType = "integer"
	TokenTypeFloat      TokenType = "float"
	TokenTypeString     TokenType = "string"
)

// Holds a map of escape rune (a rune following a '\' in a string) to the
// actual string contents of the escape.
// TODO: make this more complete
var escapedRuneToString = map[rune]string{
	'\\': "\\",
	'\'': "'",
	'"':  "\"",
	'n':  "\n",
	'r':  "\r",
	't':  "\t",
}

// Returns true if the provided TokenType holds a numeric Value. Only works
// for TokenType values within the ezcommon package.
func IsNumericTokenType(tokenType TokenType) bool {
	switch tokenType {
	case TokenTypeInteger, TokenTypeFloat:
		return true
	}
	return false
}

// A node that matches any whitespace rune and skips it (i.e. generates no Token).
var SkipWhitespaceNode = NewCallbackNode(
	func(ctx Context) bool {
		return unicode.IsSpace(ctx.PeekRune(0))
	},
	func(ctx Context) (*Token, error) {
		if r := ctx.NextRune(); !unicode.IsSpace(r) {
			return nil, fmt.Errorf("expected a whitespace rune but got a '%c' rune", r)
		}
		return nil, nil
	},
)

// A node that matches a C-style identifier. This means a word that begins with either
// a letter or '_', and is followed by 0 or more letters, digits, or '_'s.
var IdentifierNode = NewCallbackNode(
	func(ctx Context) bool {
		r := ctx.PeekRune(0)
		return unicode.IsLetter(r) || r == '_'
	},
	func(ctx Context) (*Token, error) {
		str := ReadRunesUntil(ctx, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_'
		})
		if len(str) <= 0 {
			return nil, fmt.Errorf("expected an identifier but got an empty identifier")
		}
		return NewToken(TokenTypeIdentifier, str), nil
	},
)

// A node that matches an integer or float number. Base-10 numbers can be optionally
// preceded by a '-' or '+' symbol, representing the sign of the number. The
// following non-base-10 numbers can also be specified in the following formats:
// - 0b### (base-2, binary)
// - 0o### (base-8, octal)
// - 0x### (base-16, hexadecimal)
var NumberNode = NewCallbackNode(
	func(ctx Context) bool {
		r := ctx.PeekRune(0)
		return unicode.IsDigit(r) || ((r == '+' || r == '-') && unicode.IsDigit(ctx.PeekRune(1)))
	},
	func(ctx Context) (*Token, error) {
		firstRune := ctx.NextRune()
		str := ""
		base := 10
		header := ""
		if firstRune == '0' && !unicode.IsDigit(ctx.PeekRune(0)) {
			// Change-of-base (disallow preceding +/-)
			r := ctx.NextRune()
			switch r {
			case 'b':
				base = 2
			case 'o':
				base = 8
			case 'x':
				base = 16
			default:
				return nil, fmt.Errorf("invalid numeric base '0%c' (must be one of 'b', 'o', 'x')", r)
			}
			header = fmt.Sprintf("0%c", r)
		} else {
			str = string(firstRune)
		}

		str += ReadRunesUntil(ctx, func(r rune) bool {
			return !unicode.IsDigit(r) && !unicode.IsLetter(r) && r != '.' && r != '_'
		})
		str = strings.ReplaceAll(str, "_", "")
		if len(str) <= 0 {
			return nil, fmt.Errorf("expected a number but got '%v%v'", header, str)
		}

		if base != 10 || !strings.Contains(str, ".") {
			// Integer
			intVal, err := strconv.ParseInt(str, base, 64)
			if err != nil {
				return nil, fmt.Errorf("expected an integer but got '%v%v': %v", header, str, err)
			}
			return NewToken(TokenTypeInteger, intVal), nil
		} else {
			// Float
			floatVal, err := strconv.ParseFloat(str, 64)
			if err != nil {
				return nil, fmt.Errorf("expected a float but got '%v%v': %v", header, str, err)
			}
			return NewToken(TokenTypeFloat, floatVal), nil
		}
	},
)

// A node that matches a string encased in double quotes ("). Allows for common
// escape characters by using a backslash ('\<escape>') in the string.
var DoubleQuotedEscapedStringNode = newStringNode('"')

// A node that matches a string encased in single quotes ('). Allows for common
// escape characters by using a backslash ('\<escape>') in the string.
var SingleQuotedEscapedStringNode = newStringNode('\'')

// Returns a CallbackNode that can parse an escaped string encased in
// the provided quoteRune rune.
func newStringNode(quoteRune rune) Node {
	return NewCallbackNode(
		func(ctx Context) bool {
			return ctx.PeekRune(0) == quoteRune
		},
		func(ctx Context) (*Token, error) {
			ctx.NextRune()
			str := ""
			nextIsEscaped := false

			for r := ctx.PeekRune(0); r != NilRune && (nextIsEscaped || r != quoteRune); r = ctx.PeekRune(0) {
				check := ctx.NextRune()
				if nextIsEscaped {
					replStr, ok := escapedRuneToString[check]
					if !ok {
						return nil, fmt.Errorf("unknown escape '%c' while tokenizing string '%v'", check, str)
					}
					str += replStr
					nextIsEscaped = false
				} else if check == '\\' {
					nextIsEscaped = true
				} else {
					str += string(check)
				}
			}
			r := ctx.NextRune()
			if r != quoteRune || nextIsEscaped {
				return nil, fmt.Errorf("unterminated string '%v'", str)
			}

			return NewToken(TokenTypeString, str), nil
		},
	)
}
