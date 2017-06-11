package distribution

import "github.com/shivakar/random/prng"

// Distribution is an interface for a random variable of a probability
// distribution
type Distribution interface {
	// Init initializes a given Distribution with the supplied PRNG Engine and
	// distribution parameters
	Init(prng.Engine, ...float64)

	// Float64 returns the next random number satisfying the underlying
	// probability distribution
	Float64() float64

	// PDF or probability distribution function returns the relative likelihood
	// for the random variable to take on the given value
	PDF(float64) float64

	// CDF or cumulative distribution function returns the probability that
	// a real-valued random variable X of the probability distribution will be
	// found to have a value less than or equal to x
	CDF(float64) float64

	// GetParams returns the current parameters of the Distribution
	GetParams() []float64
}
