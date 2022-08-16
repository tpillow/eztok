package eztok

// A callback function used in TokenizerNode. See TokenizerNode.
type TokenizerPeeker func(ctx *TokenizerCtx) bool

// A callback function used in TokenizerNode. See TokenizerNode.
type TokenizerParser func(ctx *TokenizerCtx) (*Token, error)

// Used to represent part of a TokenizerDef that can check if the current state of a
// TokenizerCtx (i.e. next rune) can be processed, and will attempt to process it.
type TokenizerNode struct {
	// Returns true if the current state of the TokenizerCtx might be able to be processed
	// by this TokenizerNode. This should not consume any tokens from the TokenizerCtx.
	Peeker TokenizerPeeker
	// Returns a new Token, or error, based on the current state of the TokenizerCtx.
	// This will only be called if the corresponding Peeker returned true.
	Parser TokenizerParser
}

// Creates a new TokenizerNode with the given TokenizerPeeker and TokenizerParser.
func NewTokenizerNode(peeker TokenizerPeeker, parser TokenizerParser) *TokenizerNode {
	return &TokenizerNode{peeker, parser}
}

// Creates a new TokenizerNode that will create a token of type tokenType when the
// single rune r is next in the TokenizerCtx state.
func NewSingleRuneTokenizerNode(tokenType TokenType, r rune) *TokenizerNode {
	return NewTokenizerNode(
		func(ctx *TokenizerCtx) bool {
			return ctx.PeekIs(r)
		}, func(ctx *TokenizerCtx) (*Token, error) {
			ctx.AssertNextIs(r)
			return NewToken(tokenType, r), nil
		})
}
