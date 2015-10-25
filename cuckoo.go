package cuckoo

import (
	"crypto/sha256"
	"encoding/base64"
	"hash/fnv"
	"math/rand"
	"sync"
)

type entry []string

var (
	bucketSize int

	bucketLock sync.RWMutex
	buckets    []entry
)

const maxNumKicks = 500

func fingerprint(x string) string {
	newHash := sha256.New()
	return base64.StdEncoding.EncodeToString(newHash.Sum([]byte(x)))
}

func hash(x string) int {
	newHash := fnv.New32a()
	_, err := newHash.Write([]byte(x))
	if err != nil {
		panic(err)
	}
	return int(newHash.Sum32()) % bucketSize
}

// hasSpace returns an index pair if buckets[a] or buckets[b] has space to insert one more new value
func hasSpace(indexes ...int) (int, int) {
	for _, index := range indexes {
		for i := range buckets[index] {
			if buckets[index][i] == "" {
				return index, i
			}
		}
	}
	return -1, -1
}

// hasValue checks if the specific value in one of the bucket by given indexes
func hasValue(value string, indexes ...int) bool {
	for _, index := range indexes {
		if index >= len(buckets) {
			continue
		}
		for _, fp := range buckets[index] {
			if fp == value {
				return true
			}
		}
	}
	return false
}

// removeValue removes a value from buckets if it exists
func removeValue(value string, indexes ...int) bool {
	for _, index := range indexes {
		if index >= len(buckets) {
			continue
		}
		for i := range buckets[index] {
			if buckets[index][i] == value {
				buckets[index][i] = ""
				return true
			}
		}
	}

	return false
}

func addToBucket(index, candidate int, fp string) { buckets[index][candidate] = fp }

func randIndex(index ...int) int { return rand.Intn(len(index)) }

func randEntry(index int) int { return rand.Intn(len(buckets[index])) }

// Init initializes cuckoo filter buckets
func Init(size int) {
	bucketSize = size
	buckets = make([]entry, size)
	for i := range buckets {
		buckets[i] = make([]string, 4)
	}
}

// Insert inserts a new value into hash map
func Insert(value string) bool {
	fp := fingerprint(value)
	hashIndexOne := hash(value)
	hashIndexTwo := hashIndexOne ^ hash(fp)

	bucketLock.Lock()
	defer bucketLock.Unlock()

	if index, candidate := hasSpace(hashIndexOne, hashIndexTwo); index != -1 {
		addToBucket(index, candidate, fp)
		return true
	}
	// Must relocate existing items
	index := randIndex(hashIndexOne, hashIndexTwo)
	for n := 0; n < maxNumKicks; n++ {
		ent := randEntry(index)
		fp, buckets[index][ent] = buckets[index][ent], fp
		index = index ^ hash(fp)
		if index, candidate := hasSpace(index); index != -1 {
			addToBucket(index, candidate, fp)
			return true
		}
	}
	// Hash table is considered full
	return false
}

// LookUp checks if a given key in the hash map
func LookUp(value string) bool {
	fp := fingerprint(value)
	hashIndexOne := hash(value)
	hashIndexTwo := hashIndexOne ^ hash(fp)

	bucketLock.RLock()
	defer bucketLock.RUnlock()

	return hasValue(fp, hashIndexOne, hashIndexTwo)
}

// Delete removes an existing element from buckets if it exists
func Delete(value string) bool {
	fp := fingerprint(value)
	hashIndexOne := hash(value)
	hashIndexTwo := hashIndexOne ^ hash(fp)

	bucketLock.Lock()
	defer bucketLock.Unlock()

	return removeValue(fp, hashIndexOne, hashIndexTwo)
}
