package eztok

import (
	"bufio"
	"log"
)

// A Context whose input rune stream is a bufio.Reader.
type ReaderContext struct {
	reader      *bufio.Reader
	originName  string
	nextLineNum int
	nextColNum  int
	runeQueue   []rune
}

// Returns a new ReaderContext with the provided parameters. The originName will
// be used as the Origin.Name for all runes.
func NewReaderContext(reader *bufio.Reader, originName string) *ReaderContext {
	return &ReaderContext{reader, originName, 1, 1, []rune{}}
}

// Return the rune that is relative runes ahead of the current
// rune in the input. Returns NilRune if there is none.
func (ctx *ReaderContext) PeekRune(relative int) rune {
	if relative < 0 {
		log.Panicf("PeekRuneAhead cannot peek negatively; tried peeking a relative '%v' runes", relative)
	}
	if relative < len(ctx.runeQueue) {
		return ctx.runeQueue[relative]
	} else if len(ctx.runeQueue) > 0 && ctx.runeQueue[len(ctx.runeQueue)-1] == NilRune {
		return NilRune
	}

	for r, size, err := ctx.reader.ReadRune(); err == nil && size > 0; r, size, err = ctx.reader.ReadRune() {
		ctx.runeQueue = append(ctx.runeQueue, r)
		if relative < len(ctx.runeQueue) {
			return ctx.runeQueue[relative]
		}
	}
	ctx.runeQueue = append(ctx.runeQueue, NilRune)
	return NilRune
}

// Consume (i.e. advance the input stream by 1 rune) and return the
// consumed rune. Returns NilRune if there is none.
func (ctx *ReaderContext) NextRune() rune {
	if len(ctx.runeQueue) <= 0 {
		ctx.PeekRune(0)
		if len(ctx.runeQueue) <= 0 {
			log.Panicf("NextRune expects len(runeQueue) > 0 after a call to PeekRune always")
		}
	}

	r := ctx.runeQueue[0]
	if r == NilRune {
		return r
	}
	ctx.runeQueue = ctx.runeQueue[1:]

	if r == '\n' {
		ctx.nextLineNum++
		ctx.nextColNum = 1
	} else {
		ctx.nextColNum++
	}
	return r
}

// Returns the Origin information of the rune that would be returned
// by a call to NextRune().
func (ctx *ReaderContext) GetNextOrigin() *Origin {
	return NewOrigin(ctx.originName, ctx.nextLineNum, ctx.nextColNum)
}
