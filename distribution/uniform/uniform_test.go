package uniform_test

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/shivakar/random/distribution/internal/distribtest"
	"github.com/shivakar/random/distribution/uniform"
	"github.com/shivakar/random/prng"
	"github.com/shivakar/random/prng/mt19937"
	"github.com/shivakar/random/prng/splitmix64"
	"github.com/shivakar/random/prng/xoroshiro128plus"
	"github.com/shivakar/random/prng/xorshift1024star"
	"github.com/shivakar/random/prng/xorshift128plus"
	"github.com/stretchr/testify/assert"
)

func Test_Uniform_GetParams(t *testing.T) {
	assert := assert.New(t)
	data := [...]struct {
		a float64
		b float64
	}{
		{0.0, 1.0},
		{3.0, 5.0},
		{-1.0, 4.0},
	}

	r := mt19937.New(0)
	for _, rec := range data {
		d := uniform.New(r, rec.a, rec.b)
		params := d.GetParams()
		assert.Equal(rec.a, params[0])
		assert.Equal(rec.b, params[1])
	}
}

func Test_Uniform_CDF(t *testing.T) {
	assert := assert.New(t)
	data := [...]struct {
		a float64
		b float64
		x []float64
		c []float64
	}{
		{0.0, 1.0,
			[]float64{0.642390331345, 0.761060069636, 0.597690381519,
				0.398903465024, 0.934476726782, 0.188859882506,
				0.530799008153, 0.628251167891, -1.0, 2.0},
			[]float64{0.642390331345, 0.761060069636, 0.597690381519,
				0.398903465024, 0.934476726782, 0.188859882506,
				0.530799008153, 0.628251167891, 0.0, 1.0},
		},
		{3.0, 5.0,
			[]float64{3.64277435426, 4.53121866463, 3.78183111983,
				4.469080263, 4.13137381507, 3.17604899425,
				4.26949916613, 3.61177034539, 2.0, 6.0},
			[]float64{0.321387177129, 0.765609332315, 0.390915559913,
				0.734540131502, 0.565686907535, 0.0880244971231,
				0.634749583064, 0.305885172697, 0.0, 1.0},
		},
		{-1.0, 4.0,
			[]float64{-0.142084622793, -0.800853684496, 0.752082935439,
				0.880942211634, 3.43867912977, 1.93094511565,
				-0.439590605221, -0.579982744963, -2.0, 5.0},
			[]float64{0.171583075441, 0.0398292631008, 0.350416587088,
				0.376188442327, 0.887735825955, 0.586189023131,
				0.112081878956, 0.0840034510074, 0.0, 1.0},
		},
		{-4.0, 0.0,
			[]float64{-1.72059755691, -0.573654329105, -3.80530030703,
				-0.699455861662, -3.92525543522, -3.37660839104,
				-0.640607362997, -0.29541193617, -5.0, 1.0},
			[]float64{0.569850610774, 0.856586417724, 0.0486749232425,
				0.825136034585, 0.0186861411962, 0.155847902239,
				0.839848159251, 0.926147015957, 0.0, 1.0},
		},
		{-10.0, -3.0,
			[]float64{-5.14899209812, -5.60614144756, -8.25595182256,
				-4.6540948228, -8.42692523546, -5.25977779523,
				-4.51085880636, -3.28525287482, -11.0, -2.0},
			[]float64{0.69300112884, 0.62769407892, 0.249149739635,
				0.7637007396, 0.224724966363, 0.677174600682,
				0.784163027663, 0.959249589312, 0.0, 1.0},
		},
	}

	for _, rec := range data {
		d := uniform.New(mt19937.New(0), rec.a, rec.b)
		for i := range rec.x {
			assert.InDelta(rec.c[i], d.CDF(rec.x[i]), float64(1e-10))
		}
	}
}

func Test_Uniform_PDF(t *testing.T) {
	assert := assert.New(t)
	data := [...]struct {
		a float64
		b float64
		x []float64
		c []float64
	}{
		{0.0, 1.0,
			[]float64{0.990064972772, 0.747309151714, 0.662051656946,
				0.647578669044, 0.966606825768, 0.850051931344,
				0.773112707833, 0.193343243202, -1.0, 2.0},
			[]float64{1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 0.0, 0.0},
		},
		{3.0, 5.0,
			[]float64{3.65506069727, 4.42938261542, 3.28898727256,
				4.37249149751, 4.86887798859, 3.27634337611,
				4.28289737021, 4.9866123335, 2.0, 6.0},
			[]float64{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.0, 0.0},
		},
		{-1.0, 4.0,
			[]float64{2.59357665745, 3.55069813554, 0.0649731494036,
				-0.308744847397, 0.615539536244, -0.320356574673,
				0.939473422489, 2.85412791568, -2.0, 5.0},
			[]float64{0.2, 0.2, 0.2, 0.2, 0.2, 0.2, 0.2, 0.2, 0.0, 0.0},
		},
		{-4.0, 0.0,
			[]float64{-0.226993339583, -3.41094373091, -3.2156141227,
				-3.2859079379, -0.382120126524, -2.50653616431,
				-1.32464605551, -1.68583420097, -5.0, 1.0},
			[]float64{0.25, 0.25, 0.25, 0.25, 0.25, 0.25, 0.25, 0.25, 0.0, 0.0},
		},
		{-10.0, -3.0,
			[]float64{-4.17595298522, -3.24106802865, -8.25823498894,
				-9.7428294559, -7.76176239276, -7.55696717771,
				-7.75328208605, -3.58267612937, -11.0, -2.0},
			[]float64{0.142857142857, 0.142857142857, 0.142857142857,
				0.142857142857, 0.142857142857, 0.142857142857,
				0.142857142857, 0.142857142857, 0.0, 0.0},
		},
	}

	for _, rec := range data {
		d := uniform.New(mt19937.New(0), rec.a, rec.b)
		for i := range rec.x {
			assert.InDelta(rec.c[i], d.PDF(rec.x[i]), float64(1e-10))
		}
	}
}

func Test_Uniform_Float64(t *testing.T) {
	assert := assert.New(t)
	engines := [...]struct {
		r    prng.Engine
		name string
	}{
		{mt19937.New(0), "mt19937"},
		{splitmix64.New(0), "splitmix64"},
		{xoroshiro128plus.New(0), "tinymt64"},
		{xorshift128plus.New(0), "xorshift128plus"},
		{xorshift1024star.New(0), "xorshift1024star"},
	}
	for _, engine := range engines {
		d := uniform.New(engine.r, 0.0, 1.0)
		_, pval := distribtest.KSTest(d)
		// Significance level of 1%
		assert.True(pval > 0.001, fmt.Sprintf("%s, pval:%.6f", engine.name, pval))
	}
}

// Benchmarks

func Benchmark_BuiltInPRNG_Float64(b *testing.B) {
	d := rand.New(rand.NewSource(time.Now().UnixNano()))
	x := 0.0
	y := 1.0
	for i := 0; i < b.N; i++ {
		_ = x + (y-x)*d.Float64()
	}
}

func Benchmark_Uniform_MT19937_Float64(b *testing.B) {
	d := uniform.New(mt19937.New(0), 0.0, 1.0)
	for i := 0; i < b.N; i++ {
		_ = d.Float64()
	}
}

func Benchmark_Uniform_SplitMix64_Float64(b *testing.B) {
	d := uniform.New(splitmix64.New(0), 0.0, 1.0)
	for i := 0; i < b.N; i++ {
		_ = d.Float64()
	}
}

func Benchmark_Uniform_Xoroshiro128Plus_Float64(b *testing.B) {
	d := uniform.New(xoroshiro128plus.New(0), 0.0, 1.0)
	for i := 0; i < b.N; i++ {
		_ = d.Float64()
	}
}

func Benchmark_Uniform_Xorshift128Plus_Float64(b *testing.B) {
	d := uniform.New(xorshift128plus.New(0), 0.0, 1.0)
	for i := 0; i < b.N; i++ {
		_ = d.Float64()
	}
}

func Benchmark_Uniform_Xorshift1024Star_Float64(b *testing.B) {
	d := uniform.New(xorshift1024star.New(0), 0.0, 1.0)
	for i := 0; i < b.N; i++ {
		_ = d.Float64()
	}
}

// Example - Uniform Distribution
func ExampleUniformDistribution() {
	// Example based on
	// http://www.cplusplus.com/reference/random/normal_distribution/
	r := xorshift128plus.New(20170611)
	d := uniform.New(r, 3.0, 7.0)

	var p [10]int

	nrolls := 10000
	nstars := 100
	for i := 0; i < nrolls; i++ {
		n := d.Float64()
		if n >= 0.0 && n < 10.0 {
			p[int(n)]++
		}
	}

	fmt.Println("Uniform Distribution: a=3.0, sigma=7.0")

	for i := 0; i < 10; i++ {
		v := p[i] * nstars / nrolls

		fmt.Printf("%2d-%2d: %s (%d)\n", i, i+1,
			strings.Repeat("*", v), v)
	}

	// Output:
	// Uniform Distribution: a=3.0, sigma=7.0
	//  0- 1:  (0)
	//  1- 2:  (0)
	//  2- 3:  (0)
	//  3- 4: ************************ (24)
	//  4- 5: ************************* (25)
	//  5- 6: ************************* (25)
	//  6- 7: ************************ (24)
	//  7- 8:  (0)
	//  8- 9:  (0)
	//  9-10:  (0)

}
