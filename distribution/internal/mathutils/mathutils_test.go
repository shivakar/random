package mathutils_test

import (
	"math"
	"testing"

	"github.com/shivakar/random/distribution/internal/mathutils"
	"github.com/stretchr/testify/assert"
)

func Test_PolyEval(t *testing.T) {
	assert := assert.New(t)

	/* The following examples are used as test:
		f(x) = -2x^2 + 5x - 7; f(3) == -10
	    f(x) = 3x^3 - 4x^2 + 2x + 9; f(-2) = -35
	    f(x) = 4x^3 - 3x^2 - 5x + 7; f(-3) = -113
	    f(x) = -3x^4 + 9x^2 - 5x - 1; f(-1) = 10
	    Equations are written out this way since the function
	    needs the highest degree coefficient first
	*/
	data := []struct {
		x    float64
		coef []float64
		fx   float64
	}{
		{3, []float64{-2, 5, -7}, -10},
		{-2, []float64{3, -4, 2, 9}, -35},
		{-3, []float64{4, -3, -5, 7}, -113},
		{-1, []float64{-3, 0, 9, -5, -1}, 10},
	}

	for _, rec := range data {
		assert.Equal(rec.fx, mathutils.PolyEval(rec.x, rec.coef))
	}
}

func Test_NormICDF(t *testing.T) {
	assert := assert.New(t)
	data := []struct {
		x  float64
		ix float64
	}{
		{0.962746, 1.783485},
		{0.665312, 0.427005},
		{0.876605, 1.158179},
		{0.169157, -0.957500},
		{0.613341, 0.288038},
		{0.218354, -0.777765},
		{0.695799, 0.512357},
		{0.618395, 0.301268},
		{0.634381, 0.343479},
		{0.167280, -0.964968},
		{0, math.SmallestNonzeroFloat64},
		{1, math.MaxFloat64},
		{1.2664165549e-14, -7.620200},
	}

	for _, v := range data {
		r := mathutils.NormICDF(v.x)
		assert.InDelta(v.ix, r, 1e-5)
	}
}
