package lognormal_test

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/shivakar/random/distribution/internal/distribtest"
	"github.com/shivakar/random/distribution/lognormal"
	"github.com/shivakar/random/distribution/normal"
	"github.com/shivakar/random/prng"
	"github.com/shivakar/random/prng/mt19937"
	"github.com/shivakar/random/prng/splitmix64"
	"github.com/shivakar/random/prng/xoroshiro128plus"
	"github.com/shivakar/random/prng/xorshift1024star"
	"github.com/shivakar/random/prng/xorshift128plus"
	"github.com/stretchr/testify/assert"
)

func Test_LogNormal_GetParams(t *testing.T) {
	assert := assert.New(t)
	data := [...]struct {
		l float64
		s float64
	}{
		{0.0, 1.0},
		{3.0, 5.0},
		{-1.0, 4.0},
	}

	r := xorshift128plus.New(0)
	for _, rec := range data {
		d := lognormal.New(r, rec.l, rec.s)
		params := d.GetParams()
		assert.Equal(rec.l, params[0])
		assert.Equal(rec.s, params[1])
	}
}

func Test_LogNormal_PDF(t *testing.T) {
	assert := assert.New(t)
	data := [...]struct {
		l float64
		s float64
		x []float64
		c []float64
	}{
		{0.000000, 1.000000,
			[]float64{24.8439591804217, 17.4328395707792, 18.5918919830183, 22.2167154867299, 4.1094083811660, 9.9799294228352, 23.0437736463086, 23.3484622639934, 10.2860587764505, 13.1237462103738},
			[]float64{0.0000921586001, 0.0003849515349, 0.0002996706604, 0.0001466585874, 0.0357610319017, 0.0028346421941, 0.0001261599834, 0.0001194762624, 0.0025644333221, 0.0011057964714},
		},
		{1.000000, 1.000000,
			[]float64{17.6614286850543, 25.0966695455328, 13.2624151758366, 5.6353528866048, 13.4458357384545, 18.2469586954401, 24.5292224943635, 8.0258403160566, 7.5091158893571, 8.4620076047561},
			[]float64{0.0039211944617, 0.0013442091164, 0.0085666312938, 0.0542711567386, 0.0082670298489, 0.0035687426543, 0.0014466466449, 0.0276622074409, 0.0317042704497, 0.0247407578963},
		},
		{-1.000000, 1.500000,
			[]float64{18.2243855401159, 1.6528790765271, 1.8997557355857, -2.5009179920674, 9.9009670804427, 8.6259808917944, 3.1985878036220, 11.7054802620758, 12.8510979170728, 18.4525259413339},
			[]float64{0.0004945061989, 0.0974317816377, 0.0769133300931, 0.0000000000000, 0.0024145499090, 0.0033766284908, 0.0294072011182, 0.0015885967773, 0.0012510130058, 0.0004779496690},
		},
		{-5.000000, 5.000000,
			[]float64{-4.0384763636312, 22.0077728438500, 6.9910400199531, 23.6992094023067, 22.3373711724452, 19.5728825454956, -1.5531841432436, -1.3552951211211, 7.5944544213060, 11.2742369006336},
			[]float64{0.0000000000000, 0.0009787974459, 0.0043500701445, 0.0008873182075, 0.0009597219571, 0.0011428143846, 0.0000000000000, 0.0000000000000, 0.0039128592576, 0.0023513332502},
		},
		{10.000000, 50.200000,
			[]float64{2.3303452674682, 3.5686209698457, 13.8168691625990, 18.8281736619250, 13.3808961744407, 37.8650999155177, 6.8792644124165, 7.7181081498806, 32.6923997173702, 39.2302177016140},
			[]float64{0.0033540194704, 0.0021935224189, 0.0005689984643, 0.0004179242456, 0.0005874821889, 0.0002081973641, 0.0011403826013, 0.0010168119248, 0.0002410485516, 0.0002009705234},
		},
		{20.100000, 30.500000,
			[]float64{19.0670163141520, 8.1881563956131, 12.6651672011228, 15.2877478826575, 7.6977993216828, 24.3709206452246, 6.3618057023919, -2.1096951266485, 11.0306129365391, 2.1348642411043},
			[]float64{0.0005856719974, 0.0013421953335, 0.0008750064854, 0.0007274671012, 0.0014259868072, 0.0004602740639, 0.0017190403487, 0.0000000000000, 0.0010020403532, 0.0050108894653},
		},
	}

	for _, rec := range data {
		d := lognormal.New(xorshift128plus.New(0), rec.l, rec.s)
		for i := range rec.x {
			assert.InDelta(rec.c[i], d.PDF(rec.x[i]), float64(1e-12))
		}
	}
}

func Test_LogNormal_CDF(t *testing.T) {
	assert := assert.New(t)
	data := [...]struct {
		l float64
		s float64
		x []float64
		c []float64
	}{
		{0.000000, 1.000000,
			[]float64{16.6366473210180, 23.5730688758199, -2.5466090945665, 22.4608082901078, 1.5595321628359, 14.2689107554608, 11.6923792282603, 17.9876806339965, 17.5361460246574, 11.2316028114131},
			[]float64{0.9975352725417, 0.9992114383028, 0.0000000000000, 0.9990701594661, 0.6716181900975, 0.9960706743232, 0.9930325520315, 0.9980718728325, 0.9979101041226, 0.9922126338335},
		},
		{1.000000, 1.000000,
			[]float64{3.4864885989214, -0.7498225679431, 7.6760585872951, 17.8248210002225, 24.8368720600904, 4.1572119067560, 16.6890692523069, 21.9849047909547, -0.2581256054288, 18.3440396684832},
			[]float64{0.5982790362318, 0.0000000000000, 0.8503896948314, 0.9699862744265, 0.9865280408756, 0.6645250314289, 0.9652191364376, 0.9817070870908, 0.0000000000000, 0.9718886023360},
		},
		{-1.000000, 1.500000,
			[]float64{1.5105008540553, 7.0794230524088, 24.5438756765817, 1.2098229005897, 12.8854693788461, -1.8398049891060, 3.3497580035380, 26.9927264412805, 5.5914274686711, 14.8155766558866},
			[]float64{0.8268083137471, 0.9756644499998, 0.9974473084749, 0.7863001952351, 0.9911235871522, 0.0000000000000, 0.9295694678627, 0.9979064549354, 0.9651731755627, 0.9931260930263},
		},
		{-5.000000, 5.000000,
			[]float64{19.5105898130270, 8.2640299680656, 20.3273666318598, 16.5684729566588, 21.3150020805480, 7.3577594989543, 15.8141081636345, 18.6813337474560, 9.3105843406485, 22.1188966608954},
			[]float64{0.9445534225631, 0.9225423755820, 0.9454657023522, 0.9407971274512, 0.9465062150693, 0.9191161606926, 0.9396904788294, 0.9435741752402, 0.9259437010302, 0.9473071786174},
		},
		{10.000000, 50.200000,
			[]float64{45.9587331977840, 45.8830571956273, 36.8873374585398, 38.9230282712713, 22.6183842646468, 0.9163801215958, 39.9063011183630, 30.5479664853475, 13.3832431836696, 0.7626035884834},
			[]float64{0.4510720356881, 0.4510590378633, 0.4493383025622, 0.4497617823770, 0.4454851899806, 0.4203716856749, 0.4499584789726, 0.4478521828802, 0.4413569251268, 0.4189415733792},
		},
		{20.100000, 30.500000,
			[]float64{7.9644735032236, 8.1471602671073, 15.5942239482673, 13.3103610767171, 17.5672461022108, 27.0648964201508, -2.5358189764421, 21.8919035768612, 13.4481058454035, 1.8254888311942},
			[]float64{0.2772656019426, 0.2775147631614, 0.2846936519238, 0.2829344566823, 0.2860205597768, 0.2908588408472, 0.0000000000000, 0.2884794747216, 0.2830486705353, 0.2613192021241},
		},
	}

	for _, rec := range data {
		d := lognormal.New(xorshift1024star.New(0), rec.l, rec.s)
		for i := range rec.x {
			assert.InDelta(rec.c[i], d.CDF(rec.x[i]), float64(1e-12))
		}
	}
}

func Test_LogNormal_Float64(t *testing.T) {
	assert := assert.New(t)
	engines := [...]struct {
		r    prng.Engine
		name string
	}{
		{mt19937.New(0), "mt19937"},
		{splitmix64.New(0), "splitmix64"},
		{xorshift128plus.New(0), "xorshift128plus"},
		{xorshift1024star.New(0), "xorshift1024star"},
		{xoroshiro128plus.New(0), "xoroshiro128plus"},
	}
	for _, engine := range engines {
		data := []struct {
			l float64
			s float64
		}{
			{0, 1},
			{1, 1},
			{-1, 1.5},
			{-5, 5},
			{10, 50.2},
			{20.10, 30.5},
		}
		for _, rec := range data {
			d := lognormal.New(engine.r, rec.l, rec.s)

			_, pval := distribtest.KSTest(d)
			assert.True(pval > 0.001, fmt.Sprintf("%s, pval:%.6f for (%.6f, %.6f)",
				engine.name, pval, rec.l, rec.s))

			// For normality tests Anderson-Darling test is considered
			// to be better than Kolmogorov-Smirnov test.
			// So doing both
			a2, pval, critvals, sigvals := distribtest.ADTest(d)
			assert.True(pval > 0.001, fmt.Sprintf("%s, pval:%.6f for (%.6f, %.6f)",
				engine.name, pval, rec.l, rec.s))

			// LogNormality is rejected if a2 > critval at any significance level
			for i, v := range critvals {
				assert.True(a2 < v,
					fmt.Sprintf("%s, a2: %.6f, crit:%.6f, sig: %f for (%.6f, %.6f)",
						engine.name, a2, v, sigvals[i], rec.l, rec.s))
			}
		}
	}
}

// Benchmarks
func Benchmark_BuiltInPRNG_LogNormFloat64(b *testing.B) {
	d := rand.New(rand.NewSource(time.Now().UnixNano()))
	x := 0.0
	y := 1.0
	for i := 0; i < b.N; i++ {
		_ = math.Exp(x + (y-x)*d.NormFloat64())
	}
}

func Benchmark_LogNormal_MT19937_Float64(b *testing.B) {
	d := normal.New(mt19937.New(0), 0.0, 1.0)
	for i := 0; i < b.N; i++ {
		_ = d.Float64()
	}
}

func Benchmark_LogNormal_SplitMix64_Float64(b *testing.B) {
	d := normal.New(splitmix64.New(0), 0.0, 1.0)
	for i := 0; i < b.N; i++ {
		_ = d.Float64()
	}
}

func Benchmark_LogNormal_Xoroshiro128Plus_Float64(b *testing.B) {
	d := normal.New(xoroshiro128plus.New(0), 0.0, 1.0)
	for i := 0; i < b.N; i++ {
		_ = d.Float64()
	}
}

func Benchmark_LogNormal_Xorshift128Plus_Float64(b *testing.B) {
	d := normal.New(xorshift128plus.New(0), 0.0, 1.0)
	for i := 0; i < b.N; i++ {
		_ = d.Float64()
	}
}

func Benchmark_LogNormal_Xorshift1024Star_Float64(b *testing.B) {
	d := normal.New(xorshift1024star.New(0), 0.0, 1.0)
	for i := 0; i < b.N; i++ {
		_ = d.Float64()
	}
}

// Example - LogNormal Distribution
func ExampleLogNormal() {
	n := normal.New(xorshift128plus.New(20170611), 3.0, 2.0)
	l := lognormal.New(xorshift128plus.New(20170611), 3.0, 2.0)

	var np [10]int
	var lp [10]int

	nrolls := 10000
	nstars := 100
	for i := 0; i < nrolls; i++ {
		v := n.Float64()
		if v >= 0.0 && v < 10.0 {
			np[int(v)]++
		}

		v = l.Float64()
		if v >= 0.0 && v < 10.0 {
			lp[int(v)]++
		}
	}

	fmt.Println("LogNormal Distribution: mu=3.0, sigma=2.0")
	for i := 0; i < 10; i++ {
		v := lp[i] * nstars / nrolls
		fmt.Printf("%2d-%2d: %s (%d)\n", i, i+1,
			strings.Repeat("*", v), v)
	}
	fmt.Println("Normal Distribution: mu=3.0, sigma=2.0")
	for i := 0; i < 10; i++ {
		v := np[i] * nstars / nrolls

		fmt.Printf("%2d-%2d: %s (%d)\n", i, i+1,
			strings.Repeat("*", v), v)
	}

	// Output:
	// LogNormal Distribution: mu=3.0, sigma=2.0
	//  0- 1: ****** (6)
	//  1- 2: ***** (5)
	//  2- 3: **** (4)
	//  3- 4: **** (4)
	//  4- 5: *** (3)
	//  5- 6: ** (2)
	//  6- 7: ** (2)
	//  7- 8: ** (2)
	//  8- 9: ** (2)
	//  9-10: ** (2)
	// Normal Distribution: mu=3.0, sigma=2.0
	//  0- 1: ******** (8)
	//  1- 2: *************** (15)
	//  2- 3: ******************* (19)
	//  3- 4: ******************* (19)
	//  4- 5: ************** (14)
	//  5- 6: ******** (8)
	//  6- 7: **** (4)
	//  7- 8: * (1)
	//  8- 9:  (0)
	//  9-10:  (0)

}
