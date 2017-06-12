package cauchy

import (
	"math"

	"github.com/shivakar/random/distribution"
	"github.com/shivakar/random/distribution/internal/mathutils"
	"github.com/shivakar/random/prng"
)

var (
	cauchy *Cauchy
	_      distribution.Distribution = cauchy
)

// Cauchy is a continuous random variable for a lognormal probability
// distribution parameterized by location and scale
type Cauchy struct {
	rng prng.Engine
	l   float64
	s   float64
}

// New returns a new Cauchy Distribution
func New(r prng.Engine, l float64, s float64) *Cauchy {
	c := new(Cauchy)
	c.Init(r, l, s)
	return c
}

// Init initializes the normal random variable with the supplied
// PRNG Engine and distribution parameters
//
// Two float64 parameters are required for Cauchy distribution,
// namely, location and scale
//
// where:
//    -inf < l < inf and
//    scale > 0
func (c *Cauchy) Init(r prng.Engine, params ...float64) {
	if len(params) != 2 {
		panic("Error initializing Cauchy Distribution." +
			"Expecting two float64 parameters.")
	}
	if params[1] <= 0 {
		panic("Parameter 'scale' should be greater than 0")
	}
	c.rng = r
	c.l = params[0]
	c.s = params[1]
}

// PDF or probability distribution function returns the relative likelihood
// for the random variable to take on the given value
//
// For Cauchy distribution:
// PDF(x) = 1 / (pi*scale * [ 1 + ((x-location)/scale)^2 ])
func (c *Cauchy) PDF(x float64) float64 {
	d := (x - c.l) / c.s
	return 1.0 / (math.Pi * c.s * (1.0 + d*d))
}

// CDF or cumulative distribution function returns the probability that
// a real-valued random variable X of the probability distribution will be
// found to have a value less than or equal to x
//
// For Cauchy distribution:
// CDF(x) = 0.5 + 1/pi * arctan((x-location)/scale)
func (c *Cauchy) CDF(x float64) float64 {
	return 0.5 + (mathutils.Ipi * math.Atan2(x-c.l, c.s))
}

// GetParams returns the current parameters of the Distribution
func (c *Cauchy) GetParams() []float64 { return []float64{c.l, c.s} }

// Float64 returns the next random number satisfying the underlying
// probability distribution
//
// For Cauchy distribution, Float64 return a float64 in [-Inf, Inf]
//
// n = location + scale * tan(pi * U(0, 1) - 0.5)
//
// Implemented using an internal Normal Distribution variable and
// returning exponentiation of the result.
func (c *Cauchy) Float64() float64 {
	return c.l + c.s*math.Tan(math.Pi*(c.rng.Float64OO()-0.5))
}
