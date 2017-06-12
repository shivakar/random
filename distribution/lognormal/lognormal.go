package lognormal

import (
	"math"

	"github.com/shivakar/random/distribution"
	"github.com/shivakar/random/distribution/internal/mathutils"
	"github.com/shivakar/random/distribution/normal"
	"github.com/shivakar/random/prng"
)

var (
	logNormal *LogNormal
	_         distribution.Distribution = logNormal
)

// LogNormal is a continuous random variable for a lognormal probability
// distribution parameterized by mean and standard deviation
type LogNormal struct {
	m    float64
	s    float64
	norm *normal.Normal
}

// New returns a new LogNormal Distribution
func New(r prng.Engine, m float64, s float64) *LogNormal {
	n := new(LogNormal)
	n.Init(r, m, s)
	return n
}

// Init initializes the normal random variable with the supplied
// PRNG Engine and distribution parameters
//
// Two float64 parameters are required for LogNormal distribution,
// namely, mu and sigma
//
// where:
//    -inf < mu < inf and
//    sigma > 0
func (l *LogNormal) Init(r prng.Engine, params ...float64) {
	if len(params) != 2 {
		panic("Error initializing LogNormal Distribution." +
			"Expecting two float64 parameters.")
	}
	if params[1] <= 0 {
		panic("Parameter 'sigma' should be greater than 0")
	}
	l.m = params[0]
	l.s = params[1]
	l.norm = normal.New(r, l.m, l.s)
}

// PDF or probability distribution function returns the relative likelihood
// for the random variable to take on the given value
//
// For LogNormal distribution:
//        = 0 for x <= 0
// PDF(x) = 1/(x*sigma*sqrt(2*pi)) exp(-((ln x-mu)^2)/(2*sigma^2)) for x > 0
func (l *LogNormal) PDF(x float64) float64 {
	if x <= 0.0 {
		return 0.0
	}
	p := (math.Log(x) - l.m) / l.s
	return mathutils.SqrtI2Pi * math.Exp(-(p * p / 2.0)) / (x * l.s)
}

// CDF or cumulative distribution function returns the probability that
// a real-valued random variable X of the probability distribution will be
// found to have a value less than or equal to x
//
// For LogNormal distribution:
// CDF(x) = 0.5 * [ 1 + erf((ln x-mu)/sqrt(2)*sigma) ]
//
// Implemented using an internal Normal Distribution variable and passing
// Log(x) to Normal.CDF
func (l *LogNormal) CDF(x float64) float64 {
	if x <= 0.0 {
		return 0.0
	}
	return l.norm.CDF(math.Log(x))
}

// GetParams returns the current parameters of the Distribution
func (l *LogNormal) GetParams() []float64 { return []float64{l.m, l.s} }

// Float64 returns the next random number satisfying the underlying
// probability distribution
//
// For LogNormal distribution, Float64 return a float64 in [-Inf, Inf]
// Implementation here uses polynomial approximation of the inverse of the
// LogNormal distribution CDF
//
// n = exp(N()*sigma + mu)
//
// Implemented using an internal Normal Distribution variable and
// returning exponentiation of the result.
func (l *LogNormal) Float64() float64 {
	return math.Exp(l.norm.Float64())
}
