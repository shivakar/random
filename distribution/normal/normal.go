package normal

import (
	"math"

	"github.com/shivakar/random/distribution"
	"github.com/shivakar/random/distribution/internal/mathutils"
	"github.com/shivakar/random/prng"
)

var (
	normal *Normal
	_      distribution.Distribution = normal
)

// Normal is a continuous random variable for a normal probability
// distribution parameterized by mean and standard deviation.
type Normal struct {
	// rng is source/engine for normal uint64 and float64 numbers
	rng prng.Engine
	m   float64
	s   float64
}

// New returns a new Normal Distribution
func New(r prng.Engine, m float64, s float64) *Normal {
	n := new(Normal)
	n.Init(r, m, s)
	return n
}

// Init initializes the normal random variable with the supplied
// PRNG Engine and distribution parameters
//
// Two float64 parameters are required for Normal distribution,
// namely, mu and sigma
//
// where:
//    -inf < mu < inf and
//    sigma > 0
func (n *Normal) Init(r prng.Engine, params ...float64) {
	if len(params) != 2 {
		panic("Error initializing Normal Distribution." +
			"Expecting two float64 parameters.")
	}
	if params[1] <= 0 {
		panic("Parameter 'sigma' should be greater than 0")
	}
	n.rng = r
	n.m = params[0]
	n.s = params[1]
}

// PDF or probability distribution function returns the relative likelihood
// for the random variable to take on the given value
//
// For normal distribution:
//
// PDF(x) = 1/(sqrt(2*pi*sigma^2)) * e^(-(x-mu)^2/2*sigma^2)
//
func (n *Normal) PDF(x float64) float64 {
	x = (x - n.m) / n.s
	return math.Exp(-(x*x)/2.0) / (mathutils.Sqrt2Pi * n.s)
}

// CDF or cumulative distribution function returns the probability that
// a real-valued random variable X of the probability distribution will be
// found to have a value less than or equal to x
//
// For normal distribution:
// CDF(x) = 0.5 * [ 1 + erf((x-mu)/sqrt(2)*sigma) ]
func (n *Normal) CDF(x float64) float64 {
	x = (x - n.m) / n.s
	return 0.5 + 0.5*math.Erf(x*mathutils.SqrtI2)
}

// GetParams returns the current parameters of the Distribution
func (n *Normal) GetParams() []float64 { return []float64{n.m, n.s} }

// Float64 returns the next random number satisfying the underlying
// probability distribution
//
// For normal distribution, Float64 return a float64 in [-Inf, Inf]
// Implementation here uses polynomial approximation of the inverse of the
// Normal distribution CDF
//
// n = N()*sigma + mu
func (n *Normal) Float64() float64 {
	return mathutils.NormICDF(n.rng.Float64OO())*n.s + n.m
}
