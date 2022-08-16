package basic

import (
	"fmt"
	"log"
	"strconv"
	"unicode"

	"github.com/tpillow/eztok/pkg/eztok"
)

const (
	TokenTypeStr   eztok.TokenType = "String"
	TokenTypeInt   eztok.TokenType = "Integer"
	TokenTypeFloat eztok.TokenType = "Float"
	TokenTypeIdent eztok.TokenType = "Identifier"
)

var TokenizerNodeSkipWhitespace = eztok.NewTokenizerNode(
	func(ctx *eztok.TokenizerCtx) bool {
		return unicode.IsSpace(ctx.Peek())
	}, func(ctx *eztok.TokenizerCtx) (*eztok.Token, error) {
		r := ctx.Next()
		if !unicode.IsSpace(r) {
			log.Panic("Whitespace is always expected in whitespace parse")
		}
		return nil, nil
	})

// TODO: need to ensure it doesn't overlap with float
var TokenizerNodeInt = eztok.NewTokenizerNode(
	func(ctx *eztok.TokenizerCtx) bool {
		r := ctx.Peek()
		rStr := string(r)
		return unicode.IsDigit(r) || rStr == "-" || rStr == "+"
	}, func(ctx *eztok.TokenizerCtx) (*eztok.Token, error) {
		first := ctx.Next()
		str := ""
		base := 0
		// Check for special prefix cases (not using built-in with ParseInt)
		if first == '0' && ctx.PeekIs('b') {
			// Binary number disallow +/- prefix: 0b---
			ctx.AssertNextIs('b')
			str = ctx.ReadUntilNot(unicode.IsDigit)
			base = 2
		} else if first == '0' && ctx.PeekIs('o') {
			// Octal number disallow +/- prefix: 0o---
			ctx.AssertNextIs('o')
			str = ctx.ReadUntilNot(unicode.IsDigit)
			base = 8
		} else if first == '0' && ctx.PeekIs('x') {
			// Hex number disallow +/- prefix: 0x---
			ctx.AssertNextIs('x')
			str = ctx.ReadUntilNot(unicode.IsDigit)
			base = 16
		} else {
			// Base 10 integer, optional +/- prefix
			str = string(first) + ctx.ReadUntilNot(unicode.IsDigit)
			base = 10
		}
		if base <= 0 {
			log.Panicf("Invalid integer base value: %v", base)
		}

		intVal, err := strconv.ParseInt(str, base, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer literal '%v': %v", str, err)
		}
		return eztok.NewToken(TokenTypeInt, intVal), nil
	})

var TokenizerNodeFloat = eztok.NewTokenizerNode(
	func(ctx *eztok.TokenizerCtx) bool {
		r := ctx.Peek()
		rStr := string(r)
		return unicode.IsDigit(r) || rStr == "-" || rStr == "+"
	}, func(ctx *eztok.TokenizerCtx) (*eztok.Token, error) {
		str := string(ctx.Next()) + ctx.ReadUntil(func(r rune) bool {
			return !unicode.IsDigit(r) && string(r) != "."
		})

		floatVal, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid float literal '%v': %v", str, err)
		}
		return eztok.NewToken(TokenTypeInt, floatVal), nil
	})

var TokenizerNodeCIdentifier = eztok.NewTokenizerNode(
	func(ctx *eztok.TokenizerCtx) bool {
		r := ctx.Peek()
		return unicode.IsLetter(r) || string(r) == "_"
	}, func(ctx *eztok.TokenizerCtx) (*eztok.Token, error) {
		str := ctx.ReadUntil(func(r rune) bool {
			return !(unicode.IsDigit(r) || unicode.IsLetter(r) || string(r) == "_")
		})
		return eztok.NewToken(TokenTypeIdent, str), nil
	})

var TokenizerNodeEscapedDoubleQuotedStr = eztok.NewTokenizerNode(
	func(ctx *eztok.TokenizerCtx) bool {
		return ctx.PeekIs('"')
	}, func(ctx *eztok.TokenizerCtx) (*eztok.Token, error) {
		ctx.AssertNextIs('"')
		// TODO: string escapes
		unescapedStr := ctx.ReadUntil(func(r rune) bool {
			return r == '"'
		})
		if !ctx.NextIs('"') {
			return nil, fmt.Errorf("unterminated string '%v'", unescapedStr)
		}
		return eztok.NewToken(TokenTypeStr, unescapedStr), nil
	})
