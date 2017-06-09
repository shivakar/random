package prng_test

import (
	"math/rand"
	"testing"
	"time"
)

// Benchmarks for built-in PRNG
func Benchmark_BuiltInPRNG_Int63(b *testing.B) {
	rng := rand.New(rand.NewSource(int64(time.Now().UnixNano())))
	for i := 0; i < b.N; i++ {
		_ = rng.Int63()
	}
}

/*
func Benchmark_BuiltInPRNG_Uint64(b *testing.B) {
	rng := rand.New(rand.NewSource(int64(time.Now().UnixNano())))
	for i := 0; i < b.N; i++ {
		_ = rng.Uint64()
	}
}*/

func Benchmark_BuiltInPRNG_Float64(b *testing.B) {
	rng := rand.New(rand.NewSource(int64(time.Now().UnixNano())))
	for i := 0; i < b.N; i++ {
		_ = rng.Float64()
	}
}
