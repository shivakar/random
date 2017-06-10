package splitmix64

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/shivakar/random/prng"
)

var (
	splitmix64 *SplitMix64
	_          prng.Engine = splitmix64
)

// SplitMix64 implements the avalanching function based PRNG
// by Sebastiano Vigna
type SplitMix64 struct {
	seed  uint64
	state uint64
}

// New returns a new instance of the SplitMix64 PRNG Engine.
// If the seed provided is 0, the engine is initialized with current time
func New(seed uint64) *SplitMix64 {
	r := new(SplitMix64)
	r.Seed(seed)
	return r
}

/*
 * Implement 'Engine' interface
 */

// Uint64 returns a pseudo-random 64-bit value in [0, 2^64) as a uint64.
// Uint64 advances the internal state of the engine.
func (s *SplitMix64) Uint64() uint64 {
	s.state += uint64(0x9E3779B97F4A7C15)
	z := s.state
	z = (z ^ (z >> 30)) * uint64(0xBF58476D1CE4E5B9)
	z = (z ^ (z >> 27)) * uint64(0x94D049BB133111EB)
	return z ^ (z >> 31)
}

// Float64 returns a pseudo-random number in [0.0, 1.0) as a float64.
// Float64 advances the internal state of the engine.
func (s *SplitMix64) Float64() float64 {
	return float64(s.Uint64()>>11) / float64(1<<53)
}

// Float64OO returns a pseudo-random number in (0.0, 1.0) as a float64.
// Float6400 advances the internal state of the engine.
func (s *SplitMix64) Float64OO() float64 {
	return (float64(s.Uint64()>>12) + float64(0.5)) / float64(1<<52)
}

// Seed uses the provided value to initialize the engine
func (s *SplitMix64) Seed(seed uint64) {
	if seed == uint64(0) {
		seed = uint64(time.Now().UnixNano())
	}
	s.seed = seed
	s.state = seed
}

// GetSeed returns the seed used to initialize the engine
func (s *SplitMix64) GetSeed() uint64 { return s.seed }

// GetState returns the internal state of the engine as []byte
// GetState can be used to save the state, e.g. to a file
func (s *SplitMix64) GetState() []byte {
	const msg = "splitmix64: Error encoding state"
	buf := new(bytes.Buffer)
	data := []interface{}{
		[]byte("splitmix64"),
		uint64(s.seed),
		uint64(s.state),
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
func (s *SplitMix64) SetState(b []byte) {
	const msg = "splitmix64: Error decoding state"
	buf := bytes.NewReader(b)
	nb := make([]byte, 10)
	_, err := buf.Read(nb)
	if err != nil {
		panic(strings.Join([]string{msg, err.Error()}, "\n"))
	}
	if string(nb) != "splitmix64" {
		err = fmt.Errorf("Expected 'splitmix64', got '%s'", string(nb))
		panic(strings.Join([]string{msg, err.Error()}, "\n"))
	}
	if err = binary.Read(buf, binary.LittleEndian, &s.seed); err != nil {
		panic(strings.Join([]string{msg, err.Error()}, "\n"))
	}
	if err = binary.Read(buf, binary.LittleEndian, &s.state); err != nil {
		panic(strings.Join([]string{msg, err.Error()}, "\n"))
	}
}

// Reset reverts the internal state of the engine to its default state,
// except the seed
func (s *SplitMix64) Reset() {
	s.Seed(s.seed)
}
