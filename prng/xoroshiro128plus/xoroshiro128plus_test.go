package xoroshiro128plus_test

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"

	"github.com/shivakar/random/prng/internal/prngtest"
	"github.com/shivakar/random/prng/xoroshiro128plus"
	"github.com/stretchr/testify/assert"
)

var datadir = filepath.Join("..", "..", "data", "xoroshiro128plus")
var longTest bool

// TestMain is the entry point for testing
func TestMain(m *testing.M) {
	longTest = prngtest.ParseCommandLine()
	os.Exit(m.Run())
}

func Test_xoroshiro128Plus_GetSetSeed(t *testing.T) {
	assert := assert.New(t)

	seeds := []uint64{1, 5, 10, 1024, 200000}
	r := xoroshiro128plus.New(0)
	for _, seed := range seeds {
		r.Seed(seed)
		assert.Equal(seed, r.GetSeed())
	}
	r.Seed(0)
	assert.NotEqual(0, r.GetSeed())
}

func Test_xoroshiro128Plus_GetSetState(t *testing.T) {
	assert := assert.New(t)

	// Checking seed remains same after getting and setting
	// states
	seeds := []uint64{1, 5, 10, 1024}
	states := make([][]byte, len(seeds))
	for i, seed := range seeds {
		r := xoroshiro128plus.New(seed)
		states[i] = r.GetState()
	}
	for i, state := range states {
		r := xoroshiro128plus.New(0)
		r.SetState(state)
		assert.Equal(seeds[i], r.GetSeed())
	}

	// Checking that the streams remain same after getting and setting states
	r1 := xoroshiro128plus.New(0)
	for i := 0; i < 10; i++ {
		_ = r1.Uint64()
	}
	r2 := xoroshiro128plus.New(0)
	r2.SetState(r1.GetState())
	for i := 0; i < 10; i++ {
		_ = r1.Uint64()
		_ = r2.Uint64()
	}
	assert.Equal(r1.Uint64(), r2.Uint64())
	assert.Equal(r1.Float64(), r2.Float64())
	assert.Equal(r1.Float64OO(), r2.Float64OO())

	// Checking cases where SetState should panic
	assert.Panics(func() {
		r1 := xoroshiro128plus.New(0)
		r1.SetState([]byte("Hello"))
	})
	assert.Panics(func() {
		r1 := xoroshiro128plus.New(0)
		r1.SetState(nil)
	})

	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, []byte("xoroshiro128plust"))
	assert.Panics(func() {
		r1 := xoroshiro128plus.New(0)
		r1.SetState(buf.Bytes())
	})
	buf = new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, []byte("xoroshiro128plus"))
	_ = binary.Write(buf, binary.LittleEndian, uint64(10))
	_ = binary.Write(buf, binary.LittleEndian, []byte("h"))
	assert.Panics(func() {
		r1 := xoroshiro128plus.New(0)
		r1.SetState(buf.Bytes())
	})
	buf = new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, []byte("xoroshiro128plus"))
	_ = binary.Write(buf, binary.LittleEndian, uint64(10))
	_ = binary.Write(buf, binary.LittleEndian, uint64(10))
	_ = binary.Write(buf, binary.LittleEndian, []byte("h"))
	assert.Panics(func() {
		r1 := xoroshiro128plus.New(0)
		r1.SetState(buf.Bytes())
	})
}

func Test_xoroshiro128Plus_Uint64(t *testing.T) {
	e := xoroshiro128plus.New(0)
	filenames := prngtest.GetDataFiles(datadir, "*uint64*.txt")
	prngtest.CompareDraws(t, e, filenames, longTest)
}

func Test_xoroshiro128Plus_Float64(t *testing.T) {
	e := xoroshiro128plus.New(0)
	filenames := prngtest.GetDataFiles(datadir, "*float64*.txt")
	prngtest.CompareDraws(t, e, filenames, longTest)
}

func Test_xoroshiro128Plus_Float64OO(t *testing.T) {
	e := xoroshiro128plus.New(0)
	filenames := prngtest.GetDataFiles(datadir, "*float64oo*.txt")
	prngtest.CompareDraws(t, e, filenames, longTest)
}

// Benchmarks
func Benchmark_xoroshiro128Plus_Uint64(b *testing.B) {
	rng := xoroshiro128plus.New(0)
	for i := 0; i < b.N; i++ {
		_ = rng.Uint64()
	}
}

func Benchmark_xoroshiro128Plus_Float64(b *testing.B) {
	rng := xoroshiro128plus.New(0)
	for i := 0; i < b.N; i++ {
		_ = rng.Float64()
	}
}

func Benchmark_xoroshiro128Plus_Float64OO(b *testing.B) {
	rng := xoroshiro128plus.New(0)
	for i := 0; i < b.N; i++ {
		_ = rng.Float64()
	}
}
