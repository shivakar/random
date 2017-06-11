package uniform

import (
	"github.com/shivakar/random/distribution"
	"github.com/shivakar/random/prng"
)

var (
	uniform *Uniform
	_       distribution.Distribution = uniform
)

// Uniform is a continuous random variable for a uniform probability
// distribution parameterized by minimum and maximum values, namely,
// a and b, respectively.
type Uniform struct {
	// rng is source/engine for uniform uint64 and float64 numbers
	rng prng.Engine
	a   float64
	b   float64
}

// New returns a new Uniform Distribution
func New(r prng.Engine, a float64, b float64) *Uniform {
	u := new(Uniform)
	u.Init(r, a, b)
	return u
}

// Init initializes the uniform random variable with the supplied
// PRNG Engine and distribution parameters
//
// Two float64 parameters are required for Uniform distribution,
// namely, a and b
//
// where:
//    -inf < a < b < inf
func (u *Uniform) Init(r prng.Engine, params ...float64) {
	if len(params) != 2 {
		panic("Error initializing Uniform Distribution." +
			"Expecting two float64 parameters.")
	}
	if !(params[0] < params[1]) {
		panic("Parameter 'a' is not less than parameter 'b'")
	}
	u.rng = r
	u.a = params[0]
	u.b = params[1]
}

// PDF or probability distribution function returns the relative likelihood
// for the random variable to take on the given value
//
// For uniform distribution:
//
// PDF(x) = 1/(b-a) for a <= x <= b
//        = 0       for x < a or x > b
//
func (u *Uniform) PDF(x float64) float64 {
	if x < u.a || x > u.b {
		return 0
	}
	return 1.0 / (u.b - u.a)
}

// CDF or cumulative distribution function returns the probability that
// a real-valued random variable X of the probability distribution will be
// found to have a value less than or equal to x
//
// For uniform distribution:
//        = 0           for x < a
// CDF(x) = (x-a)/(b-a) for a <= x <= b
//        = 1           for x > b
func (u *Uniform) CDF(x float64) float64 {
	if x < u.a {
		return 0.0
	}
	if x >= u.b {
		return 1.0
	}
	return (x - u.a) / (u.b - u.a)
}

// GetParams returns the current parameters of the Distribution
func (u *Uniform) GetParams() []float64 { return []float64{u.a, u.b} }

// Float64 returns the next random number satisfying the underlying
// probability distribution
//
// For uniform distribution, since the underlying PRNG Engine already
// produces continuous uniform numbers between [0, 1), a simple linear
// transformation is sufficient to generate random variates
//
// n = a + (b-a)*U
func (u *Uniform) Float64() float64 {
	return u.a + (u.b-u.a)*u.rng.Float64()
}
