// Package prngtest defines utilities for testing PRNG Engines
package prngtest

import (
	"bufio"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/shivakar/random/prng"
	"github.com/stretchr/testify/assert"
)

// fileinfo struct contains metadata about a PRNG draws file
type fileinfo struct {
	seed       uint64
	rng        string
	function   string
	start      uint64
	numSamples uint64
	filename   string
}

// parse parses metadata for a PRNG draws file from its filename
// filename is expected to be in the following format
// engine-seed-functionToCall-startIndex-numSamples.txt
func (f *fileinfo) parse(filename string) {
	var err error
	bn := filepath.Base(filename)
	fields := strings.Split(bn, "-")
	nSamples := strings.Split(fields[len(fields)-1], ".")[0]

	f.rng = fields[0]
	if f.seed, err = strconv.ParseUint(fields[1], 10, 64); err != nil {
		panic(err)
	}
	f.function = fields[2]
	if f.start, err = strconv.ParseUint(fields[3], 10, 64); err != nil {
		panic(err)
	}
	if f.numSamples, err = strconv.ParseUint(nSamples, 10, 64); err != nil {
		panic(err)
	}
	f.filename = filename
}

// GetDataFiles return the datafiles matching given pattern in dataDir
func GetDataFiles(dataDir string, pattern string) []string {
	filenames, err := filepath.Glob(filepath.Join(dataDir, pattern))
	if err != nil {
		panic(err)
	}
	return filenames
}

// CompareDraws compares output of an Engine against expected output
func CompareDraws(t *testing.T, e prng.Engine, datafiles []string, longTest bool) {
	assert := assert.New(t)
	assert.NotZero(len(datafiles))
	for _, filename := range datafiles {
		finfo := new(fileinfo)
		finfo.parse(filename)
		e.Reset()
		e.Seed(finfo.seed)
		if !longTest && finfo.start >= 1e9 {
			continue
		}
		for i := uint64(0); i < finfo.start; i++ {
			switch finfo.function {
			case "uint64":
				_ = e.Uint64()
			case "float64":
				_ = e.Float64()
			case "float64oo":
				_ = e.Float64OO()
			}
		}

		var file *os.File
		var err error
		if file, err = os.Open(filename); err != nil {
			panic(err)
		}
		defer file.Close()
		s := bufio.NewScanner(file)
		for s.Scan() {
			switch finfo.function {
			case "uint64":
				v, _ := strconv.ParseUint(s.Text(), 10, 64)
				assert.Equal(v, e.Uint64())
			case "float64":
				v, _ := strconv.ParseFloat(s.Text(), 64)
				assert.InDelta(v, e.Float64(), float64(1e-15))
			case "float64oo":
				v, _ := strconv.ParseFloat(s.Text(), 64)
				assert.InDelta(v, e.Float64OO(), float64(1e-15))
			}
		}
	}
}

// ParseCommandLine parses command line arguments and returns relevant values
func ParseCommandLine() (longTest bool) {
	flag.BoolVar(&longTest, "long", false, "Include long running tests")
	flag.Parse()

	if longTest {
		log.Printf("Note: Running 'long' test.\n\tIncluding tests that require more than 1e9 PRNG draws.\n")

	} else {
		log.Printf("Note: Running tests that require less than 1e9 PRNG draws.\n\tTo run all tests run with -long flag.\n")

	}
	return longTest
}
