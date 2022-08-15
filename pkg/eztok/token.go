package eztok

type TokenType string

type Token struct {
	TokenType TokenType
	Value     any
}

func NewToken(tokenType TokenType, value any) *Token {
	return &Token{tokenType, value}
}
