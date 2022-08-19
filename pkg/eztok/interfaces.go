package eztok

// Alias to rune(0). Used as nil in the context of runes.
const NilRune rune = 0

// Represents something that can look ahead by n input runes, and
// consume input runes in-order. A Context should also track Origin
// information for the current location in the input.
type Context interface {
	// Return the rune that is relative runes ahead of the current
	// rune in the input. Returns NilRune if there is none.
	PeekRune(relative int) rune
	// Consume (i.e. advance the input stream by 1 rune) and return the
	// consumed rune. Returns NilRune if there is none.
	NextRune() rune
	// Returns the Origin information of the rune that would be returned
	// by a call to NextRune().
	GetNextOrigin() *Origin
}

// Represents a Node in a node-based tokenizer.
type Node interface {
	// Returns true if, given the current Context state, this node
	// thinks that it can parse a token.
	CanParseToken(ctx Context) bool
	// Parses a token given the current Context state. This will only ever
	// be called if the corresponding CanParseToken(...) call returned true
	// for the current Context state.
	// Returns an error if the token could not be properly parsed. Returns a
	// nil Token if the token was properly parsed, but should be ignored
	// (i.e. consume input, but do not generate a token). Returns a valid Token
	// if the token was properly parsed.
	ParseToken(ctx Context) (*Token, error)
}

// Represents something that can convert a Context into a stream of Token
// objects.
type Tokenizer interface {
	// Convert the given Context state into a stream of Token objects.
	// Returns an error if tokenization fails. Returns a list of Token objects
	// if tokenization succeeds.
	Tokenize(ctx Context) ([]*Token, error)
}

// Represents something that can look 1 token ahead in a Token stream, and
// consume & return the next Token in a Token stream.
type Traverser interface {
	// Return the Token that is relative Tokens ahead of the current Token in
	// the input. Returns nil if there is none.
	PeekToken(relative int) *Token
	// Consume (i.e. advance the input stream by 1 Token) and return the
	// consumed Token. Returns nil if no Tokens remain.
	NextToken() *Token
}
