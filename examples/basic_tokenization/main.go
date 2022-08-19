package main

import (
	"fmt"
	"log"

	"github.com/tpillow/eztok/pkg/eztok"
)

// Custom TokenType defintions.
const (
	TokenTypeEqual         eztok.TokenType = "="
	TokenTypeSemicolon     eztok.TokenType = ";"
	TokenTypeKeywordFloat  eztok.TokenType = "keyword float"
	TokenTypeKeywordString eztok.TokenType = "keyword string"
	TokenTypeKeywordBool   eztok.TokenType = "keyword bool"
)

func main() {
	// Create our tokenizer with the desired nodes.
	// Order of the nodes matters for an InOrderNodeTokenizer, since the first
	// node that is able to parse the next rune will consume & parse the next rune.
	tokenizer := eztok.NewInOrderNodeTokenizer(
		// Skips all whitespace.
		eztok.SkipWhitespaceNode,
		// These match on a specific rune (character), returning a token of the provided type.
		eztok.NewRuneMatchNode(TokenTypeEqual, '='),
		eztok.NewRuneMatchNode(TokenTypeSemicolon, ';'),
		// Matches integer and float numbers.
		eztok.NumberNode,
		// Matches a double-quoted string, allowing for escaped characters.
		eztok.DoubleQuotedEscapedStringNode,
		// These match on a full word, returning a token of the provided type.
		eztok.NewStringMatchNode(TokenTypeKeywordFloat, "float"),
		eztok.NewStringMatchNode(TokenTypeKeywordString, "string"),
		eztok.NewStringMatchNode(TokenTypeKeywordBool, "bool"),
		// Matches a C-style identifier (alpha-numeric and '_', must not start with a digit)
		// This needs to go after any StringMatchNode nodes, since it will match any word
		// following the rules just described.
		eztok.IdentifierNode,
	)

	// Tokenize an example string using our tokenizer above.
	tokens, err := eztok.TokenizeString(tokenizer, `
			float apples = -55.3;
			string title = "Post\tCard";
			bool b = true;
		`)
	// Check for tokenization errors.
	if err != nil {
		log.Fatalf("Error while tokenizing: %v", err)
	}

	// Create a traverser to assist in processing each token in-order.
	traverser := eztok.NewTokenTraverser(tokens)
	// Print out each token in-order.
	fmt.Printf("Tokenized Tokens:\n")
	for traverser.PeekToken(0) != nil {
		fmt.Printf("%v\n", traverser.NextToken().ToString())
	}
}
