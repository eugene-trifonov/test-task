package utils

import (
	"strconv"
	"testing"
)

func TestDictGetValue(t *testing.T) {
	dict := NewDict()

	dict.Put("A", "B")
	result := dict.Get("A")
	if result != "B" {
		t.Error("Unexpected behavior of Get function")
	}

	result = dict.Get("B")
	if result != nil {
		t.Error("Unexpected behavior of Get function")
	}
}

func BenchmarkGetDictValue(b *testing.B) {
	dict := NewDict()
	for i := 0; i < 1000; i++ {
		dict.Put(strconv.Itoa(i), i)
	}

	for n := 0; n < b.N; n++ {
		dict.Get("55")
	}
}

func TestDictPutValue(t *testing.T) {
	dict := NewDict()

	result := dict.Put("A", "B")
	if result != nil {
		t.Error("Unexpected behavior of Put function")
	}

	result = dict.Put("A", "C")
	if result != "B" {
		t.Error("Unexpected behavior of Put function")
	}

	result = dict.PutIfAbsent("A", "D")
	if result != "C" {
		t.Error("Unexpected behavior of PutIfAbsent function")
	}

	result = dict.Put("A", "E")
	if result != "C" {
		t.Error("Unexpected behavior of Put function")
	}
}

func BenchmarkPutDictValue(b *testing.B) {
	dict := NewDict()

	for n := 0; n < b.N; n++ {
		dict.Put(strconv.Itoa(n), n)
	}
}

func TestDictReplaceValue(t *testing.T) {
	dict := NewDict()

	result := dict.Replace("A", "B")
	if result != nil {
		t.Error("Unexpected behavior of Replace function")
	}

	dict.Put("A", "C")
	result = dict.Replace("A", "D")
	if result != "C" {
		t.Error("Unexpected behavior of Replace function")
	}

	if dict.ReplaceValue("A", "B", "C") {
		t.Error("Unexpected behavior of ReplaceValue function")
	}

	if !dict.ReplaceValue("A", "D", "E") {
		t.Error("Unexpected behavior of Replace function")
	}
}

func BenchmarkReplaceDict(b *testing.B) {
	dict := NewDict()
	for i := 0; i < 1000; i++ {
		dict.Put(strconv.Itoa(i), i)
	}

	for n := 0; n < b.N; n++ {
		dict.Replace("55", n)
	}
}

func TestDictGetKeys(t *testing.T) {
	dict := NewDict()

	dict.Put("A", "B")
	dict.Put("C", "D")
	dict.Put("E", "F")

	keys := dict.GetKeys()

	if len(keys) != 3 || dict.Get(keys[0]) == nil || dict.Get(keys[1]) == nil || dict.Get(keys[2]) == nil {
		t.Error("Unexpected behavior of GetKeys function")
	}

	dict.Remove("C")

	keys = dict.GetKeys()
	if len(keys) != 2 || dict.Get(keys[0]) == nil || dict.Get(keys[1]) == nil {
		t.Error("Unexpected behavior of GetKeys function")
	}

	dict.Remove("A")
	keys = dict.GetKeys()
	if len(keys) != 1 || dict.Get(keys[0]) == nil {
		t.Error("Unexpected behavior of GetKeys function")
	}

	dict.Remove("E")
	keys = dict.GetKeys()
	if len(keys) != 0 {
		t.Error("Unexpected behavior of GetKeys function")
	}
}

func BenchmarkGetKeysDict(b *testing.B) {
	dict := NewDict()
	for i := 0; i < 1000; i++ {
		dict.Put(strconv.Itoa(i), i)
	}

	for n := 0; n < b.N; n++ {
		dict.GetKeys()
	}
}

func TestDictRemoveValues(t *testing.T) {
	dict := NewDict()

	dict.Put("A", "1")
	dict.Put("B", "2")
	dict.Put("C", "3")

	result := dict.Remove("D")
	if result != nil {
		t.Error("Unexpected behavior of Remove function")
	}

	result = dict.Remove("A")
	if result != "1" {
		t.Error("Unexpected behavior of Remove function")
	}

	if dict.RemovePair("B", "5") {
		t.Error("Unexpected behavior of RemovePair function")
	}

	if !dict.RemovePair("B", "2") {
		t.Error("Unexpected behavior of RemovePair function")
	}

	result = dict.Remove("C")
	if result != "3" {
		t.Error("Unexpected behavior of Remove function")
	}
}
