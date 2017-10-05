package main

import (
	"fmt"
	"testing"
	"math/rand"
	"time"
)

const (
	letterBytes string = "abcdefghijklmnopqrstuvwxyz ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type testCase struct {
	description string
	data []byte
}

func TestBasics(t *testing.T) {
	const (
		count int = 10000
		numItems uint64 = uint64(count)
		fp float64 = 0.01
		maxStrLen int = 50
		minStrLen int = 20
	)

	tests := prepTestCases(count, minStrLen, maxStrLen)

	bf := NewByEstimates(numItems, fp, nil, nil)

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			bf.Add(tt.data)
		})
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			result := bf.Query(tt.data)
			if !result {
				t.Errorf("Query(%v): expected %v, actual %v", string(tt.data), true, false)
			}
		})
	}
}

func BenchmarkAdd(t *testing.B) {
	const (
		count int = 1000000
		numItems uint64 = uint64(count)
		fp float64 = 0.01
		maxStrLen int = 50
		minStrLen int = 20
	)

	bf := NewByEstimates(numItems, fp, nil, nil)
	tests := prepTestCases(count, minStrLen, maxStrLen)

	t.ResetTimer()

	for _, tt := range tests {
		bf.Add(tt.data)
	}
	
}

func BenchmarkQueryEmptyBF(t *testing.B) {
	const (
		count int = 1000000
		numItems uint64 = uint64(count)
		fp float64 = 0.01
		maxStrLen int = 50
		minStrLen int = 20
	)

	bf := NewByEstimates(numItems, fp, nil, nil)
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

	bf := NewByEstimates(numItems, fp, nil, nil)
	tests := prepTestCases(count, minStrLen, maxStrLen)

	for _, tt := range tests {
		bf.Add(tt.data)
	}
	
	t.ResetTimer()

	for _, tt := range tests {
		bf.Query(tt.data)
	}
	
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

func TestFalsePositiveRate1000_5(t *testing.T)   { testFalsePositiveRate(t, 1000, 0.5) }
func TestFalsePositiveRate10000_5(t *testing.T)   { testFalsePositiveRate(t, 10000, 0.5) }
// func TestFalsePositiveRate100000_5(t *testing.T)   { testFalsePositiveRate(t, 100000, 0.5) }
// func TestFalsePositiveRate1000000_5(t *testing.T)   { testFalsePositiveRate(t, 1000000, 0.5) }
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


func testFalsePositiveRate(t *testing.T, count int, fp float64) {
	numItems := uint64(count)
	
	const (
		maxStrLen int = 40
		minStrLen int = 30
	)

	bf := NewByEstimates(numItems, fp, nil, nil)

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
			result := bf.Query(tt.data)
			if result {
				if _, ok := mT[string(tt.data)]; !ok {
					fpCount++
				}
			}
		})
	}

	actualFpRate := float64(fpCount) / float64(totalCount)

	if actualFpRate > 2.0 * fp {
		t.Fail()
	}

	if t.Failed() {
		t.Errorf("Expected false positive rate is %v, actual is %v - %v out of %v items - %v\n", fp, actualFpRate, fpCount, totalCount, len(mT))
	}
}

func randBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return b
}



