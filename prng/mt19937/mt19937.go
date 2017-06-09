package mt19937

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/shivakar/random/prng"
)

var (
	mt19937 *MT19937
	_       prng.Engine = mt19937
)

// Constants
const (
	nn      int    = 312
	mm      int    = 156
	matrixA uint64 = 0xB5026F5AA96619E9
	um      uint64 = 0xFFFFFFFF80000000 // Most significant 33 bits
	lm      uint64 = 0x000000007FFFFFFF // Least significant 31 bits
)

// MT19937 implements the 64-bit variant of the Mersenne Twister algorithm
// based on the Mersenne prime 2^19937-1 as its period.
type MT19937 struct {
	seed  uint64
	index int
	state [nn]uint64
}

// New returns a new instance of the MT19937 PRNG Engine.
// If the seed provided is 0, the engine is initialized with current time
func New(seed uint64) *MT19937 {
	r := new(MT19937)
	r.Seed(seed)
	return r
}

/*
 * Implement 'Engine' interface
 */

// Uint64 returns a pseudo-random 64-bit value in [0, 2^64) as a uint64.
// Uint64 advances the internal state of the engine.
func (r *MT19937) Uint64() uint64 {
	if r.index >= nn {
		for i := 0; i < nn-mm; i++ {
			y := (r.state[i] & um) | (r.state[i+1] & lm)
			r.state[i] = r.state[i+mm] ^ (y >> 1) ^ ((y & 1) * matrixA)

		}

		for i := nn - mm; i < nn-1; i++ {
			y := (r.state[i] & um) | (r.state[i+1] & lm)
			r.state[i] = r.state[i+(mm-nn)] ^ (y >> 1) ^ ((y & 1) * matrixA)

		}
		y := (r.state[nn-1] & um) | (r.state[0] & lm)
		r.state[nn-1] = r.state[mm-1] ^ (y >> 1) ^ ((y & 1) * matrixA)

		r.index = 0
	}
	y := r.state[r.index]
	r.index++

	y ^= (y >> 29) & 0x5555555555555555
	y ^= (y << 17) & 0x71D67FFFEDA60000
	y ^= (y << 37) & 0xFFF7EEE000000000
	y ^= (y >> 43)

	return y
}

// Float64 returns a pseudo-random number in [0.0, 1.0) as a float64.
// Float64 advances the internal state of the engine.
func (r *MT19937) Float64() float64 {
	return float64(r.Uint64()>>11) / float64(1<<53)
}

// Float64OO returns a pseudo-random number in (0.0, 1.0) as a float64.
// Float6400 advances the internal state of the engine.
func (r *MT19937) Float64OO() float64 {
	return (float64(r.Uint64()>>12) + float64(0.5)) / float64(1<<52)
}

// Seed uses the provided value to initialize the engine
// If the seed provided is 0, the engine is initialized with current time
func (r *MT19937) Seed(seed uint64) {
	if seed == 0 {
		seed = uint64(time.Now().UnixNano())
	}
	r.seed = seed
	r.state[0] = seed
	for mti := uint64(1); mti < uint64(nn); mti++ {
		r.state[mti] = (uint64(6364136223846793005)*
			(r.state[mti-1]^(r.state[mti-1]>>62)) + mti)
	}
	r.index = nn
}

// GetSeed returns the seed used to initialize the engine
func (r *MT19937) GetSeed() uint64 { return r.seed }

// GetState returns the internal state of the engine as []byte
// GetState can be used to save the state, e.g. to a file
func (r *MT19937) GetState() []byte {
	const msg = "mt19937: Error encoding state"
	buf := new(bytes.Buffer)
	data := []interface{}{
		[]byte("mt19937"),
		uint64(r.seed),
		uint64(r.index),
	}

	for _, v := range data {
		if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
			panic(strings.Join([]string{msg, err.Error()}, "\n"))
		}
	}

	for _, v := range r.state {
		if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
			panic(strings.Join([]string{msg, err.Error()}, "\n"))
		}
	}
	return buf.Bytes()
}

// SetState sets the internal state of the engine from a []byte
// SetState can be used to resume from a saved state
func (r *MT19937) SetState(b []byte) {
	const msg = "mt19937: Error decoding state"
	buf := bytes.NewReader(b)
	nb := make([]byte, 7)
	_, err := buf.Read(nb)
	if err != nil {
		panic(strings.Join([]string{msg, err.Error()}, "\n"))
	}
	if string(nb) != "mt19937" {
		err = fmt.Errorf("Expected 'mt19937', got '%s'", string(nb))
		panic(strings.Join([]string{msg, err.Error()}, "\n"))
	}
	if err = binary.Read(buf, binary.LittleEndian, &r.seed); err != nil {
		panic(strings.Join([]string{msg, err.Error()}, "\n"))
	}
	var index uint64
	if err = binary.Read(buf, binary.LittleEndian, &index); err != nil {
		panic(strings.Join([]string{msg, err.Error()}, "\n"))
	}
	r.index = int(index)

	for i := 0; i < nn; i++ {
		if err = binary.Read(buf, binary.LittleEndian, &r.state[i]); err != nil {
			panic(strings.Join([]string{msg, err.Error()}, "\n"))
		}
	}
}

// Reset reverts the internal state of the engine to its default state,
// except the seed
func (r *MT19937) Reset() {
	r.Seed(r.seed)
}
