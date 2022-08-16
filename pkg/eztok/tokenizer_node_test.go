package eztok_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tpillow/eztok/pkg/eztok"
	"github.com/tpillow/eztok/pkg/eztok/basic"
)

func NewStrOriginToken(tokenType eztok.TokenType, value any,
	lineNum int, colNum int) *eztok.Token {

	tok := eztok.NewToken(tokenType, value)
	tok.OriginInfo = eztok.NewOriginInfo(eztok.StringOrigin, lineNum, colNum)
	return tok
}

func TokenizeAssertMatches(t *testing.T, content string,
	expTokens []*eztok.Token, nodes ...*eztok.TokenizerNode) {

	res, err := eztok.NewTokenizerDef(nodes...).TokenizeStr(content)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, expTokens, res.Tokens)
}

func TestSingleRuneTokenizerNode(t *testing.T) {
	TokenizeAssertMatches(t, ";; \n ;", []*eztok.Token{
		NewStrOriginToken(TestSemicolonType, ';', 1, 1),
		NewStrOriginToken(TestSemicolonType, ';', 1, 2),
		NewStrOriginToken(TestSemicolonType, ';', 2, 2),
	},
		basic.TokenizerNodeSkipWhitespace,
		eztok.NewSingleRuneTokenizerNode(TestSemicolonType, ';'))
}

func TestSingleRuneTokenizerNodeMulti(t *testing.T) {
	TokenizeAssertMatches(t, ";,; \n ,;", []*eztok.Token{
		NewStrOriginToken(TestSemicolonType, ';', 1, 1),
		NewStrOriginToken(TestCommaType, ',', 1, 2),
		NewStrOriginToken(TestSemicolonType, ';', 1, 3),
		NewStrOriginToken(TestCommaType, ',', 2, 2),
		NewStrOriginToken(TestSemicolonType, ';', 2, 3),
	},
		basic.TokenizerNodeSkipWhitespace,
		eztok.NewSingleRuneTokenizerNode(TestSemicolonType, ';'),
		eztok.NewSingleRuneTokenizerNode(TestCommaType, ','))
}
