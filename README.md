[![Build Status](https://travis-ci.org/shivakar/random.svg?branch=master)](https://travis-ci.org/shivakar/random) [![Coverage Status](https://coveralls.io/repos/github/shivakar/random/badge.svg?branch=master)](https://coveralls.io/github/shivakar/random?branch=master) [![GoDoc](https://godoc.org/github.com/shivakar/random?status.svg)](https://godoc.org/github.com/shivakar/random)

# random
Package random implements pseudo-random number generators and random variate generators.

## Installation

```
go get github.com/shivakar/random
```

## Features

PRNG Engines available:

* Mersenne Twister: mt19937 64-bit
    * See http://www.math.sci.hiroshima-u.ac.jp/~m-mat/MT/emt64.html for
      details and reference implementation
* SplitMix64: Pseduo RNG based on avalanching function
    * See http://prng.di.unimi.it/splitmix64.c for details
      and reference implementation
* Xorshift128Plus: Fast generator passing BigCrush
    * See http://xorshift.di.unimi.it/xorshift128plus.c for details
      and reference implementation
* Xorshift1024Star: Fast generator with maximal period of 2^1024 - 1
    * See http://xorshift.di.unimi.it/xorshift1024star.c for details
      and reference implementation

## License

Random is licensed under a MIT license.
