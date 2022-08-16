package eztok

// Represents a final result from a TokenizerDef.
type TokenizerResult struct {
	// The list of tokens that were successfully processed.
	Tokens []*Token
}

// Creates a new TokenizerResult by appending the data in TokenResult b onto TokenResult a.
func MergeTokenizerResults(a *TokenizerResult, b *TokenizerResult) *TokenizerResult {
	return &TokenizerResult{append(a.Tokens, b.Tokens...)}
}

// Creates a new TokenizerResult by appending the data in TokenResult b into TokenResult a
// starting at the aInsertIdx index.
//
// This may be useful for processing something akin to a C-style "#include" statement, which would
// require tokenization of another file to be inserted at the current position of tokens being processed.
func MergeTokenizerResultsAt(a *TokenizerResult, aInsertIdx int, b *TokenizerResult) *TokenizerResult {
	a0 := a.Tokens[:aInsertIdx]
	a1 := a.Tokens[aInsertIdx:]
	return &TokenizerResult{append(a0, append(b.Tokens, a1...)...)}
}
