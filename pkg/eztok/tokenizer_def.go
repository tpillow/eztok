package eztok

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type TokenizerPeeker func(ctx *TokenizerCtx) bool
type TokenizerParser func(ctx *TokenizerCtx) (*Token, error)

type TokenizerNode struct {
	Peeker TokenizerPeeker
	Parser TokenizerParser
}

func NewTokenizerNode(peeker TokenizerPeeker, parser TokenizerParser) *TokenizerNode {
	return &TokenizerNode{peeker, parser}
}

type TokenizerResult struct {
	Tokens []*Token
}

type TokenizerDef struct {
	TokenizerNodes []*TokenizerNode
}

func NewTokenizerDef(tokenizerNodes ...*TokenizerNode) *TokenizerDef {
	return &TokenizerDef{
		TokenizerNodes: tokenizerNodes,
	}
}

func (td *TokenizerDef) TokenizeWithCtx(ctx *TokenizerCtx) (*TokenizerResult, error) {
	toks := []*Token{}
	for !ctx.PeekIs(EOFRune) {
		for _, node := range td.TokenizerNodes {
			if node.Peeker(ctx) {
				tok, err := node.Parser(ctx)
				if err != nil {
					return nil, fmt.Errorf("%v at %v", err, ctx.AtString())
				}
				// Some functions may want to discard tokens (ex: comments) and will return nil to do so
				if tok != nil {
					toks = append(toks, tok)
				}
			}
		}
	}
	return &TokenizerResult{toks}, nil
}

func (td *TokenizerDef) TokenizeStr(content string) (*TokenizerResult, error) {
	return td.TokenizeWithCtx(NewTokenizerCtx(bufio.NewReader(strings.NewReader(content)), "<string>"))
}

func (td *TokenizerDef) TokenizeFile(path string) (*TokenizerResult, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return td.TokenizeWithCtx(NewTokenizerCtx(bufio.NewReader(file), path))
}
