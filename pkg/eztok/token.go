package eztok

import "fmt"

// TokenType is just a string that should be unique per token type.
type TokenType string

// Represents a parsed token.
type Token struct {
	// The type of the token.
	TokenType TokenType
	// The value of the token.
	Value any
	// Origin information of where the token came from.
	OriginInfo *OriginInfo
}

// Creates a new Token with the given TokenType and Value, and nil OriginInfo.
func NewToken(tokenType TokenType, value any) *Token {
	return &Token{
		TokenType:  tokenType,
		Value:      value,
		OriginInfo: nil,
	}
}

// Formats the Token to a string containing TokenType and Value information.
func (token *Token) ToString() string {
	return fmt.Sprintf("Token[%v](%v)", token.TokenType, token.Value)
}
