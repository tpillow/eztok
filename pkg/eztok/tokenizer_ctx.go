package eztok

import (
	"bufio"
	"fmt"
	"log"
	"strings"
)

const EOF_RUNE rune = 0

type ReadUntilFunc func(r rune) bool

type TokenizerCtx struct {
	Reader  *bufio.Reader
	Origin  string
	LineNum int
	ColNum  int
}

func NewTokenizerCtx(reader *bufio.Reader, origin string) *TokenizerCtx {
	ctx := &TokenizerCtx{
		Reader:  reader,
		Origin:  origin,
		LineNum: 1,
		ColNum:  0,
	}
	ctx.Peek() // Set EOF appropriately
	return ctx
}

func (ctx *TokenizerCtx) AtString() string {
	return fmt.Sprintf("%v:%v:%v", ctx.Origin, ctx.LineNum, ctx.ColNum)
}

func (ctx *TokenizerCtx) Peek() rune {
	r, _, err := ctx.Reader.ReadRune()
	if err != nil {
		return EOF_RUNE
	}
	if err := ctx.Reader.UnreadRune(); err != nil {
		log.Panicf("UnreadRune unexpectedly failed: %v (at %v)", err, ctx.AtString())
	}
	return r
}

func (ctx *TokenizerCtx) Next() rune {
	r, _, err := ctx.Reader.ReadRune()
	if err != nil {
		return EOF_RUNE
	}
	return r
}

func (ctx *TokenizerCtx) PeekIs(exp rune) bool {
	return ctx.Peek() == exp
}

func (ctx *TokenizerCtx) NextIs(exp rune) bool {
	val := ctx.Next()
	return val == exp
}

func (ctx *TokenizerCtx) ReadUntil(untilFunc ReadUntilFunc) string {
	var sb strings.Builder
	for r := ctx.Peek(); r != EOF_RUNE && !untilFunc(r); {
		sb.WriteRune(ctx.Next())
	}
	return sb.String()
}

func (ctx *TokenizerCtx) ReadUntilNot(untilFunc ReadUntilFunc) string {
	return ctx.ReadUntil(func(r rune) bool { return !untilFunc(r) })
}
