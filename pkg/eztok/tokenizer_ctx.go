package eztok

import (
	"bufio"
	"log"
	"strings"
)

const (
	// A rune representing the end of input.
	EOFRune rune = 0
	// A rune representing a newline that will increment the line number.
	NewlineRune rune = '\n'
)

// Callback function used for some TokenizerCtx methods. See ReadUntil and ReadUntilNot.
type ReadUntilFunc func(r rune) bool

// Holds information required to tokenize input from a bufio.Reader.
type TokenizerCtx struct {
	// The bufio.Reader of input to tokenize.
	Reader *bufio.Reader
	// The next rune's OriginInfo.
	CurOriginInfo *OriginInfo
}

// Creates a new TokenizerCtx with input from the provided reader.
// The originName is passed along to all OriginInfo.
func NewTokenizerCtx(reader *bufio.Reader, originName string) *TokenizerCtx {
	return &TokenizerCtx{
		Reader:        reader,
		CurOriginInfo: NewOriginInfo(originName, 1, 1),
	}
}

// Returns the next rune to process without consuming it.
// Returns EOFRune if the end of input has been reached.
func (ctx *TokenizerCtx) Peek() rune {
	r, size, err := ctx.Reader.ReadRune()
	if err != nil || size <= 0 {
		return EOFRune
	}
	if err := ctx.Reader.UnreadRune(); err != nil {
		log.Panicf("UnreadRune unexpectedly failed: %v (at %v)", err, ctx.CurOriginInfo.ToString())
	}
	return r
}

// Returns the next rune to process and consumes it, and updates CurOriginInfo.
// Returns EOFRune if the end of input has been reached.
func (ctx *TokenizerCtx) Next() rune {
	r, size, err := ctx.Reader.ReadRune()
	if err != nil || size <= 0 {
		return EOFRune
	}
	// Re-make the origin infos, since things will be storing pointers to them
	// TODO: better way to do this than re-instantiating this a bunch
	if r == NewlineRune {
		ctx.CurOriginInfo = NewOriginInfo(ctx.CurOriginInfo.Origin,
			ctx.CurOriginInfo.LineNum+1, 1)
	} else {
		ctx.CurOriginInfo = NewOriginInfo(ctx.CurOriginInfo.Origin,
			ctx.CurOriginInfo.LineNum, ctx.CurOriginInfo.ColNum+1)
	}
	return r
}

// Analagous to: Peek() == exp
func (ctx *TokenizerCtx) PeekIs(exp rune) bool {
	return ctx.Peek() == exp
}

// Analagous to: Next() == exp
func (ctx *TokenizerCtx) NextIs(exp rune) bool {
	return ctx.Next() == exp
}

// Calls Next() and panics if the rune returned is not exp.
func (ctx *TokenizerCtx) AssertNextIs(exp rune) {
	val := ctx.Next()
	if val != exp {
		log.Panicf("Expected '%v' but got '%v' at %v",
			exp, val, ctx.CurOriginInfo.ToString())
	}
}

// Reads runes into a string until untilFunc returns true.
func (ctx *TokenizerCtx) ReadUntil(untilFunc ReadUntilFunc) string {
	var sb strings.Builder
	for r := ctx.Peek(); r != EOFRune && !untilFunc(r); r = ctx.Peek() {
		sb.WriteRune(ctx.Next())
	}
	return sb.String()
}

// Reads runes into a string until untilFunc returns false.
func (ctx *TokenizerCtx) ReadUntilNot(untilFunc ReadUntilFunc) string {
	return ctx.ReadUntil(func(r rune) bool { return !untilFunc(r) })
}
