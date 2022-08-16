package main

import (
	"fmt"
	"log"

	"github.com/tpillow/eztok/pkg/eztok"
	"github.com/tpillow/eztok/pkg/eztok/basic"
	"github.com/tpillow/eztok/pkg/eztok/traverse"
)

const (
	TokenTypeEqual     eztok.TokenType = "="
	TokenTypeSemicolon eztok.TokenType = ";"
)

func TraverserErrorHandler(err error) {
	log.Fatalf("Error traversing: %v", err)
}

func main() {
	tokenizerDef := eztok.NewTokenizerDef(
		basic.TokenizerNodeSkipWhitespace,
		basic.TokenizerNodeCIdentifier,
		basic.TokenizerNodeInt,
		basic.TokenizerNodeFloat,
		eztok.NewSingleRuneTokenizerNode(TokenTypeEqual, '='),
		eztok.NewSingleRuneTokenizerNode(TokenTypeSemicolon, ';'),
	)
	res, err := tokenizerDef.TokenizeStr("int apples = 0x5050;")
	if err != nil {
		log.Panicf("Error tokenizing: %v", err)
	}

	traverser := traverse.NewTokenTraverser(res, TraverserErrorHandler)
	for !traverser.IsAtEnd() {
		fmt.Printf("%v\n", traverser.Next().ToString())
	}
}
