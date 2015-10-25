package cuckoo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	size := 8
	Init(size)

	testData := []struct {
		name string
		data string
	}{
		{
			name: "Test Hash -- test",
			data: "test",
		},
		{
			name: "Test Hash -- random",
			data: "random",
		},
	}

	for _, test := range testData {
		fp := fingerprint(test.data)
		a := hash(test.data)
		b := a ^ hash(fp)
		assert.NotEqual(t, a, b)
	}
}

func TestInsert(t *testing.T) {
	size := 500
	Init(size)
	for i := 0; i < size; i++ {
		ok := Insert(string(i))
		assert.True(t, ok)
	}
}

func TestLookUp(t *testing.T) {
	size := 500
	Init(size)
	for i := 0; i < size; i++ {
		ok := Insert(string(i))
		assert.True(t, ok)
	}

	for i := 0; i < size; i++ {
		ok := LookUp(string(i))
		assert.True(t, ok)
	}
}

func TestLookUpFail(t *testing.T) {
	size := 500
	Init(size)
	for i := 0; i < size; i++ {
		ok := Insert(string(i))
		assert.True(t, ok)
	}

	for i := size; i < 1000; i++ {
		ok := LookUp(string(i))
		assert.False(t, ok)
	}
}

func TestDelete(t *testing.T) {
	size := 100
	Init(size)

	ok := Insert("test")
	assert.True(t, ok)

	ok = LookUp("test")
	assert.True(t, ok)

	ok = Delete("test")
	assert.True(t, ok)

	ok = LookUp("test")
	assert.False(t, ok)
}

func TestDeleteFail(t *testing.T) {
	size := 100
	Init(size)

	ok := Delete("test")
	assert.False(t, ok)
}
