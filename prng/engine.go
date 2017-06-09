package prng

// Engine is an interface for a Pseudo-Random Number Generator of
// uniformly-distributed values
type Engine interface {
	// Uint64 returns a pseudo-random 64-bit value in [0, 2^64) as a uint64.
	// Uint64 advances the internal state of the engine.
	Uint64() uint64

	// Float64 returns a pseudo-random number in [0.0, 1.0) as a float64.
	// Float64 advances the internal state of the engine.
	Float64() float64

	// Float64OO returns a pseudo-random number in (0.0, 1.0) as a float64.
	// Float6400 advances the internal state of the engine.
	Float64OO() float64

	// Seed uses the provided value to initialize the engine
	Seed(uint64)

	// GetSeed returns the seed used to initialize the engine
	GetSeed() uint64

	// GetState returns the internal state of the engine as []byte
	// GetState can be used to save the state, e.g. to a file
	GetState() []byte

	// SetState sets the internal state of the engine from a []byte
	// SetState can be used to resume from a saved state
	SetState([]byte)

	// Reset reverts the internal state of the engine to its default state,
	// except the seed
	Reset()
}
