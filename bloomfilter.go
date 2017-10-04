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
package main

import (
	"fmt"
	"sync"
	"math"
	"strconv"
	"time"
)

type HashFunction func([]byte) uint64

type BloomFilter struct {
	mtx sync.RWMutex
	Hash1 HashFunction
	Hash2 HashFunction
	NumHashFunctions uint8
	Size uint64				//in bits
	bits []uint64
}

func main() {
	x := uint64(0)
	fmt.Println(x)

	x |= (1 << 0)
	fmt.Println(x)

	x |= (1 << 2)
	fmt.Println(x)

	x &^= (1 << 7)

	fmt.Println(x)
	x &^= (1 << 2)
	
		fmt.Println(x)
		x &^= (1 << 0)
		
			fmt.Println(x)

	aaa := strconv.FormatUint(math.MaxUint64, 2)
	fmt.Println(aaa)

	bf := NewBloomFilter(nil, nil, 0, 0)

	bf.Add([]byte("a"))
	bf.Add([]byte("b"))
	bf.Add([]byte("c"))
	bf.Add([]byte("d"))
	bf.Add([]byte("e"))
	bf.Add([]byte("f"))
	bf.Add([]byte("g"))

	result := bf.Contains([]byte("a"))
	fmt.Println(result)

	result = bf.Contains([]byte("b"))
	fmt.Println(result)
	result = bf.Contains([]byte("c"))
	fmt.Println(result)
	result = bf.Contains([]byte("d"))
	fmt.Println(result)
	result = bf.Contains([]byte("e"))
	fmt.Println(result)
	result = bf.Contains([]byte("f"))
	fmt.Println(result)
	result = bf.Contains([]byte("g"))
	fmt.Println(result)
	result = bf.Contains([]byte("h"))
	fmt.Println(result)
	result = bf.Contains([]byte("j"))
	fmt.Println(result)

	result = bf.Contains([]byte("rauf"))
	fmt.Println(result)

	bf.Add([]byte("rauf"))

	result = bf.Contains([]byte("rauf"))
	fmt.Println(result)



	tix := time.Tick(1 * time.Second)
	for {
		select {
		case <-tix :
			result := bf.Contains([]byte("rauf"))
			fmt.Println(result)
		default:
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func NewBloomFilter(hash1 HashFunction, hash2 HashFunction, size uint64, numHashFunctions uint8) *BloomFilter {
		if hash1 == nil {
			hash1 = func(input []byte) uint64 {
				retVal := uint64(0)
				for i := 0; i < len(input); i++ {
					retVal = retVal * 101 + uint64(input[i])
				}
				return retVal
			}
		}

		if hash2 == nil {
			hash2 = func(input []byte) uint64 {
				retVal := uint64(0)
				for i := 0; i < len(input); i++ {
					retVal = retVal * 31 + uint64(input[i])
				}
				return retVal
			}
		}
		
		if size == 0 {
			size = uint64(8 * 1024 * 1024 * 1024 * 40) // 100 mb in bits
		}

		l := (size - (size % 64)) / 64

		if size % 64 > 0 {
			l++
		}

		fmt.Printf("l -> %v\n", l)

		bits := make([]uint64, l)

		if numHashFunctions == uint8(0) {
			numHashFunctions = uint8(11)
		}
		
		bf := BloomFilter{
			Hash1 : hash1, 
			Hash2 : hash2,
			NumHashFunctions : numHashFunctions,
			Size : size,
			bits : bits,
		}

		return &bf
}

func (bf *BloomFilter) getBitLocations(data []byte) []uint64 {
	hash1Val := bf.Hash1(data)
	hash2Val := bf.Hash2(data)

	retVal := make([]uint64, bf.NumHashFunctions)
	
	for i := uint8(0); i < bf.NumHashFunctions; i++ {
		retVal[i] = (hash1Val + uint64(i) * hash2Val) % (bf.Size / 64)
	}
	// fmt.Print(string(data))
	// fmt.Print(" ")
	// fmt.Print(retVal)
	// fmt.Println("")
	
	return retVal
}

func (bf *BloomFilter) Add(data []byte) {
	bf.mtx.Lock()
	bitLocations := bf.getBitLocations(data)

	for i := 0; i < len(bitLocations); i++ {
		currLoc := bitLocations[i]
		sliceLoc := (currLoc - (currLoc % 64)) / 64
		bf.bits[sliceLoc] |= (1 << (currLoc % 64))
	}
	bf.mtx.Unlock()
}

func (bf *BloomFilter) Contains(data []byte) bool {
	bf.mtx.RLock()
	bitLocations := bf.getBitLocations(data)
	
	for i := 0; i < len(bitLocations); i++ {
		currLoc := bitLocations[i]
		sliceLoc := (currLoc - (currLoc % 64)) / 64
		fmt.Printf("bf.size -> %v bf.bits.len -> %v currLoc -> %v sliceLoc -> %v\n", bf.Size, len(bf.bits), currLoc, sliceLoc)
		if bf.bits[sliceLoc] & (1 << (currLoc % 64)) == 0 {
			bf.mtx.RUnlock()
			return false
		}
	}
	bf.mtx.RUnlock()
	return true
}