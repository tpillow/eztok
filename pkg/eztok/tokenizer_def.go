package eztok

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Defines how a tokenizer should look at runes.
type TokenizerDef struct {
	// All TokenizerNode instances that can process runes for this tokenizer.
	// These are checked in-order from first element (index 0) to last element.
	// The first TokenizerNode whose Peeker returns true for a rune will process the rune.
	TokenizerNodes []*TokenizerNode // TODO: check sorting / iterating happens in order
}

// Creates a new TokenizerDef with the provided TokenizerNode instances.
func NewTokenizerDef(tokenizerNodes ...*TokenizerNode) *TokenizerDef {
	return &TokenizerDef{
		TokenizerNodes: tokenizerNodes,
	}
}

// Using the provided TokenizerCtx, this will tokenize all runes until end of input.
// A TokenizerResult will be returned if tokenization succeeded, else it will be nil
// and an error will be returned.
func (td *TokenizerDef) TokenizeWithCtx(ctx *TokenizerCtx) (*TokenizerResult, error) {
	toks := []*Token{}
	for !ctx.PeekIs(EOFRune) {
		valid := false
		for _, node := range td.TokenizerNodes {
			if node.Peeker(ctx) {
				beforeOriginInfo := ctx.CurOriginInfo
				tok, err := node.Parser(ctx)
				if err != nil {
					return nil, fmt.Errorf("%v at %v", err, ctx.CurOriginInfo.ToString())
				}
				// Some functions may want to discard tokens (ex: comments) & return nil discard
				if tok != nil {
					if tok.OriginInfo == nil {
						tok.OriginInfo = beforeOriginInfo
					}
					toks = append(toks, tok)
				}
				valid = true
				break
			}
		}
		if !valid {
			return nil, fmt.Errorf("unexpected character '%c' at %v",
				ctx.Peek(), ctx.CurOriginInfo.ToString())
		}
	}
	return &TokenizerResult{toks}, nil
}

// Calls TokenizeWithCtx with a TokenizerCtx created from a strings.Reader.
// Use this to tokenize directly from a string. The origin name will be StringOrigin.
func (td *TokenizerDef) TokenizeStr(content string) (*TokenizerResult, error) {
	return td.TokenizeWithCtx(NewTokenizerCtx(
		bufio.NewReader(strings.NewReader(content)), StringOrigin))
}

// Calls TokenizeWithCtx with a TokenizerCtx created from a os.File from the provided path.
// Use this to tokenize a file's contents. The origin name will be the path.
func (td *TokenizerDef) TokenizeFile(path string) (*TokenizerResult, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return td.TokenizeWithCtx(NewTokenizerCtx(
		bufio.NewReader(file), path))
}
