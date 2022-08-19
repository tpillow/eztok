package eztok

import (
	"fmt"
	"log"
)

// Returns a new CallbackNode whose CanParseToken function returns true if
// Context.PeekRune(0) equals runeValue and whose ParseToken function returns a
// Token with a TokenType of tokenType and a Value of runeValue.
func NewRuneMatchNode(tokenType TokenType, runeValue rune) *CallbackNode {
	return NewCallbackNode(
		func(ctx Context) bool {
			return ctx.PeekRune(0) == runeValue
		},
		func(ctx Context) (*Token, error) {
			if nextRune := ctx.NextRune(); nextRune != runeValue {
				return nil, fmt.Errorf("expected a '%c' rune but got a '%c' rune",
					runeValue, nextRune)
			}
			return NewToken(tokenType, runeValue), nil
		},
	)
}

// Returns a new CallbackNode whose CanParseToken function returns true if
// the next runes in the Context equal the string value and whose ParseToken
// function returns a Token with a TokenType of tokenType and a Value of value.
func NewStringMatchNode(tokenType TokenType, value string) *CallbackNode {
	if len(value) <= 0 {
		log.Panicf("Cannot create a NewStringMatchNode with an empty string to match on.")
	}
	return NewCallbackNode(
		func(ctx Context) bool {
			// Peek full length of value first, to better cache for some context implementations
			for i := len(value) - 1; i >= 0; i-- {
				if ctx.PeekRune(i) != rune(value[i]) {
					return false
				}
			}
			return true
		},
		func(ctx Context) (*Token, error) {
			for i := 0; i < len(value); i++ {
				if r := ctx.NextRune(); r != rune(value[i]) {
					return nil, fmt.Errorf("expected word '%v' but got mismatched rune '%c' at position %v", value, r, i)
				}
			}
			return NewToken(tokenType, value), nil
		},
	)
}
