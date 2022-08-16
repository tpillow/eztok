package eztok_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tpillow/eztok/pkg/eztok"
)

const (
	TestType0         eztok.TokenType = "type0"
	TestType1         eztok.TokenType = "type1"
	TestSemicolonType eztok.TokenType = ";"
	TestCommaType     eztok.TokenType = ","
)

func TestTokenEquality(t *testing.T) {
	tok0a := eztok.NewToken(TestType0, 5)
	tok0b := eztok.NewToken(TestType0, 5)
	tok1 := eztok.NewToken(TestType1, 5)

	assert.Equal(t, &eztok.Token{TestType0, 5, nil}, tok0a)
	assert.Equal(t, tok0a, tok0b)
	assert.Equal(t, &eztok.Token{TestType1, 5, nil}, tok1)
	assert.NotEqual(t, tok0a, tok1)
}
