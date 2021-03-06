package bloomfilter

import (
	"fmt"
	"testing"
	"math/rand"
	"time"
)

const (
	letterBytes string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%&*() "
	acceptableAdditionalFalsePositiveErrorRate float64 = 0.5
)

type testCase struct {
	description string
	data []byte
}

func TestBloomFilterInit(t *testing.T) {
	type initByEstimates struct {
		numItems uint64
		fpRate float64
	}
	type initBySizeAndNumHashFuncs struct {
		size uint64
		numHashFunctions uint8
	}
	var testsForEstimates = []initByEstimates{
		{uint64(0), 0.1},
		{uint64(0), 0.01},
		{uint64(0), 0.001},
		{uint64(1), 0.0},
		{uint64(1), 1.0},
		{uint64(10), -0.1},
		{uint64(10), 1.1},
		{uint64(100), 0.0},
		{uint64(100), 1.0},
		{uint64(1000), -0.1},
		{uint64(1000), 1.1},
		{uint64(1), 0.1},
		{uint64(1), 0.01},
		{uint64(1), 0.001},
		{uint64(100), 0.1},
		{uint64(100), 0.01},
		{uint64(100), 0.001},
		{uint64(10000), 0.1},
		{uint64(10000), 0.01},
		{uint64(10000), 0.001},
	}
	var testsForSizeAndNumHashFuncs = []initBySizeAndNumHashFuncs{
		{uint64(0), uint8(1)},
		{uint64(0), uint8(2)},
		{uint64(0), uint8(3)},
		{uint64(0), uint8(4)},
		{uint64(0), uint8(11)},
		{uint64(1), uint8(0)},
		{uint64(10), uint8(0)},
		{uint64(100), uint8(0)},
		{uint64(1000), uint8(0)},
		{uint64(1000000), uint8(0)},
		{uint64(1000000), uint8(1)},
		{uint64(1000000), uint8(2)},
		{uint64(1000000), uint8(3)},
		{uint64(100000000), uint8(11)},
	}
	for i := 0; i < len(testsForEstimates); i++ {
		bf, err := NewByEstimates(testsForEstimates[i].numItems, testsForEstimates[i].fpRate, nil, nil)
		if testsForEstimates[i].numItems == 0 {
			if bf != nil {
				t.Logf("expected nil bloomfilter for values of numItems %v, fpRate %v", testsForEstimates[i].numItems, testsForEstimates[i].fpRate)
				t.Fail()
			}
			if err.Error() != ErrInvalidNumberOfItems.Error() {
				t.Logf("expected error message for values of numItems %v, fpRate %v is %v", testsForEstimates[i].numItems, testsForEstimates[i].fpRate, ErrInvalidNumberOfItems.Error())
				t.Fail()
			}
		} else if testsForEstimates[i].fpRate <= 0.0 || testsForEstimates[i].fpRate >= 1.0 {
			if bf != nil {
				t.Logf("expected nil bloomfilter for values of numItems %v, fpRate %v", testsForEstimates[i].numItems, testsForEstimates[i].fpRate)
				t.Fail()
			}
			if err.Error() != ErrInvalidFalsePositiveRate.Error() {
				t.Logf("expected error message for values of numItems %v, fpRate %v is %v", testsForEstimates[i].numItems, testsForEstimates[i].fpRate, ErrInvalidFalsePositiveRate.Error())
				t.Fail()
			}
		} else {
			if bf == nil {
				t.Logf("expected non-nil bloomfilter for values of numItems %v, fpRate %v", testsForEstimates[i].numItems, testsForEstimates[i].fpRate)
				t.Fail()
			}
			if err != nil {
				t.Logf("expected nil error for values of numItems %v, fpRate %v", testsForEstimates[i].numItems, testsForEstimates[i].fpRate)
				t.Fail()
			}
		}
	}
	for i := 0; i < len(testsForSizeAndNumHashFuncs); i++ {
		bf, err := NewBySizeAndNumHashFuncs(testsForSizeAndNumHashFuncs[i].size, testsForSizeAndNumHashFuncs[i].numHashFunctions, nil, nil)
		if testsForSizeAndNumHashFuncs[i].size == 0 {
			if bf != nil {
				t.Logf("expected nil bloomfilter for values of size %v, numHashFunctions %v", testsForSizeAndNumHashFuncs[i].size, testsForSizeAndNumHashFuncs[i].numHashFunctions)
				t.Fail()
			}
			if err.Error() != ErrInvalidSize.Error() {
				t.Logf("expected error message for values of size %v, numHashFunctions %v is %v", testsForSizeAndNumHashFuncs[i].size, testsForSizeAndNumHashFuncs[i].numHashFunctions, ErrInvalidSize.Error())
				t.Fail()
			}
		} else if testsForSizeAndNumHashFuncs[i].numHashFunctions == 0 {
			if bf != nil {
				t.Logf("expected nil bloomfilter for values of size %v, numHashFunctions %v", testsForSizeAndNumHashFuncs[i].size, testsForSizeAndNumHashFuncs[i].numHashFunctions)
				t.Fail()
			}
			if err.Error() != ErrInvalidNumberOfHashFunctions.Error() {
				t.Logf("expected error message for values of size %v, numHashFunctions %v is %v", testsForSizeAndNumHashFuncs[i].size, testsForSizeAndNumHashFuncs[i].numHashFunctions, ErrInvalidNumberOfHashFunctions.Error())
				t.Fail()
			}
		} else {
			if bf == nil {
				t.Logf("expected non-nil bloomfilter for values of size %v, numHashFunctions %v", testsForSizeAndNumHashFuncs[i].size, testsForSizeAndNumHashFuncs[i].numHashFunctions)
				t.Fail()
			}
			if err != nil {
				t.Logf("expected nil error for values of size %v, numHashFunctions %v", testsForSizeAndNumHashFuncs[i].size, testsForSizeAndNumHashFuncs[i].numHashFunctions)
				t.Fail()
			}
		}
	}
	for i := 0; i < len(testsForEstimates); i++ {
		bfts, err := NewTSByEstimates(testsForEstimates[i].numItems, testsForEstimates[i].fpRate, nil, nil)
		if testsForEstimates[i].numItems == 0 {
			if bfts != nil {
				t.Logf("expected nil bloomfilter for values of numItems %v, fpRate %v", testsForEstimates[i].numItems, testsForEstimates[i].fpRate)
				t.Fail()
			}
			if err.Error() != ErrInvalidNumberOfItems.Error() {
				t.Logf("expected error message for values of numItems %v, fpRate %v is %v", testsForEstimates[i].numItems, testsForEstimates[i].fpRate, ErrInvalidNumberOfItems.Error())
				t.Fail()
			}
		} else if testsForEstimates[i].fpRate <= 0.0 || testsForEstimates[i].fpRate >= 1.0 {
			if bfts != nil {
				t.Logf("expected nil bloomfilter for values of numItems %v, fpRate %v", testsForEstimates[i].numItems, testsForEstimates[i].fpRate)
				t.Fail()
			}
			if err.Error() != ErrInvalidFalsePositiveRate.Error() {
				t.Logf("expected error message for values of numItems %v, fpRate %v is %v", testsForEstimates[i].numItems, testsForEstimates[i].fpRate, ErrInvalidFalsePositiveRate.Error())
				t.Fail()
			}
		} else {
			if bfts == nil {
				t.Logf("expected non-nil bloomfilter for values of numItems %v, fpRate %v", testsForEstimates[i].numItems, testsForEstimates[i].fpRate)
				t.Fail()
			}
			if err != nil {
				t.Logf("expected nil error for values of numItems %v, fpRate %v", testsForEstimates[i].numItems, testsForEstimates[i].fpRate)
				t.Fail()
			}
		}
	}
	for i := 0; i < len(testsForSizeAndNumHashFuncs); i++ {
		bfts, err := NewTSBySizeAndNumHashFuncs(testsForSizeAndNumHashFuncs[i].size, testsForSizeAndNumHashFuncs[i].numHashFunctions, nil, nil)
		if testsForSizeAndNumHashFuncs[i].size == 0 {
			if bfts != nil {
				t.Logf("expected nil bloomfilter for values of size %v, numHashFunctions %v", testsForSizeAndNumHashFuncs[i].size, testsForSizeAndNumHashFuncs[i].numHashFunctions)
				t.Fail()
			}
			if err.Error() != ErrInvalidSize.Error() {
				t.Logf("expected error message for values of size %v, numHashFunctions %v is %v", testsForSizeAndNumHashFuncs[i].size, testsForSizeAndNumHashFuncs[i].numHashFunctions, ErrInvalidSize.Error())
				t.Fail()
			}
		} else if testsForSizeAndNumHashFuncs[i].numHashFunctions == 0 {
			if bfts != nil {
				t.Logf("expected nil bloomfilter for values of size %v, numHashFunctions %v", testsForSizeAndNumHashFuncs[i].size, testsForSizeAndNumHashFuncs[i].numHashFunctions)
				t.Fail()
			}
			if err.Error() != ErrInvalidNumberOfHashFunctions.Error() {
				t.Logf("expected error message for values of size %v, numHashFunctions %v is %v", testsForSizeAndNumHashFuncs[i].size, testsForSizeAndNumHashFuncs[i].numHashFunctions, ErrInvalidNumberOfHashFunctions.Error())
				t.Fail()
			}
		} else {
			if bfts == nil {
				t.Logf("expected non-nil bloomfilter for values of size %v, numHashFunctions %v", testsForSizeAndNumHashFuncs[i].size, testsForSizeAndNumHashFuncs[i].numHashFunctions)
				t.Fail()
			}
			if err != nil {
				t.Logf("expected nil error for values of size %v, numHashFunctions %v", testsForSizeAndNumHashFuncs[i].size, testsForSizeAndNumHashFuncs[i].numHashFunctions)
				t.Fail()
			}
		}
	}
}
func TestBloomFilterBasics(t *testing.T) {
	var (
		count = 1000000
		numItems = uint64(count)
		fp = 0.01
		maxStrLen = 50
		minStrLen = 20
	)

	tests := prepTestCases(count, minStrLen, maxStrLen)

	bf, err := NewByEstimates(numItems, fp, nil, nil)
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			bf.Add(tt.data)
		})
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			if result := bf.Query(tt.data); !result {
				t.Errorf("Query(%v): expected %v, actual %v", string(tt.data), true, false)
			}
		})
	}
}

func TestBloomFilterTSBasics(t *testing.T) {
	var (
		count = 1000000
		numItems = uint64(count)
		fp = 0.01
		maxStrLen = 50
		minStrLen = 20
	)

	tests := prepTestCases(count, minStrLen, maxStrLen)

	bf, err := NewTSByEstimates(numItems, fp, nil, nil)
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			bf.Add(tt.data)
		})
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			if result := bf.Query(tt.data); !result {
				t.Errorf("Query(%v): expected %v, actual %v", string(tt.data), true, false)
			}
		})
	}
}

// This test should fail when "go test -race" command is issued.
// BloomFilter structure is NOT thread safe.
// func TestBloomFilterParallel(t *testing.T) {
// 	var (
// 		count = 10
// 		numItems = uint64(count)
// 		fp = 0.01
// 		maxStrLen = 50
// 		minStrLen = 20
// 	)
// 
// 	tests := prepTestCases(count, minStrLen, maxStrLen)
// 
// 	bf := NewByEstimates(numItems, fp, nil, nil)
// 
// 	t.Parallel()
// 	for _, tt := range tests {
// 		t.Run(tt.description, func(t *testing.T) {
// 			t.Parallel()
// 			bf.Add(tt.data)
// 		})
// 		t.Run(tt.description, func(t *testing.T) {
// 			t.Parallel()
// 			bf.Query(tt.data)
// 		})
// 	}
// }

// This test should NOT fail when "go test -race" command is issued.
// BloomFilterTS structure is thread safe.
func TestBloomFilterTSParallel(t *testing.T) {
	var (
		count = 10
		numItems = uint64(count)
		fp = 0.01
		maxStrLen = 50
		minStrLen = 20
	)

	tests := prepTestCases(count, minStrLen, maxStrLen)

	bf, err := NewTSByEstimates(numItems, fp, nil, nil)
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			bf.Add(tt.data)
		})
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			bf.Query(tt.data)
		})
	}
}

// Some of the following tests may fail even when an additional acceptable false positive rate is provided
func TestFalsePositiveRate1000_5(t *testing.T)   { testFalsePositiveRate(t, 1000, 0.5) }
func TestFalsePositiveRate10000_5(t *testing.T)   { testFalsePositiveRate(t, 10000, 0.5) }
func TestFalsePositiveRate100000_5(t *testing.T)   { testFalsePositiveRate(t, 100000, 0.5) }
func TestFalsePositiveRate1000000_5(t *testing.T)   { testFalsePositiveRate(t, 1000000, 0.5) }
// func TestFalsePositiveRate1000_1(t *testing.T)   { testFalsePositiveRate(t, 1000, 0.1) }
// func TestFalsePositiveRate10000_1(t *testing.T)   { testFalsePositiveRate(t, 10000, 0.1) }
// func TestFalsePositiveRate100000_1(t *testing.T)   { testFalsePositiveRate(t, 100000, 0.1) }
// func TestFalsePositiveRate1000000_1(t *testing.T)   { testFalsePositiveRate(t, 1000000, 0.1) }
// func TestFalsePositiveRate1000_01(t *testing.T)   { testFalsePositiveRate(t, 1000, 0.01) }
// func TestFalsePositiveRate10000_01(t *testing.T)   { testFalsePositiveRate(t, 10000, 0.01) }
// func TestFalsePositiveRate100000_01(t *testing.T)   { testFalsePositiveRate(t, 100000, 0.01) }
// func TestFalsePositiveRate1000000_01(t *testing.T)   { testFalsePositiveRate(t, 1000000, 0.01) }
// func TestFalsePositiveRate1000_001(t *testing.T)   { testFalsePositiveRate(t, 1000, 0.001) }
// func TestFalsePositiveRate10000_001(t *testing.T)   { testFalsePositiveRate(t, 10000, 0.001) }
// func TestFalsePositiveRate100000_001(t *testing.T)   { testFalsePositiveRate(t, 100000, 0.001) }
// func TestFalsePositiveRate1000000_001(t *testing.T)   { testFalsePositiveRate(t, 1000000, 0.001) }
// func TestFalsePositiveRate1000_0001(t *testing.T)   { testFalsePositiveRate(t, 1000, 0.0001) }
// func TestFalsePositiveRate10000_0001(t *testing.T)   { testFalsePositiveRate(t, 10000, 0.0001) }
// func TestFalsePositiveRate100000_0001(t *testing.T)   { testFalsePositiveRate(t, 100000, 0.0001) }
// func TestFalsePositiveRate1000000_0001(t *testing.T)   { testFalsePositiveRate(t, 1000000, 0.0001) }

// TODO: consider adding a distrubution test for default hash functions

func BenchmarkAdd(t *testing.B) {
	var (
		count = 1000000
		numItems = uint64(count)
		fp = 0.01
		maxStrLen = 50
		minStrLen = 20
	)

	bf, err := NewByEstimates(numItems, fp, nil, nil)
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}

	tests := prepTestCases(count, minStrLen, maxStrLen)

	t.ResetTimer()

	for _, tt := range tests {
		bf.Add(tt.data)
	}
	
}

func BenchmarkQueryEmptyBF(t *testing.B) {
	var (
		count = 1000000
		numItems = uint64(count)
		fp = 0.01
		maxStrLen = 50
		minStrLen = 20
	)

	bf, err := NewByEstimates(numItems, fp, nil, nil)
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}

	tests := prepTestCases(count, minStrLen, maxStrLen)

	t.ResetTimer()
	for _, tt := range tests {
		bf.Query(tt.data)
	}
}

func BenchmarkQuery(t *testing.B) {
	var (
		count = 1000000
		numItems = uint64(count)
		fp = 0.01
		maxStrLen = 50
		minStrLen = 20
	)

	bf, err := NewByEstimates(numItems, fp, nil, nil)
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}

	tests := prepTestCases(count, minStrLen, maxStrLen)

	for _, tt := range tests {
		bf.Add(tt.data)
	}
	
	t.ResetTimer()

	for _, tt := range tests {
		bf.Query(tt.data)
	}
	
}

func testFalsePositiveRate(t *testing.T, count int, fp float64) {
	var (
		numItems = uint64(count)
		maxStrLen = 40
		minStrLen = 30
	)

	bf, err := NewByEstimates(numItems, fp, nil, nil)
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}

	// input strings (actually byte slices) are random, 
	// for uniqueness and false positive rate testing,
	// let's make use of a map.
	// this test may not be considered to be easy on memory.
	mT := make(map[string]bool)

	tests := prepTestCases(count, minStrLen, maxStrLen)

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			bf.Add(tt.data)
			mT[string(tt.data)] = true
		})
	}

	testFP := prepTestCases(count, minStrLen, maxStrLen)

	totalCount := 0
	fpCount := 0

	for _, tt := range testFP {
		t.Run(tt.description, func(t *testing.T) {
			totalCount++
			if result := bf.Query(tt.data); result {
				if _, ok := mT[string(tt.data)]; !ok {
					fpCount++
				}
			}
		})
	}

	actualFpRate := float64(fpCount) / float64(totalCount)
	acceptableFpRate := (1 + acceptableAdditionalFalsePositiveErrorRate) * fp
	if actualFpRate > acceptableFpRate {
		t.Errorf("expected false positive rate is %v, acceptable is %v, actual is %v - %v out of %v items\n", fp, acceptableFpRate, actualFpRate, fpCount, totalCount)
	}
}

func randBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return b
}

func prepTestCases(count, minStrLen, maxStrLen int) []testCase {
	rand.Seed(time.Now().UnixNano())
	
	testCases := make([]testCase, count)

	for i := 0; i < count; i++ {
		s := randBytes(rand.Intn(maxStrLen - minStrLen + 1) + minStrLen)
		testCases[i].data = s
		testCases[i].description = fmt.Sprintf("test case: %v, data: %v", i+1, string(s))
	}

	return testCases
}