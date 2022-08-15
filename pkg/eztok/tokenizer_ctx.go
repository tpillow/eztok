package eztok

import (
	"bufio"
	"fmt"
	"log"
	"strings"
)

const EOFRune rune = 0

type ReadUntilFunc func(r rune) bool

type TokenizerCtx struct {
	Reader  *bufio.Reader
	Origin  string
	LineNum int
	ColNum  int
}

func NewTokenizerCtx(reader *bufio.Reader, origin string) *TokenizerCtx {
	return &TokenizerCtx{
		Reader:  reader,
		Origin:  origin,
		LineNum: 1,
		ColNum:  0,
	}
}

func (ctx *TokenizerCtx) AtString() string {
	return fmt.Sprintf("%v:%v:%v", ctx.Origin, ctx.LineNum, ctx.ColNum)
}

func (ctx *TokenizerCtx) Peek() rune {
	r, _, err := ctx.Reader.ReadRune()
	if err != nil {
		return EOFRune
	}
	if err := ctx.Reader.UnreadRune(); err != nil {
		log.Panicf("UnreadRune unexpectedly failed: %v (at %v)", err, ctx.AtString())
	}
	return r
}

func (ctx *TokenizerCtx) Next() rune {
	r, _, err := ctx.Reader.ReadRune()
	if err != nil {
		return EOFRune
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

func (ctx *TokenizerCtx) AssertNextIs(exp rune) {
	val := ctx.Next()
	if val != exp {
		log.Panicf("Expected character '%v' but got '%v' at %v", exp, val, ctx.AtString())
	}
}

func (ctx *TokenizerCtx) ReadUntil(untilFunc ReadUntilFunc) string {
	var sb strings.Builder
	for r := ctx.Peek(); r != EOFRune && !untilFunc(r); {
		sb.WriteRune(ctx.Next())
	}
	return sb.String()
}

func (ctx *TokenizerCtx) ReadUntilNot(untilFunc ReadUntilFunc) string {
	return ctx.ReadUntil(func(r rune) bool { return !untilFunc(r) })
}
