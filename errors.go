package bloomfilter

import (
	"errors"
)

var (
	// ErrInvalidNumberOfItems is returned when numItems in a bloom filter is not positive
	ErrInvalidNumberOfItems = errors.New("number of items must be positive")

	// ErrInvalidFalsePositiveRate is returned when false positive rate is not greater than 0.0 
	// or less than 1.0
	ErrInvalidFalsePositiveRate = errors.New("false positive rate must be in range of (0.0, 1.0)")

	// ErrInvalidSize is returned when bloom filter size in bits is 0
	ErrInvalidSize = errors.New("size is the number of bits in bloom filter structure and should be positive")

	// ErrInvalidNumberOfHashFunctions is returned when number of hash functions is not positive
	ErrInvalidNumberOfHashFunctions = errors.New("number of hash functions should be positive")
)