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
	"math/rand"
	// "time"
)

type HashFunction func([]byte) uint64

type BloomFilter struct {
	hash1 HashFunction
	hash2 HashFunction
	numHashFunctions uint8
	size uint64				//in bits
	bits []uint64
	threadSafe bool
}

type BloomFilterTS struct {
	bf *BloomFilter
	mtx sync.RWMutex
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

	// bf := NewBySizeAndNumHashFuncs(0, 0, nil, nil)

	bf := NewByEstimates(90000000, 0.01, nil, nil)
	bf.Add([]byte("a"))
	bf.Add([]byte("b"))
	bf.Add([]byte("c"))
	bf.Add([]byte("d"))
	bf.Add([]byte("e"))
	bf.Add([]byte("f"))
	bf.Add([]byte("g"))

	result := bf.Query([]byte("a"))
	fmt.Println(result)

	result = bf.Query([]byte("b"))
	fmt.Println(result)
	result = bf.Query([]byte("c"))
	fmt.Println(result)
	result = bf.Query([]byte("d"))
	fmt.Println(result)
	result = bf.Query([]byte("e"))
	fmt.Println(result)
	result = bf.Query([]byte("f"))
	fmt.Println(result)
	result = bf.Query([]byte("g"))
	fmt.Println(result)
	result = bf.Query([]byte("h"))
	fmt.Println(result)
	result = bf.Query([]byte("j"))
	fmt.Println(result)

	result = bf.Query([]byte("rauf"))
	fmt.Println(result)

	bf.Add([]byte("rauf"))

	result = bf.Query([]byte("rauf"))
	fmt.Println(result)

	for x:= 0 ; x<90000000 ; x++{
		str := RandStringBytes( 1000)
		bf.Add([]byte(str))
	}

	for {}



	// tix := time.Tick(1 * time.Second)
	// for {
	// 	select {
	// 	case <-tix :
	// 		result := bf.Query([]byte("rauf"))
	// 		fmt.Println(result)
	// 	default:
	// 		time.Sleep(50 * time.Millisecond)
	// 	}
	// }
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}

func NewTSByEstimates(numItems uint64, fpRate float64, hash1 HashFunction, hash2 HashFunction) *BloomFilterTS {
	size := uint64(math.Ceil(-1 * float64(numItems) * math.Log(fpRate) / math.Pow(math.Log(2), 2)))
	numHashFunctions := uint8(math.Ceil(math.Log(2) * float64(size) / float64(numItems)))
	fmt.Printf("size %v, numHashFunctions %v", size, numHashFunctions)
	
	bf := NewBySizeAndNumHashFuncs(size, numHashFunctions, hash1, hash2)

	return &BloomFilterTS{bf : bf}
} 

func NewTSBySizeAndNumHashFuncs(size uint64, numHashFunctions uint8, hash1 HashFunction, hash2 HashFunction) *BloomFilterTS {
	bf := NewBySizeAndNumHashFuncs(size, numHashFunctions, hash1, hash2)

	return &BloomFilterTS{bf : bf}
}

func NewByEstimates(numItems uint64, fpRate float64, hash1 HashFunction, hash2 HashFunction) *BloomFilter {
	size := uint64(math.Ceil(-1 * float64(numItems) * math.Log(fpRate) / math.Pow(math.Log(2), 2)))
	numHashFunctions := uint8(math.Ceil(math.Log(2) * float64(size) / float64(numItems)))
	fmt.Printf("size %v, numHashFunctions %v", size, numHashFunctions)
	return NewBySizeAndNumHashFuncs(size, numHashFunctions, hash1, hash2)
}

func NewBySizeAndNumHashFuncs(size uint64, numHashFunctions uint8, hash1 HashFunction, hash2 HashFunction) *BloomFilter {
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
			size = uint64(8 * 1024 * 1024 * 100) // 10 mb in bits
		}

		l := (size - (size % 64)) / 64

		if size % 64 > 0 {
			l++
		}

		fmt.Printf("size -> %v l -> %v\n", size, l)

		bits := make([]uint64, l, l)

		fmt.Printf("bits len is %v, cap is %v, memory should be %v bytes, %v kbytes, %v mbytes \n", len(bits), cap(bits), len(bits) * 8, len(bits) * 8 / 1024, len(bits) * 8 / 1024 / 1024)

		if numHashFunctions == uint8(0) {
			numHashFunctions = uint8(11)
		}
		
		bf := BloomFilter{
			hash1 : hash1, 
			hash2 : hash2,
			numHashFunctions : numHashFunctions,
			size : size,
			bits : bits,
		}

		return &bf
}



func (bf *BloomFilter) getBitLocations(data []byte) []uint64 {
	hash1Val := bf.hash1(data)
	hash2Val := bf.hash2(data)

	retVal := make([]uint64, bf.numHashFunctions)
	
	for i := uint8(0); i < bf.numHashFunctions; i++ {
		retVal[i] = (hash1Val + uint64(i) * hash2Val) % (bf.size / 64)
	}
	// fmt.Print(string(data))
	// fmt.Print(" ")
	// fmt.Print(retVal)
	// fmt.Println("")
	
	return retVal
}

func (bfts *BloomFilterTS) Add(data []byte) {
	bfts.mtx.Lock()
	bfts.bf.Add(data)
	bfts.mtx.Unlock()
}

func (bfts *BloomFilterTS) Query(data []byte) bool {
	bfts.mtx.RLock()
	defer bfts.mtx.RUnlock()
	retVal := bfts.bf.Query(data)
	return retVal
}

func (bf *BloomFilter) Add(data []byte) {
	bitLocations := bf.getBitLocations(data)

	for i := 0; i < len(bitLocations); i++ {
		currLoc := bitLocations[i]
		sliceLoc := (currLoc - (currLoc % 64)) / 64
		bf.bits[sliceLoc] |= (1 << (currLoc % 64))
	}
}

func (bf *BloomFilter) Query(data []byte) bool {
	bitLocations := bf.getBitLocations(data)
	
	for i := 0; i < len(bitLocations); i++ {
		currLoc := bitLocations[i]
		sliceLoc := (currLoc - (currLoc % 64)) / 64
		fmt.Printf("bf.size -> %v bf.bits.len -> %v currLoc -> %v sliceLoc -> %v\n", bf.size, len(bf.bits), currLoc, sliceLoc)
		if bf.bits[sliceLoc] & (1 << (currLoc % 64)) == 0 {
			return false
		}
	}
	return true
}