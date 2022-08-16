package eztok_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tpillow/eztok/pkg/eztok"
)

func TestTokenizeStrNoNodesEmpty(t *testing.T) {
	res, err := eztok.NewTokenizerDef().TokenizeStr("")
	assert.Nil(t, err)
	assert.Len(t, res.Tokens, 0)
}

func TestTokenizeStrNoNodesNotEmpty(t *testing.T) {
	res, err := eztok.NewTokenizerDef().TokenizeStr("hello")
	assert.Equal(t, err,
		fmt.Errorf("unexpected character 'h' at %v:1:1", eztok.StringOrigin))
	assert.Nil(t, res)
}
