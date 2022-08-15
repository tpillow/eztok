package eztok

import (
	"fmt"
	"log"
	"strconv"
	"unicode"
)

const (
	BasicTokenTypeStr       TokenType = "String"
	BasicTokenTypeInt       TokenType = "Integer"
	BasicTokenTypeIdent     TokenType = "Identifier"
	BasicTokenTypeComma     TokenType = ","
	BasicTokenTypeSemicolon TokenType = ";"
)

func BasicTokenizerIsPossibleNumericRune(r rune) bool {
	rStr := string(r)
	return unicode.IsDigit(r) || rStr == "." || rStr == "_" || rStr == "-"
}

func BasicTokenizerIsIdentRune(c rune) bool {
	return unicode.IsDigit(c) || unicode.IsLetter(c) || string(c) == "_"
}

func NewSingleRuneTokenizerNode(r rune, tokenType TokenType) *TokenizerNode {
	return NewTokenizerNode(
		func(ctx *TokenizerCtx) bool {
			return ctx.PeekIs(r)
		}, func(ctx *TokenizerCtx) (*Token, error) {
			ctx.AssertNextIs(r)
			return NewToken(tokenType, r), nil
		})
}

var BasicTokenizerNodeSemicolon = NewSingleRuneTokenizerNode(';', BasicTokenTypeSemicolon)

var BasicTokenizerNodeComma = NewSingleRuneTokenizerNode(';', BasicTokenTypeComma)

var BasicTokenizerNodeInt = NewTokenizerNode(
	func(ctx *TokenizerCtx) bool {
		r := ctx.Peek()
		return unicode.IsDigit(r) || string(r) == "-"
	}, func(ctx *TokenizerCtx) (*Token, error) {
		str := ctx.ReadUntilNot(BasicTokenizerIsPossibleNumericRune)
		if len(str) <= 0 {
			log.Panicf("Length of int string should never be 0")
		}
		intVal, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		return NewToken(BasicTokenTypeInt, intVal), nil
	})

var BasicTokenizerNodeIdent = NewTokenizerNode(
	func(ctx *TokenizerCtx) bool {
		r := ctx.Peek()
		return !unicode.IsDigit(r) && BasicTokenizerIsIdentRune(r)
	}, func(ctx *TokenizerCtx) (*Token, error) {
		str := ctx.ReadUntilNot(BasicTokenizerIsIdentRune)
		if len(str) <= 0 {
			log.Panicf("Length of ident string should never be 0")
		}
		return NewToken(BasicTokenTypeIdent, str), nil
	})

var BasicTokenizerNodeStr = NewTokenizerNode(
	func(ctx *TokenizerCtx) bool {
		return ctx.PeekIs('"')
	}, func(ctx *TokenizerCtx) (*Token, error) {
		ctx.AssertNextIs('"')
		// TODO: string escapes
		unescapedStr := ctx.ReadUntil(func(r rune) bool {
			return r == '"'
		})
		if !ctx.NextIs('"') {
			return nil, fmt.Errorf("unterminated string '%v'", unescapedStr)
		}
		return NewToken(BasicTokenTypeStr, unescapedStr), nil
	})
