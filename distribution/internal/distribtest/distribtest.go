// Package distribtest defines utilities for testing Distributions
package distribtest

import (
	"math"
	"sort"

	"github.com/shivakar/random/distribution"
)

// mean returns the mean/average of the data
func mean(data []float64) float64 {
	tot := float64(0)
	for _, v := range data {
		tot += v
	}
	return tot / float64(len(data))
}

// stddev returns the standard deviation of the data
//
// If ddof=0; then population standard deviation is returned
// If ddof>0; then variance is divided by (n-ddof)
func stddev(data []float64, ddof float64) float64 {
	m := mean(data)
	variance := float64(0)
	for _, v := range data {
		variance += (v - m) * (v - m)
	}
	variance /= float64(len(data)) - ddof

	return math.Sqrt(variance)
}

// ksStatistic return the Kolmogorov-Smirnov statistic
//
// See https://en.wikipedia.org/wiki/Kolmogorov%E2%80%93Smirnov_test#Kolmogorov.E2.80.93Smirnov_statistic
func ksStatistic(rvals []float64, cdf func(float64) float64) float64 {
	d := float64(0)
	da := float64(0)
	db := float64(0)

	sort.Float64s(rvals)
	N := float64(len(rvals))

	for i, v := range rvals {
		cdfval := cdf(v)
		ca := (float64(i+1) / N) - cdfval
		cb := cdfval - (float64(i) / N)
		if ca > da {
			da = ca
		}
		if cb > db {
			db = cb
		}
	}
	if da > db {
		d = da
	} else {
		d = db
	}
	return d
}

// Kolmogorov limiting distribution of two-sided test.
// Returns probability that sqrt(n) * max deviation > y,
// or that max deviation > y/sqrt(n).
// The approximation is useful for the tail of the distribution when n is
// large
//
// Originally from http://www.netlib.org/cephes/index.html
// See browseable source at https://github.com/scipy/scipy/blob/master/scipy/special/cephes/kolmogorov.c
func kolmogorov(y float64) float64 {
	if y < 1.1e-16 {
		return float64(1.0)
	}
	x := float64(-2.0) * y * y

	sign := float64(1.0)
	p := float64(0.0)
	r := float64(1.0)
	t := float64(1.0)
	for (t / p) > float64(1.1e-16) {
		t := math.Exp(x * r * r)
		p += sign * t
		if t == float64(0.0) {
			break
		}
		r += float64(1.0)
		sign = -sign
	}
	return (p + p)
}

// KSTest performs the Kolmogorov-Smirnov test for goodness of fit
// Return the KS test statistic and pvalue
func KSTest(dist distribution.Distribution) (float64, float64) {
	nSamples := 10000
	rvs := make([]float64, nSamples)
	for i := 0; i < nSamples; i++ {
		rvs[i] = dist.Float64()
	}
	d := ksStatistic(rvs, dist.CDF)
	p := kolmogorov(d * math.Sqrt(float64(len(rvs))))

	return d, p
}

// adStatistic returns statistics from Anderson-Darling test for goodness-of-fit
//
// Implementation based on https://en.wikipedia.org/wiki/Anderson%E2%80%93Darling_test
func adStatistic(x []float64, cdf func(float64) float64) (float64, float64, []float64, []float64) {
	sort.Float64s(x)
	// xbar := mean(x)
	// s := stddev(x, 1.0)
	n := float64(len(x))

	// y := make([]float64, len(x))
	z := make([]float64, len(x))
	for i, v := range x {
		// y[i] = (v - xbar) / s
		// z[i] = cdf(y[i])
		z[i] = cdf(v)
		// fmt.Println("Yi:", y[i], "Zi:", z[i])
	}

	// Calculate AD statistic a2 using alternate expression which
	// deals with only one observation in each iteration
	a2 := float64(0)
	for j := range z {
		i := j + 1 // Adjusting for zero based index
		lz := math.Log(z[j])
		t1 := float64((2*i)-1) * lz
		t2 := float64((2*(n-float64(i)))+1) * (float64(1.0) - lz)
		a2 += t1 + t2
	}
	a2 = -n - (a2 / n)

	// Calculate adjusted statistic AD*
	// D'Agostino (1986) in Table 4.7 on p. 123 and
	// on pages 372-373 gives the following adjusted statistic
	//
	// Ralph B. D'Agostino (1986). "Tests for the Normal Distribution".
	// In D'Agostino, R.B. and Stephens, M.A. Goodness-of-Fit Techniques.
	// New York: Marcel Dekker. ISBN 0-8247-7487-6.
	//
	// Not Adjusting assuming both mean and variance are known
	// a2 = a2 * (float64(1.0) + (float64(0.75) / n) + (float64(2.25) / (n * n)))

	// Significance levels and critvals from D'Agostino (1986)
	// sigvals := []float64{10, 5, 2.5, 1, 0.5}
	// critvals := []float64{0.631, 0.752, 0.873, 1.035, 1.159}

	// Significance levels and critical values from
	// https://en.wikipedia.org/wiki/Anderson%E2%80%93Darling_test#Test_for_normality
	sigvals := []float64{15, 10, 5, 2.5, 1}
	critvals := []float64{1.610, 1.933, 2.492, 3.070, 3.857}

	// p-value
	//  D'Agostino (1986) in Table 4.9 on p. 127
	//
	// If AD*>=0.6, then p = exp(1.2937 - 5.709(AD*)+ 0.0186(AD*)^2)
	// If 0.34 <= AD* < .6, then p = exp(0.9177 - 4.279(AD*) - 1.38(AD*)^2)
	// If 0.2 <= AD* < 0.34, then p = 1 - exp(-8.318 + 42.796(AD*)- 59.938(AD*)^2)
	// If AD* < 0.2, then p = 1 - exp(-13.436 + 101.14(AD*)- 223.73(AD*)^2)

	pval := float64(0)
	if a2 >= float64(0.6) {
		pval = math.Exp(1.2937 - 5.709*a2 + 0.0816*a2*a2)
	} else if a2 >= float64(0.34) {
		pval = math.Exp(0.9177 - 4.279*a2 - 1.38*a2*a2)
	} else if a2 >= float64(0.2) {
		pval = float64(1.0) - math.Exp(-8.318+42.796*a2-59.938*a2*a2)
	} else {
		pval = float64(1.0) - math.Exp(-13.436+101.14*a2-223.73*a2*a2)
	}

	return a2, pval, critvals, sigvals
}

// ADTest performs the Anderson-Darling test for goodness of fit
// Return the AD statistic, pvalue, critical values and significance levels
func ADTest(dist distribution.Distribution) (float64, float64, []float64, []float64) {
	nSamples := 10000
	rvs := make([]float64, nSamples)
	for i := 0; i < nSamples; i++ {
		rvs[i] = dist.Float64()
	}
	stat, pval, critVals, sigs := adStatistic(rvs, dist.CDF)
	return stat, pval, critVals, sigs
}
