package xorshift1024star

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
	xorshift1024star *Xorshift1024star
	_                prng.Engine = xorshift1024star
)

// Xorshift1024star implements a Xorshift PRNG with 1024 bits of state and
// a maximal period of 2^1024-1. The algorithm uses multiplication as the
// non-linear transformation function
type Xorshift1024star struct {
	seed  uint64
	state [16]uint64
	index int
}

// New returns a new instance of the Xorshift1024star PRNG Engine.
// If the seed provided is 0, the engine is initialized with current time
func New(seed uint64) *Xorshift1024star {
	r := new(Xorshift1024star)
	r.Seed(seed)
	return r
}

/*
 * Implement 'Engine' interface
 */

// Uint64 returns a pseudo-random 64-bit value in [0, 2^64) as a uint64.
// Uint64 advances the internal state of the engine.
func (x *Xorshift1024star) Uint64() uint64 {
	s0 := x.state[x.index]
	x.index = (x.index + 1) & 15
	s1 := x.state[x.index]
	s1 ^= s1 << 31
	x.state[x.index] = s1 ^ s0 ^ (s1 >> 11) ^ (s0 >> 30)
	return x.state[x.index] * uint64(1181783497276652981)
}

// Float64 returns a pseudo-random number in [0.0, 1.0) as a float64.
// Float64 advances the internal state of the engine.
func (x *Xorshift1024star) Float64() float64 {
	return float64(x.Uint64()>>11) / float64(1<<53)
}

// Float64OO returns a pseudo-random number in (0.0, 1.0) as a float64.
// Float6400 advances the internal state of the engine.
func (x *Xorshift1024star) Float64OO() float64 {
	return (float64(x.Uint64()>>12) + float64(0.5)) / float64(1<<52)
}

// Seed uses the provided value to initialize the engine
func (x *Xorshift1024star) Seed(seed uint64) {
	if seed == 0 {
		seed = uint64(time.Now().UnixNano())
	}
	x.seed = seed
	x.index = 0
	ms := splitmix64.New(seed)
	for i := 0; i < len(x.state); i++ {
		x.state[i] = ms.Uint64()
	}
}

// GetSeed returns the seed used to initialize the engine
func (x *Xorshift1024star) GetSeed() uint64 { return x.seed }

// GetState returns the internal state of the engine as []byte
// GetState can be used to save the state, e.g. to a file
func (x *Xorshift1024star) GetState() []byte {
	const msg = "xorshift1024star: Error encoding state"
	buf := new(bytes.Buffer)
	data := []interface{}{
		[]byte("xorshift1024star"),
		uint64(x.seed),
		uint64(x.index),
	}
	for _, v := range data {
		if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
			panic(strings.Join([]string{msg, err.Error()}, "\n"))
		}
	}
	for _, v := range x.state {
		if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
			panic(strings.Join([]string{msg, err.Error()}, "\n"))
		}
	}
	return buf.Bytes()
}

// SetState sets the internal state of the engine from a []byte
// SetState can be used to resume from a saved state
func (x *Xorshift1024star) SetState(b []byte) {
	const msg = "xorshift1024star: Error decoding state"
	const algo = "xorshift1024star"
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
	var index uint64
	if err = binary.Read(buf, binary.LittleEndian, &index); err != nil {
		panic(strings.Join([]string{msg, err.Error()}, "\n"))
	}
	x.index = int(index)
	for i := 0; i < len(x.state); i++ {
		if err = binary.Read(buf, binary.LittleEndian, &x.state[i]); err != nil {
			panic(strings.Join([]string{msg, err.Error()}, "\n"))
		}
	}
}

// Reset reverts the internal state of the engine to its default state,
// except the seed
func (x *Xorshift1024star) Reset() {
	x.Seed(x.seed)
}
