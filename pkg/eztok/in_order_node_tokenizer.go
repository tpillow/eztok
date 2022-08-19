package eztok

import "fmt"

// A Tokenizer that operates on a slice of Node objects, attempting to
// parse Token objects by checking each Node in-order.
type InOrderNodeTokenizer struct {
	// The slice of Node objects to operate on. A Node at a lower index
	// will attempt to be processed before a Node with a higher index.
	Nodes []Node
}

// Returns a new InOrderNodeTokenizer with the given parameters.
func NewInOrderNodeTokenizer(initialNodes ...Node) *InOrderNodeTokenizer {
	return &InOrderNodeTokenizer{initialNodes}
}

// For as long as Context.PeekRune(0) does not return NilRune, all
// InOrderNodeTokenzer.Nodes will have their CanParseToken function called
// until one returns true for the current Context state. In which case, that
// Node will have its ParseToken function called with the current Context state
// and no further nodes will be checked until the next iteration.
func (tizer *InOrderNodeTokenizer) Tokenize(ctx Context) ([]*Token, error) {
	toks := []*Token{}
	for ctx.PeekRune(0) != NilRune {
		valid := false
		for _, node := range tizer.Nodes {
			if node.CanParseToken(ctx) {
				beforeOriginInfo := ctx.GetNextOrigin()
				tok, err := node.ParseToken(ctx)
				if err != nil {
					return nil, fmt.Errorf("%v at %v", err, ctx.GetNextOrigin().ToString())
				}
				if tok != nil {
					if tok.Origin == nil {
						tok.Origin = beforeOriginInfo
					}
					toks = append(toks, tok)
				}
				valid = true
				break
			}
		}
		if !valid {
			return nil, fmt.Errorf("unexpected rune '%c' at %v",
				ctx.PeekRune(0), ctx.GetNextOrigin().ToString())
		}
	}
	return toks, nil
}
