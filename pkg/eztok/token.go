package eztok

import "fmt"

// Alias to a string. Represents the type of a Token.
type TokenType string

// Represents a Token.
type Token struct {
	// The type of this Token.
	TokenType TokenType
	// The value of this Token. May be nil.
	Value any
	// The Origin information of where this Token came from.
	Origin *Origin
}

// Returns a new Token object with the given parameters and a nil Origin.
func NewToken(tokenType TokenType, value any) *Token {
	return &Token{tokenType, value, nil}
}

// Returns a string representation of the Token containing its TokenType
// and Value.
func (token *Token) ToString() string {
	return fmt.Sprintf("%v (%v)", token.TokenType, token.Value)
}
