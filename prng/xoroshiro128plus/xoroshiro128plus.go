package xoroshiro128plus

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/shivakar/random/prng"
	"github.com/shivakar/random/prng/splitmix64"
)

var (
	xoroshiro128plus *Xoroshiro128Plus
	_                prng.Engine = xoroshiro128plus
)

/*
// rotl rotates the given input x by k bits to the left
func rotl(x uint64, k uint) uint64 {
	return (x << k) | (x >> (64 - k))
}
*/
// Using const shift amounts as per https://golang.org/src/runtime/hash64.go
func rotl55(x uint64) uint64 {
	return (x << 55) | (x >> (64 - 55))
}
func rotl36(x uint64) uint64 {
	return (x << 36) | (x >> (64 - 36))
}

// Xoroshiro128Plus implements a xoroshiro PRNG with 128 bits of state and
// a maximal period of 2^128-1. The algorithm uses addition as the non-linear
// transformation function
type Xoroshiro128Plus struct {
	seed  uint64
	state [2]uint64
}

// New returns a new instance of the xoroshiro128Plus PRNG Engine.
// If the seed provided is 0, the engine is initialized with current time
func New(seed uint64) *Xoroshiro128Plus {
	r := new(Xoroshiro128Plus)
	r.Seed(seed)
	return r
}

/*
 * Implement 'Engine' interface
 */

// Uint64 returns a pseudo-random 64-bit value in [0, 2^64) as a uint64.
// Uint64 advances the internal state of the engine.
func (x *Xoroshiro128Plus) Uint64() uint64 {
	s0 := x.state[0]
	s1 := x.state[1]
	result := s0 + s1

	s1 ^= s0
	x.state[0] = rotl55(s0) ^ s1 ^ (s1 << 14)
	x.state[1] = rotl36(s1)

	return result
}

// Float64 returns a pseudo-random number in [0.0, 1.0) as a float64.
// Float64 advances the internal state of the engine.
func (x *Xoroshiro128Plus) Float64() float64 {
	return float64(x.Uint64()>>11) / float64(1<<53)
}

// Float64OO returns a pseudo-random number in (0.0, 1.0) as a float64.
// Float6400 advances the internal state of the engine.
func (x *Xoroshiro128Plus) Float64OO() float64 {
	return (float64(x.Uint64()>>12) + float64(0.5)) / float64(1<<52)
}

// Seed uses the provided value to initialize the engine
func (x *Xoroshiro128Plus) Seed(seed uint64) {
	if seed == 0 {
		seed = uint64(time.Now().UnixNano())
	}
	x.seed = seed
	ms := splitmix64.New(seed)
	x.state[0] = ms.Uint64()
	x.state[1] = ms.Uint64()
}

// GetSeed returns the seed used to initialize the engine
func (x *Xoroshiro128Plus) GetSeed() uint64 { return x.seed }

// GetState returns the internal state of the engine as []byte
// GetState can be used to save the state, e.g. to a file
func (x *Xoroshiro128Plus) GetState() []byte {
	const msg = "xoroshiro128plus: Error encoding state"
	buf := new(bytes.Buffer)
	data := []interface{}{
		[]byte("xoroshiro128plus"),
		uint64(x.seed),
		uint64(x.state[0]),
		uint64(x.state[1]),
	}
	for _, v := range data {
		if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
			panic(strings.Join([]string{msg, err.Error()}, "\n"))
		}
	}
	return buf.Bytes()
}

// SetState sets the internal state of the engine from a []byte
// SetState can be used to resume from a saved state
func (x *Xoroshiro128Plus) SetState(b []byte) {
	const msg = "xoroshiro128plus: Error decoding state"
	const algo = "xoroshiro128plus"
	buf := bytes.NewReader(b)
	nb := make([]byte, len(algo))
	_, err := buf.Read(nb)
	if err != nil {
		panic(strings.Join([]string{msg, err.Error()}, "\n"))
	}
	if string(nb) != algo {
		err = fmt.Errorf("Expected '%s', got '%s'", algo, string(nb))
		panic(strings.Join([]string{msg, err.Error()}, "\n"))
	}
	if err = binary.Read(buf, binary.LittleEndian, &x.seed); err != nil {
		panic(strings.Join([]string{msg, err.Error()}, "\n"))
	}
	if err = binary.Read(buf, binary.LittleEndian, &x.state[0]); err != nil {
		panic(strings.Join([]string{msg, err.Error()}, "\n"))
	}
	if err = binary.Read(buf, binary.LittleEndian, &x.state[1]); err != nil {
		panic(strings.Join([]string{msg, err.Error()}, "\n"))
	}
}

// Reset reverts the internal state of the engine to its default state,
// except the seed
func (x *Xoroshiro128Plus) Reset() {
	x.Seed(x.seed)
}
