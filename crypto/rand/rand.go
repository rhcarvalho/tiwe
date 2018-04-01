// Package rand provides a cryptographically secure source of random numbers.
// It does so by exporting a math/rand.Rand value backed by crypto/rand and
// convenience functions that replicate those found in math/rand.
package rand

import (
	"math/rand"
)

// Rand is a source of random numbers backed by crypto/rand. It is safe for
// concurrent use by multiple goroutines.
var Rand = rand.New(source)

// Top-level convenience functions, as seen in math/Rand.

// Int63 returns a non-negative random 63-bit integer as an int64.
func Int63() int64 { return Rand.Int63() }

// Uint32 returns a random 32-bit value as a uint32.
func Uint32() uint32 { return Rand.Uint32() }

// Uint64 returns a random 64-bit value as a uint64.
func Uint64() uint64 { return Rand.Uint64() }

// Int31 returns a non-negative random 31-bit integer as an int32.
func Int31() int32 { return Rand.Int31() }

// Int returns a non-negative random int.
func Int() int { return Rand.Int() }

// Int63n returns, as an int64, a non-negative random number in [0,n).
// It panics if n <= 0.
func Int63n(n int64) int64 { return Rand.Int63n(n) }

// Int31n returns, as an int32, a non-negative random number in [0,n).
// It panics if n <= 0.
func Int31n(n int32) int32 { return Rand.Int31n(n) }

// Intn returns, as an int, a non-negative random number in [0,n).
// It panics if n <= 0.
func Intn(n int) int { return Rand.Intn(n) }

// Float64 returns, as a float64, a random number in [0.0,1.0).
func Float64() float64 { return Rand.Float64() }

// Float32 returns, as a float32, a random number in [0.0,1.0).
func Float32() float32 { return Rand.Float32() }

// Perm returns, as a slice of n ints, a random permutation of the integers
// [0,n).
func Perm(n int) []int { return Rand.Perm(n) }

// Shuffle randomizes the order of elements.
// n is the number of elements. Shuffle panics if n < 0.
// swap swaps the elements with indexes i and j.
func Shuffle(n int, swap func(i, j int)) { Rand.Shuffle(n, swap) }

// Read generates len(p) random bytes and writes them into p. It always returns
// len(p) and a nil error.
func Read(p []byte) (n int, err error) { return Rand.Read(p) }

// NormFloat64 returns a normally distributed float64 in the range
// [-math.MaxFloat64, +math.MaxFloat64] with standard normal distribution
// (mean = 0, stddev = 1).
// To produce a different normal distribution, callers can adjust the output
// using:
//
//  sample = NormFloat64() * desiredStdDev + desiredMean
//
func NormFloat64() float64 { return Rand.NormFloat64() }

// ExpFloat64 returns an exponentially distributed float64 in the range
// (0, +math.MaxFloat64] with an exponential distribution whose rate parameter
// (lambda) is 1 and whose mean is 1/lambda (1).
// To produce a distribution with a different rate parameter, callers can adjust
// the output using:
//
//  sample = ExpFloat64() / desiredRateParameter
//
func ExpFloat64() float64 { return Rand.ExpFloat64() }
