package eztok

import (
	"fmt"
)

// Returns true if the Token returned by Traverser.PeekToken has a Token.TokenType
// matching tokenType. Returns false otherwise.
func PeekTokenTypeIs(trav Traverser, tokenType TokenType) bool {
	tok := trav.PeekToken(0)
	if tok == nil {
		return false
	}
	return tok.TokenType == tokenType
}

// Returns the Token returned by Traverser.NextToken if its Token.TokenType matches
// tokenType. Otherwise, log.Fatalf is called with a corresponding error and nil is
// returned.
func ExpectTokenTypeIs(trav Traverser, tokenType TokenType) (*Token, error) {
	tok := trav.NextToken()
	if tok == nil {
		return nil, fmt.Errorf("expected token with type '%v' but got end of input", tokenType)
	}
	if tok.TokenType != tokenType {
		return nil, fmt.Errorf("expected token with type '%v' but got type '%v' at %v",
			tokenType, tok.TokenType, tok.Origin.ToString())
	}
	return tok, nil
}

// Returns true if the Token returned by Traverser.PeekToken has a Token.TokenType
// matching tokenType and a Token.Value matching value. Returns false otherwise.
func PeekTokenValueIs(trav Traverser, tokenType TokenType, value any) bool {
	tok := trav.PeekToken(0)
	if tok == nil {
		return false
	}
	return tok.TokenType == tokenType && tok.Value == value
}

// Returns the Token returned by Traverser.NextToken if its Token.TokenType matches
// tokenType and its Token.Value matches value. Otherwise, log.Fatalf is called with
// a corresponding error and nil is returned.
func ExpectTokenValueIs(trav Traverser, tokenType TokenType, value any) (*Token, error) {
	tok := trav.NextToken()
	if tok == nil {
		return nil, fmt.Errorf("expected token with type '%v' and value '%v' but got end of input", tokenType, value)
	}
	if tok.TokenType != tokenType || tok.Value != value {
		return nil, fmt.Errorf("expected token with type '%v' and value '%v' but got type '%v' with value '%v' at %v",
			tokenType, value, tok.TokenType, tok.Value, tok.Origin.ToString())
	}
	return tok, nil
}
