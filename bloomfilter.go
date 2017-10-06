// Copyright (c) 2017 Mehmet Rauf Celik
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Package bloomfilter provides necessary methods and data structures for creating
// Bloom Filters.
//
// A Bloom Filter is a space-efficient probabilistic data structure, conceived by 
// Burton Howard Bloom in 1970, that is used to test whether an element is a member of a set.
// It is possible to add an elemen to a bloom filter or query an element for existence, 
// while removing an element is not possible.
// When testing (querying) an element in the bloom filter, either true or false is returned.
// However it should be noted that false positives are possible, while false negatives are not.
// For additional information see: 
// Bloom filter - Wikipedia https://en.wikipedia.org/wiki/Bloom_filter
//
// This implementation of bloomfilter provides two different ways of creating two different
// bloom filter data structures.
//
// It is possible to create a new bloom filter data structure by providing either estimated number of items
// that the bloom filter will hold and estimated false positove error rate along with two custom hash functions:
//
//     bf := NewByEstimates(numItems uint64, fpRate float64, hash1 HashFunction64, hash2 HashFunction64)
// 
// or maximum bloom filter size (in bits) and number or hash functions that will be created by double hashing
// of two hash functions:
// 
//     bf := NewBySizeAndNumHashFuncs(size uint64, numHashFunctions uint8, hash1 HashFunction64, hash2 HashFunction64)
//
// In both cases, when hash1 and/or hash2 values are nil, default hash function implementations of FNV-1 and/or FNV-1a
// from the standard library will be used.
//
// Both of the above bloom filter data structures are non-thread safe, however it is possible to create a thread safe
// implementation by:
//
//     bfts := NewTSByEstimates(numItems uint64, fpRate float64, hash1 HashFunction64, hash2 HashFunction64)
//
// or
// 
//     bfts := NewTSBySizeAndNumHashFuncs(size uint64, numHashFunctions uint8, hash1 HashFunction64, hash2 HashFunction64)
//
// Once a bloom filter structure is created, one can add an element by;
//
//     bf.Add([]byte("data"))
//
// and test an element for existence by;
//
//     exists := bf.Query([]byte("data"))
//
package bloomfilter

import (
	"log"
	"math"
	"hash/fnv"
	"sync"
)

// HashFunction64 is a function type that takes a byte slice as input and returns hash value of uint64 type.
type HashFunction64 func([]byte) uint64

// BloomFilter is non-thread safe bloom filter data structure.
type BloomFilter struct {
	hash1            HashFunction64
	hash2            HashFunction64
	numHashFunctions uint8
	size             uint64 //in bits
	bits             []uint64
}

// BloomFilterTS is a BloomFilter structure with a RWMutex for thread safety.
type BloomFilterTS struct {
	bf  *BloomFilter
	mtx sync.RWMutex
}

// Add takes a byte slice as input and adds it to the BloomFilter structure's bit array.
func (bf *BloomFilter) Add(data []byte) {
	bitLocations := bf.getBitLocations(data)

	for i := 0; i < len(bitLocations); i++ {
		currLoc := bitLocations[i]
		sliceLoc := (currLoc - (currLoc % 64)) / 64
		bf.bits[sliceLoc] |= (1 << (currLoc % 64))
	}
}

// Query tests the byte slice input's existence in the BloomFilter structure and returns a boolean value. 
// The result is either true for existence or false for inexistence. 
// However it should be noted that false positives are possible, while false negatives are not.
func (bf *BloomFilter) Query(data []byte) bool {
	bitLocations := bf.getBitLocations(data)

	for i := 0; i < len(bitLocations); i++ {
		currLoc := bitLocations[i]
		sliceLoc := (currLoc - (currLoc % 64)) / 64
		if bf.bits[sliceLoc]&(1<<(currLoc%64)) == 0 {
			return false
		}
	}
	return true
}

// Add for thread safe BloomFilterTS structure serves the same purpose as Add for BloomFilter structure.
// Structure is locked for
func (bfts *BloomFilterTS) Add(data []byte) {
	bfts.mtx.Lock()
	bfts.bf.Add(data)
	bfts.mtx.Unlock()
}

// Query for thread safe BloomFilterTS structure serves the same purpose as Query for BloomFilter structure.
func (bfts *BloomFilterTS) Query(data []byte) bool {
	bfts.mtx.RLock()
	defer bfts.mtx.RUnlock()
	retVal := bfts.bf.Query(data)
	return retVal
}

// NewByEstimates requires estimated number of items and estimated false positive rate to create a BloomFilter structure.
// This function calculates size in bits and ideal number of hash functions that will be created by double hashing of 
// hash function hash1 and hash function hash2.
// HashFunction64 hash1 and HashFunction64 hash2 can be nil and when they are nil, a default HashFunction64 for each will be used.
func NewByEstimates(numItems uint64, fpRate float64, hash1 HashFunction64, hash2 HashFunction64) *BloomFilter {
	if (numItems == 0) {
		log.Fatalln("number of items must be positive")
	}
	size := uint64(math.Ceil(-1 * float64(numItems) * math.Log(fpRate) / math.Pow(math.Log(2), 2)))
	numHashFunctions := uint8(math.Ceil(math.Log(2) * float64(size) / float64(numItems)))
	
	return NewBySizeAndNumHashFuncs(size, numHashFunctions, hash1, hash2)
}

func defaultHash1() func(input []byte) uint64 {
	return func(input []byte) uint64 {
		h := fnv.New64a()
		h.Write(input)
		return h.Sum64()
	}
}

func defaultHash2() func(input []byte) uint64 {
	return func(input []byte) uint64 {
		h := fnv.New64()
		h.Write(input)
		return h.Sum64()
	}
}

// NewBySizeAndNumHashFuncs requires maximum size in bits and number of hash functions that will be created via double hashing of
// hash function hash1 and hash function hash2.
// HashFunction64 hash1 and HashFunction64 hash2 can be nil and when they are nil, a default HashFunction64 for each will be used.
// This function returns a new BloomFilter structure.
func NewBySizeAndNumHashFuncs(size uint64, numHashFunctions uint8, hash1 HashFunction64, hash2 HashFunction64) *BloomFilter {
	if hash1 == nil {
		hash1 = defaultHash1()
	}

	if hash2 == nil {
		hash2 = defaultHash2()
	}

	if size == 0 {
		size = uint64(8 * 1024 * 1024 * 100) // 10 mb in bits
	}

	l := (size - (size % 64)) / 64

	if size%64 > 0 {
		l++
	}

	bits := make([]uint64, l, l)

	if numHashFunctions == uint8(0) {
		numHashFunctions = uint8(11)
	}

	bf := BloomFilter{
		hash1:            hash1,
		hash2:            hash2,
		numHashFunctions: numHashFunctions,
		size:             size,
		bits:             bits,
	}

	return &bf
}

// NewTSByEstimates returns a new BloomFilterTS structure. For more details, please see NewByEstimates function.
func NewTSByEstimates(numItems uint64, fpRate float64, hash1 HashFunction64, hash2 HashFunction64) *BloomFilterTS {
	size := uint64(math.Ceil(-1 * float64(numItems) * math.Log(fpRate) / math.Pow(math.Log(2), 2)))
	numHashFunctions := uint8(math.Ceil(math.Log(2) * float64(size) / float64(numItems)))

	bf := NewBySizeAndNumHashFuncs(size, numHashFunctions, hash1, hash2)

	return &BloomFilterTS{bf: bf}
}

// NewTSBySizeAndNumHashFuncs returns a new BloomFilterTS structure. For more details, please see NewBySizeAndNumHashFuncs function.
func NewTSBySizeAndNumHashFuncs(size uint64, numHashFunctions uint8, hash1 HashFunction64, hash2 HashFunction64) *BloomFilterTS {
	bf := NewBySizeAndNumHashFuncs(size, numHashFunctions, hash1, hash2)

	return &BloomFilterTS{bf: bf}
}

func (bf *BloomFilter) getBitLocations(data []byte) []uint64 {
	hash1Val := bf.hash1(data)
	hash2Val := bf.hash2(data)

	retVal := make([]uint64, bf.numHashFunctions)

	for i := uint8(0); i < bf.numHashFunctions; i++ {
		retVal[i] = (hash1Val + uint64(i)*hash2Val) % (bf.size)
	}

	return retVal
}
