package traverse

import (
	"fmt"

	"github.com/tpillow/eztok/pkg/eztok"
)

// An error handler callback function used by the TokenTraverser.
type ErrHandler func(err error)

// Holds data required to traverse a TokenizerResult.
type TokenTraverser struct {
	// The TokenizerResult to traverse.
	// NOTE: this will modify that data contained within the Result as necessary.
	Result *eztok.TokenizerResult
	// Holds the last error emitted by the TokenTraverser. Will be nil if no error has occurred.
	LastErr error
	// This callback will be called when the TokenTraverser emits an error, if it is not nil.
	ErrHandler ErrHandler
}

// Creates a new TokenTraverser to traverse the given result. Any error that is emitted will be
// processed by the errHandler callback if provided.
func NewTokenTraverser(result *eztok.TokenizerResult, errHandler ErrHandler) *TokenTraverser {
	return &TokenTraverser{result, nil, errHandler}
}

// Emits an error from the TokenTraverser. This will set LastErr to err and call the ErrHandler if not nil.
func (tt *TokenTraverser) EmitErr(err error) {
	tt.LastErr = err
	if tt.ErrHandler != nil {
		tt.ErrHandler(err)
	}
}

// Returns the next Token to process without consuming it, or nil if none remain.
func (tt *TokenTraverser) Peek() *eztok.Token {
	if len(tt.Result.Tokens) > 0 {
		return tt.Result.Tokens[0]
	}
	return nil
}

// Returns true if there are no more tokens to process. Analagous to Peek() == nil.
func (tt *TokenTraverser) IsAtEnd() bool {
	return tt.Peek() == nil
}

// Returns the next Token to process by consuming it, or nil if none remain.
func (tt *TokenTraverser) Next() *eztok.Token {
	if len(tt.Result.Tokens) > 0 {
		tok := tt.Result.Tokens[0]
		tt.Result.Tokens = tt.Result.Tokens[1:]
		return tok
	}
	return nil
}

// Returns the next Token to process by consuming it. If no tokens remain, an error will be
// emitted.
func (tt *TokenTraverser) NextErrIfEOF(expectingTokenType eztok.TokenType) *eztok.Token {
	tok := tt.Next()
	if tok == nil {
		tt.EmitErr(fmt.Errorf("expected a '%v' token but got EOF", expectingTokenType))
		return nil
	}
	return tok
}

// Analagous to (assuming non-nil): Peek().TokenType == tokenType
func (tt *TokenTraverser) PeekIsType(tokenType eztok.TokenType) bool {
	tok := tt.Peek()
	if tok == nil {
		return false
	}
	return tok.TokenType == tokenType
}

// Analagous to (assuming non-nil): Peek().TokenType in tokenTypes
func (tt *TokenTraverser) PeekAnyIsType(tokenTypes ...eztok.TokenType) bool {
	tok := tt.Peek()
	if tok == nil {
		return false
	}
	for _, tt := range tokenTypes {
		if tok.TokenType == tt {
			return true
		}
	}
	return false
}

// Returns the next Token to process by consuming it. If the token is nil or its TokenType
// does not match the provided tokenType, an error is emitted.
func (tt *TokenTraverser) ExpectType(tokenType eztok.TokenType) *eztok.Token {
	tok := tt.NextErrIfEOF(tokenType)
	if tok.TokenType != tokenType {
		tt.EmitErr(fmt.Errorf(
			"expected a '%v' token but got a '%v' token with value '%v' at %v",
			tokenType, tok.TokenType, tok.Value, tok.OriginInfo.ToString()))
		return nil
	}
	return tok
}

// Analagous to (assuming non-nil): Peek().TokenType == tokenType && Peek().Value == value
func (tt *TokenTraverser) PeekIsTypeAndVal(tokenType eztok.TokenType, value any) bool {
	tok := tt.Peek()
	if tok == nil {
		return false
	}
	return tok.TokenType == tokenType && tok.Value == value
}

// Returns the next Token to process by consuming it. If the token is nil or its TokenType
// does not match the provided tokenType or its Value does not match the provided value, an error
// is emitted.
func (tt *TokenTraverser) ExpectTypeAndVal(tokenType eztok.TokenType, value any) *eztok.Token {
	tok := tt.NextErrIfEOF(tokenType)
	if tok.TokenType != tokenType || tok.Value != value {
		tt.EmitErr(
			fmt.Errorf("expected a '%v' token with value '%v' but got a '%v' token with value '%v' at %v",
				tokenType, value, tok.TokenType, tok.Value, tok.OriginInfo.ToString()))
		return nil
	}
	return tok
}
