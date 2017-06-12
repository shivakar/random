package cauchy_test

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/shivakar/random/distribution/cauchy"
	"github.com/shivakar/random/distribution/internal/distribtest"
	"github.com/shivakar/random/distribution/normal"
	"github.com/shivakar/random/prng"
	"github.com/shivakar/random/prng/mt19937"
	"github.com/shivakar/random/prng/splitmix64"
	"github.com/shivakar/random/prng/xoroshiro128plus"
	"github.com/shivakar/random/prng/xorshift1024star"
	"github.com/shivakar/random/prng/xorshift128plus"
	"github.com/stretchr/testify/assert"
)

func Test_Cauchy_GetParams(t *testing.T) {
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
		d := cauchy.New(r, rec.l, rec.s)
		params := d.GetParams()
		assert.Equal(rec.l, params[0])
		assert.Equal(rec.s, params[1])
	}
}

func Test_Cauchy_PDF(t *testing.T) {
	assert := assert.New(t)
	data := [...]struct {
		l float64
		s float64
		x []float64
		c []float64
	}{
		{-10.000000, 0.500000,
			[]float64{-6.3616384072018, -10.2537240807910, -10.7636935661186, -12.9820467599308, -10.6954220919106, -7.9884210420360, -9.3694558521774, -9.0449657488496, -10.4132398923392, -9.3559148127672},
			[]float64{0.0118000279625, 0.5062568041885, 0.1910101068057, 0.0174080560217, 0.2169470617190, 0.0370433605985, 0.2457665270325, 0.1369557310116, 0.3782493973656, 0.2393862760792},
		},
		{-100.340000, 1.000000,
			[]float64{-100.9799752477762, -97.9004390023405, -101.2326551261166, -100.3638767044087, -100.6596373202898, -124.6801623109993, -102.5842324682186, -99.7823080679022, -100.5670356905037, -97.6996015428361},
			[]float64{0.2258208290948, 0.0457903784407, 0.1771504949694, 0.3181285220635, 0.2888034141911, 0.0005363777436, 0.0527301749218, 0.2427955450714, 0.3027068020872, 0.0399299680070},
		},
		{0.000000, 1.000000,
			[]float64{2.0267398730002, 0.0572587489646, -0.9844330378949, -1.2931781814731, 1.0525638563355, -0.0453383628937, 0.6725491647712, -3.0909978147386, -0.1378628196714, -0.1938864110308},
			[]float64{0.0623199237511, 0.3172696970710, 0.1616517837199, 0.1191141405488, 0.1510087266208, 0.3176569210519, 0.2191730230002, 0.0301593536893, 0.3123728782032, 0.3067775237985},
		},
		{1.000000, 0.500000,
			[]float64{2.7839307237914, 1.6557684754541, 0.5318568869974, 0.6667324959746, 0.6530575986121, 1.9131084299296, -22.3651239134107, 1.0591503527194, 0.2044839799393, 0.3434145403563},
			[]float64{0.0463682944254, 0.2340402722587, 0.3392352934973, 0.4407903298984, 0.4297199016427, 0.1468534679371, 0.0002913970302, 0.6278332108534, 0.1802749180409, 0.2336718537021},
		},
		{10.000000, 50.200000,
			[]float64{-127.9358750048804, -17.6471483036061, -51.6563887281969, -132.8589790105821, -151.5736529283328, 31.1102124983726, 65.7506792219868, 4.2232347675185, 169.4577279964030, -11.1776730570079},
			[]float64{0.0007416179325, 0.0048651604215, 0.0025277274721, 0.0006969054876, 0.0005582026361, 0.0053880224812, 0.0028391347294, 0.0062579648098, 0.0005717703735, 0.0053828445854},
		},
	}

	for _, rec := range data {
		d := cauchy.New(xorshift128plus.New(0), rec.l, rec.s)
		for i := range rec.x {
			assert.InDelta(rec.c[i], d.PDF(rec.x[i]), float64(1e-12))
		}
	}
}

func Test_Cauchy_CDF(t *testing.T) {
	assert := assert.New(t)
	data := [...]struct {
		l float64
		s float64
		x []float64
		c []float64
	}{
		{-10.000000, 0.500000,
			[]float64{-10.0632599656590, -9.9824105398248, -7.8612148440712, -10.0021487984965, -8.5777214866722, -14.7454411076474, -10.1525300232333, -13.4988693780479, -11.1964796130952, -9.5669769298785},
			[]float64{0.4599402993444, 0.5111931822574, 0.9268991183849, 0.4986320408120, 0.8923941527612, 0.0334152024700, 0.4057508142383, 0.0451816353853, 0.1259982778117, 0.7271892975725},
		},
		{-100.340000, 1.000000,
			[]float64{-100.1208969031428, -104.5398228111473, -100.2610770257642, -101.6397096975040, -101.0488420696366, -7.1331033068828, -101.0173328797836, -99.5104260401800, -101.5755075784183, -97.0925282824496},
			[]float64{0.5686577385555, 0.0744057912222, 0.5250699968270, 0.2087487567104, 0.3037187854848, 0.9965850416114, 0.3104940186477, 0.7204345374981, 0.2165897369407, 0.9049151872027},
		},
		{0.000000, 1.000000,
			[]float64{-0.5855664561156, 0.2189596690236, -13.2664911588700, -0.0608252047717, 0.3018055061254, 0.9146667529590, -21.2511946020337, -0.3569124972991, 0.1298747874919, -0.2121687646708},
			[]float64{0.3313788413622, 0.5686141740541, 0.0239482368716, 0.4806625601208, 0.5933005740102, 0.7358228837685, 0.0149674058784, 0.3908771243509, 0.5411103171972, 0.4334514463038},
		},
		{1.000000, 0.500000,
			[]float64{0.7414741272585, 0.1178679533358, 1.2359504062401, 0.0450106837164, 0.6508469449881, -1.0523796523159, 0.5951905687500, 1.9761867782631, 6.9485469032692, 2.0709596065432},
			[]float64{0.3481038710331, 0.1641385330714, 0.6403481261658, 0.1535279359379, 0.3059620428922, 0.0760648156024, 0.2833650891580, 0.8493256029467, 0.9733074797189, 0.8609638485155},
		},
		{10.000000, 50.200000,
			[]float64{17.9259743372359, 217.7643965020225, 33.1502309154175, -372.4783414423412, 14.9496809593914, -29.3126431541364, 9.1829785785344, 130.1932570646486, -34.7686857104322, -47.1777589505588},
			[]float64{0.5498458132084, 0.9245363544759, 0.6375403371714, 0.0415404929543, 0.5312839897544, 0.2885265575372, 0.4948198598269, 0.8740645742008, 0.2681845187667, 0.2293442222010},
		},
	}

	for _, rec := range data {
		d := cauchy.New(xorshift1024star.New(0), rec.l, rec.s)
		for i := range rec.x {
			assert.InDelta(rec.c[i], d.CDF(rec.x[i]), float64(1e-12))
		}
	}
}

func Test_Cauchy_Float64(t *testing.T) {
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

			{-10, 1},
			{-3.5, 3.5},
			{0, 1},
			{1.5, 2.5},
			{1, 32.3},
		}
		for _, rec := range data {
			d := cauchy.New(engine.r, rec.l, rec.s)

			_, pval := distribtest.KSTest(d)
			assert.True(pval > 0.001, fmt.Sprintf("%s, pval:%.6f for (%.6f, %.6f)",
				engine.name, pval, rec.l, rec.s))

		}
	}
}

// Benchmarks
func Benchmark_BuiltInPRNG_CauchyFloat64(b *testing.B) {
	d := rand.New(rand.NewSource(time.Now().UnixNano()))
	x := 0.0
	y := 1.0
	for i := 0; i < b.N; i++ {
		_ = x + y*math.Tan(math.Pi*(d.Float64()-0.5))
	}
}

func Benchmark_Cauchy_MT19937_Float64(b *testing.B) {
	d := cauchy.New(mt19937.New(0), 0.0, 1.0)
	for i := 0; i < b.N; i++ {
		_ = d.Float64()
	}
}

func Benchmark_Cauchy_SplitMix64_Float64(b *testing.B) {
	d := cauchy.New(splitmix64.New(0), 0.0, 1.0)
	for i := 0; i < b.N; i++ {
		_ = d.Float64()
	}
}

func Benchmark_Cauchy_Xoroshiro128Plus_Float64(b *testing.B) {
	d := cauchy.New(xoroshiro128plus.New(0), 0.0, 1.0)
	for i := 0; i < b.N; i++ {
		_ = d.Float64()
	}
}

func Benchmark_Cauchy_Xorshift128Plus_Float64(b *testing.B) {
	d := cauchy.New(xorshift128plus.New(0), 0.0, 1.0)
	for i := 0; i < b.N; i++ {
		_ = d.Float64()
	}
}

func Benchmark_Cauchy_Xorshift1024Star_Float64(b *testing.B) {
	d := cauchy.New(xorshift1024star.New(0), 0.0, 1.0)
	for i := 0; i < b.N; i++ {
		_ = d.Float64()
	}
}

// Example - Cauchy Distribution
func ExampleCauchy() {
	n := normal.New(xorshift128plus.New(20170611), 5.0, 0.5)
	l := cauchy.New(xorshift128plus.New(20170611), 5.0, 0.5)

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

	fmt.Println("Cauchy Distribution: mu=5.0, sigma=0.5")
	for i := 0; i < 10; i++ {
		v := lp[i] * nstars / nrolls
		fmt.Printf("%2d-%2d: %s (%d)\n", i, i+1,
			strings.Repeat("*", v), v)
	}
	fmt.Println("Normal Distribution: mu=5.0, sigma=0.5")
	for i := 0; i < 10; i++ {
		v := np[i] * nstars / nrolls

		fmt.Printf("%2d-%2d: %s (%d)\n", i, i+1,
			strings.Repeat("*", v), v)
	}

	// Output:
	// Cauchy Distribution: mu=5.0, sigma=0.5
	//  0- 1:  (0)
	//  1- 2: * (1)
	//  2- 3: ** (2)
	//  3- 4: ****** (6)
	//  4- 5: *********************************** (35)
	//  5- 6: *********************************** (35)
	//  6- 7: ****** (6)
	//  7- 8: ** (2)
	//  8- 9: * (1)
	//  9-10:  (0)
	// Normal Distribution: mu=5.0, sigma=0.5
	//  0- 1:  (0)
	//  1- 2:  (0)
	//  2- 3:  (0)
	//  3- 4: ** (2)
	//  4- 5: *********************************************** (47)
	//  5- 6: *********************************************** (47)
	//  6- 7: ** (2)
	//  7- 8:  (0)
	//  8- 9:  (0)
	//  9-10:  (0)
}
