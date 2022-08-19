package eztok

// Returns true if the Token returned by Traverser.PeekToken has a Token.TokenType
// matching tokenType. Returns false otherwise.
func PeekTokenTypeIs(trav Traverser, tokenType TokenType) bool {
	tok := trav.PeekToken(0)
	if tok == nil {
		return false
	}
	return tok.TokenType == tokenType
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
