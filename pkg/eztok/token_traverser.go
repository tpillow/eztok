package eztok

import "log"

// A Traverser that operates on a slice of Token objects.
type TokenTraverser struct {
	// The Token stream to operate on. The Token at index 0 will
	// always be the next Token in the stream. Consuming a Token
	// means removing the Token at index 0 of this slice, if it
	// exists.
	Tokens []*Token
}

// Returns a new TokenTraverser with the given parameters.
func NewTokenTraverser(tokens []*Token) *TokenTraverser {
	return &TokenTraverser{tokens}
}

// Return the Token that is relative Tokens ahead of the current Token in
// TokenTraverser.Tokens. Returns nil otherwise.
func (trav *TokenTraverser) PeekToken(relative int) *Token {
	if relative < 0 {
		log.Panicf("PeekToken cannot peek negatively; tried peeking a relative '%v' tokens", relative)
	}
	if relative < len(trav.Tokens) {
		return trav.Tokens[relative]
	}
	return nil
}

// Consumes the Token at index 0 of TokenTraverser.Tokens if its length
// is not 0 and returns the consumed Token. Returns nil otherwise.
func (trav *TokenTraverser) NextToken() *Token {
	if len(trav.Tokens) <= 0 {
		return nil
	}
	tok := trav.Tokens[0]
	trav.Tokens = trav.Tokens[1:]
	return tok
}
