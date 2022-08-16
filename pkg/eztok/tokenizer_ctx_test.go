package eztok_test

import (
	"bufio"
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
	"github.com/tpillow/eztok/pkg/eztok"
)

func NewCtxFromStr(content string) *eztok.TokenizerCtx {
	return eztok.NewTokenizerCtx(bufio.NewReader(strings.NewReader(content)), eztok.StringOrigin)
}

func TestTokenizerCtxEmptyStr(t *testing.T) {
	ctx := NewCtxFromStr("")
	origin := eztok.NewOriginInfo(eztok.StringOrigin, 1, 1)
	assert.Equal(t, origin, ctx.CurOriginInfo)
	assert.Equal(t, eztok.EOFRune, ctx.Peek())
	assert.Equal(t, eztok.EOFRune, ctx.Next())
	assert.Equal(t, "", ctx.ReadUntil(func(r rune) bool { return false }))
	assert.Equal(t, "", ctx.ReadUntilNot(func(r rune) bool { return true }))
	assert.Equal(t, origin, ctx.CurOriginInfo)
}

func TestTokenizerCtxPeekNext(t *testing.T) {
	testStr := "hello\n world!\n\n"
	ctx := NewCtxFromStr(testStr)
	expOrigin := eztok.NewOriginInfo(eztok.StringOrigin, 1, 1)
	for idx, testR := range testStr {
		assert.Equal(t, testR, ctx.Peek())
		assert.True(t, ctx.PeekIs(testR))
		assert.Equal(t, expOrigin, ctx.CurOriginInfo)

		if idx < len(testStr)/2 {
			assert.Equal(t, testR, ctx.Next())
		} else {
			assert.True(t, ctx.NextIs(testR))
		}

		if testR == '\n' {
			expOrigin.LineNum++
			expOrigin.ColNum = 1
		} else {
			expOrigin.ColNum++
		}
		assert.Equal(t, expOrigin, ctx.CurOriginInfo)
	}
	assert.True(t, ctx.PeekIs(eztok.EOFRune))
	assert.Equal(t, expOrigin, ctx.CurOriginInfo)

	assert.Panics(t, func() {
		ctx.AssertNextIs(';')
	})
}

func TestTokenizerCtxReadUntil(t *testing.T) {
	ctx := NewCtxFromStr("hello, world!")
	assert.Equal(t, "hello,", ctx.ReadUntil(func(r rune) bool {
		return unicode.IsSpace(r)
	}))
	assert.True(t, ctx.NextIs(' '))
	assert.Equal(t, "world!", ctx.ReadUntil(func(r rune) bool {
		return unicode.IsSpace(r)
	}))
}

func TestTokenizerCtxReadUntilNot(t *testing.T) {
	ctx := NewCtxFromStr("hello, world!")
	assert.Equal(t, "hello,", ctx.ReadUntilNot(func(r rune) bool {
		return !unicode.IsSpace(r)
	}))
	assert.True(t, ctx.NextIs(' '))
	assert.Equal(t, "world!", ctx.ReadUntilNot(func(r rune) bool {
		return !unicode.IsSpace(r)
	}))
}
