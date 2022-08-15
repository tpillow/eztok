package eztok

import (
	"log"
	"strconv"
	"unicode"
)

const (
	BasicTokenTypeStr       TokenType = "String"
	BasicTokenTypeInt       TokenType = "Integer"
	BasicTokenTypeBool      TokenType = "Boolean"
	BasicTokenTypeIdent     TokenType = "Identifier"
	BasicTokenTypeComma     TokenType = ","
	BasicTokenTypeSemicolon TokenType = ";"
)

var SingleRuneBasicBreakerSet = map[string]bool{
	string(BasicTokenTypeComma):     true,
	string(BasicTokenTypeSemicolon): true,
	"\"":                            true,
}

func BasicTokenizerIsBreaker(r rune) bool {
	_, isBreaker := SingleRuneBasicBreakerSet[string(r)]
	return isBreaker || unicode.IsSpace(r)
}

func BasicTokenizerIsIdentRune(c rune) bool {
	return unicode.IsDigit(c) || unicode.IsLetter(c) || string(c) == "_"
}

var BasicTokenizerNodeInt = NewTokenizerNode(
	func(ctx *TokenizerCtx) bool {
		r := ctx.Peek()
		return unicode.IsDigit(r) || string(r) == "-"
	}, func(ctx *TokenizerCtx) (*Token, error) {
		str := ctx.ReadUntil(BasicTokenizerIsBreaker)
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
