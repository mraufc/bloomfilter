Bloom Filter
-------------

A Bloom Filter is a space-efficient probabilistic data structure, conceived by 
Burton Howard Bloom in 1970, that is used to test whether an element is a member of a set.
It is possible to add an elemen to a bloom filter or query an element for existence, 
while removing an element is not possible.
When testing (querying) an element in the bloom filter, either true or false is returned.
However it should be noted that false positives are possible, while false negatives are not.
For additional information see: 
Bloom filter - Wikipedia https://en.wikipedia.org/wiki/Bloom_filter

This implementation of bloomfilter provides two different ways of creating two different
bloom filter data structures.

It is possible to create a new bloom filter data structure by providing either estimated number of items
that the bloom filter will hold and estimated false positove error rate along with two custom hash functions:

    bf := NewByEstimates(numItems uint64, fpRate float64, hash1 HashFunction64, hash2 HashFunction64)

or maximum bloom filter size (in bits) and number or hash functions that will be created by double hashing
of two hash functions:

    bf := NewBySizeAndNumHashFuncs(size uint64, numHashFunctions uint8, hash1 HashFunction64, hash2 HashFunction64)

In both cases, when hash1 and/or hash2 values are nil, default hash function implementations of FNV-1 and/or FNV-1a
from the standard library will be used.

Both of the above bloom filter data structures are non-thread safe, however it is possible to create a thread safe
implementation by:

    bfts := NewTSByEstimates(numItems uint64, fpRate float64, hash1 HashFunction64, hash2 HashFunction64)

or

    bfts := NewTSBySizeAndNumHashFuncs(size uint64, numHashFunctions uint8, hash1 HashFunction64, hash2 HashFunction64)

Once a bloom filter structure is created, one can add an element by;

    bf.Add([]byte("data"))

and test an element for existence by;

    exists := bf.Query([]byte("data"))

Installation
-------------

     go get -u github.com/mraufc/bloomfilter
