package main

import (
	"fmt"
	"log"

	"github.com/tpillow/eztok/pkg/eztok"
)

// cat-language
// A mini runnable "language" example using eztok.
// Can print simple numbers and strings
// Can print a cat...sort of
// Can include other pieces of code (like C's "#include")

// Rough EBNF of cat-language:
//
// Program = { IncludeStatement | Statement }
// IncludeStatement = "@include", String
// Statement = Expression, ";"
// Expression = String | Number | "cat"

// Custom TokenType values.
const (
	TokenTypeSemicolon eztok.TokenType = ";"
	TokenTypeInclude   eztok.TokenType = "@include"
	TokenTypeCat       eztok.TokenType = "cat"
)

// Available "files" to include in the example program.
// Key is "file name" and value is "file contents".
var availableIncludePathToContent = map[string]string{
	"Include1": `
		"I'm from Include1!";
		@include "Include2"
	`,
	"Include2": `
		"I'm from Include2!";
		cat;
		cat;
	`,
}

// The main program contents to run.
const ProgramText = `
	"Start of Program";
	44; -3.3;
	@include "Include1"
	0x55;
	"End of Program";
`

// The cat-language tokenizer definition.
// Order of the nodes matters for an InOrderNodeTokenizer, since the first
// node that is able to parse the next rune will consume & parse the next rune.
var tokenizer = eztok.NewInOrderNodeTokenizer(
	// Tokenization will skip/ignore any whitespace.
	eztok.SkipWhitespaceNode,
	// Tokenization will accept integer and float numbers.
	eztok.NumberNode,
	// Tokenization will accept any double-quoted string. Escape characters allowed.
	eztok.DoubleQuotedEscapedStringNode,
	// Tokenization will accept a ';' as its own token.
	eztok.NewRuneMatchNode(TokenTypeSemicolon, ';'),
	// Tokenization will accept '@include' as it's own token.
	eztok.NewStringMatchNode(TokenTypeInclude, "@include"),
	// Tokenization will accept 'cat' as it's own token.
	eztok.NewStringMatchNode(TokenTypeCat, "cat"),
)

// A helper to tokenize a string of text into our cat-language.
func Tokenize(content string) []*eztok.Token {
	toks, err := eztok.TokenizeString(tokenizer, content)
	if err != nil {
		log.Fatalf("Error tokenizing: %v.", err)
	}
	return toks
}

// Interpret an expression (just print it).
// An expression is any of: <integer>, <float>, <string>, cat
func RunExpression(trav *eztok.TokenTraverser) {
	// An expression in the cat-language is only ever one token long.
	tok := trav.NextToken()
	if tok == nil {
		log.Fatalf("Expected an expression but got end of input.")
	}
	// Print the expression (i.e. token's value).
	switch tok.TokenType {
	case eztok.TokenTypeFloat, eztok.TokenTypeInteger, eztok.TokenTypeString:
		fmt.Printf("Expression: %v\n", tok.Value)
	case TokenTypeCat:
		fmt.Printf("Expression: (^-^)\n")
	default:
		log.Fatalf("Expected an expression but got token %v.", tok.ToString())
	}
}

// Interpret a statement.
// A statement is in the form: <Expression> ;
func RunStatement(trav *eztok.TokenTraverser) {
	// Simulate the expression
	RunExpression(trav)
	// Each expression must end with a semicolon.
	eztok.ExpectTokenTypeIs(trav, TokenTypeSemicolon)
}

// Interpret an include statement.
// An include statement is in the form: @include "pathOfCodeToInclude"
func RunIncludeStatement(trav *eztok.TokenTraverser) {
	eztok.ExpectTokenTypeIs(trav, TokenTypeInclude)
	// Get the include path.
	includeStrTok, err := eztok.ExpectTokenTypeIs(trav, eztok.TokenTypeString)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	includePath := includeStrTok.Value.(string)
	// Get the text to include.
	includeText, ok := availableIncludePathToContent[includePath]
	if !ok {
		log.Fatalf("Unknown include path '%v'.", includePath)
	}
	// Tokenize the text to include
	includeToks := Tokenize(includeText)
	// Prepend all new tokens from the include to the existing tokens to traverse.
	// The traverser's next token is always at index 0, so prepending will ensure the
	// included tokens are next to process.
	trav.Tokens = append(includeToks, trav.Tokens...)
}

// Interpret a program.
// A program is 0 or more <IncludeStatement> or <Statement>.
func RunProgram(toks []*eztok.Token) {
	// Create a token traverser to go over all tokens.
	trav := eztok.NewTokenTraverser(toks)
	// For as long as we have tokens, we must have a statement.
	for trav.PeekToken(0) != nil {
		// Check if we might have an include statement.
		if eztok.PeekTokenTypeIs(trav, TokenTypeInclude) {
			// We have an include statement.
			RunIncludeStatement(trav)
		} else {
			// Otherwise, we have a normal statement.
			RunStatement(trav)
		}
	}
}

// Entry point.
func main() {
	// Tokenize the initial program text
	toks := Tokenize(ProgramText)
	// Interpret the program
	RunProgram(toks)
}
