package eztok_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tpillow/eztok/pkg/eztok"
)

func TestMergeTokenizerResults(t *testing.T) {
	t0 := eztok.NewToken(TestType0, 5)
	t1 := eztok.NewToken(TestType1, 5)
	r0 := &eztok.TokenizerResult{[]*eztok.Token{t0}}
	r1 := &eztok.TokenizerResult{[]*eztok.Token{t1}}
	assert.Equal(t, &eztok.TokenizerResult{[]*eztok.Token{t0, t1}},
		eztok.MergeTokenizerResults(r0, r1))
}

func TestMergeTokenizerResultsAt(t *testing.T) {
	t0 := eztok.NewToken(TestType0, 5)
	t1 := eztok.NewToken(TestType1, 5)
	r0 := &eztok.TokenizerResult{[]*eztok.Token{t0, t0}}
	r1 := &eztok.TokenizerResult{[]*eztok.Token{t1}}
	assert.Equal(t, &eztok.TokenizerResult{[]*eztok.Token{t0, t1, t0}},
		eztok.MergeTokenizerResultsAt(r0, 1, r1))
}
