package cache

import (
	"strconv"
	"testing"
)

func TestPutGetValues(t *testing.T) {
	cache := NewCache()
	result := cache.Get("A")
	if result != nil {
		t.Error("Result of Get method on empty cache should be <nil>")
	}

	cache.Put("A", "B")
	cache.Put("C", "D")
	result = cache.Get("B")
	if result != nil {
		t.Error("Result of Get method on empty cache should be <nil>")
	}

	result = cache.Get("A")
	if result != "B" {
		t.Error("Unexpected result found")
	}

	result = cache.Get("C")
	if result != "D" {
		t.Error("Unexpected result found")
	}
}

func BenchmarkGetValues(b *testing.B) {
	cache := NewCache()
	for i := 0; i < 1000; i++ {
		cache.Put(strconv.Itoa(i), i)
	}

	for n := 0; n < b.N; n++ {
		cache.Get("50")
	}
}

func TestPutMethods(t *testing.T) {
	cache := NewCache()
	result := cache.Put("A", "B")
	if result != nil {
		t.Error("The result must be <nil> for non-existing key")
	}

	result = cache.Put("C", "D")
	if result != nil {
		t.Error("The result must be <nil> for non-existing key")
	}

	result = cache.PutIfAbsent("C", "E")
	if result != "D" {
		t.Error("Unexpected behavior of PutIfAbsent function")
	}

	result = cache.Put("C", "E")
	if result != "D" {
		t.Error("Unexpected behavior of Put function")
	}

	result = cache.PutIfAbsent("F", "G")
	if result != nil {
		t.Error("Unexpected behavior of PutIfAbsent function")
	}
}

func BenchmarkPutMethod(b *testing.B) {
	cache := NewCache()
	for n := 0; n < b.N; n++ {
		cache.Put(strconv.Itoa(n), n)
	}
}

func TestReplacingMethods(t *testing.T) {
	cache := NewCache()

	cache.Put("A", "B")
	cache.Put("C", "D")
	cache.Put("E", "F")

	result := cache.Replace("G", "H")
	if result != nil {
		t.Error("Unexpected behavior of Replace function")
	}

	result = cache.Replace("A", "Z")
	if result != "B" {
		t.Error("Unexpected behavior of Replace function")
	}

	replaced := cache.ReplaceValue("C", "Q", "X")
	if replaced {
		t.Error("Unexpected behavior of Replace function")
	}

	replaced = cache.ReplaceValue("C", "D", "X")
	if !replaced {
		t.Error("Unexpected behavior of Replace function")
	}
}

func BenchmarkReplaceMethod(b *testing.B) {
	cache := NewCache()
	cache.Put("A", "B")

	for n := 0; n < b.N; n++ {
		cache.Replace("A", strconv.Itoa(n))
	}
}

func TestSizeMethod(t *testing.T) {
	cache := NewCache()

	cache.Put("A", "B")
	cache.Put("C", "D")
	cache.Put("E", "F")

	if cache.Size() != 3 {
		t.Error("Wrong size calculation found")
	}

	cache.Put("A", "Z")
	cache.Replace("C", "X")

	if cache.Size() != 3 {
		t.Error("Wrong size calculation found")
	}

	cache.PutIfAbsent("G", "H")
	if cache.Size() != 4 {
		t.Error("Wrong size calculation found")
	}

	cache.Remove("A")
	cache.Remove("C")
	if cache.Size() != 2 {
		t.Error("Wrong size calculation found")
	}

	cache.Remove("E")
	cache.Remove("G")
	if cache.Size() != 0 {
		t.Error("Wrong size calculation found")
	}
}

func TestRemoveMethods(t *testing.T) {
	cache := NewCache()

	cache.Put("A", "B")
	cache.Put("C", "D")

	result := cache.Remove("E")
	if result != nil {
		t.Error("Wrong behavior of Remove function")
	}

	result = cache.Remove("A")
	if result != "B" {
		t.Error("Wrong behavior of Remove function")
	}

	removed := cache.RemovePair("C", "F")
	if removed {
		t.Error("Wrong behavior of Remove function")
	}

	removed = cache.RemovePair("C", "D")
	if !removed {
		t.Error("Wrong behavior of Remove function")
	}
}

func TestGetKeysMethods(t *testing.T) {
	cache := NewCache()

	cache.Put("1", "2")
	cache.Put("3", "4")

	keys := cache.GetKeys()
	if len(keys) != 2 || cache.Get(keys[0]) == nil || cache.Get(keys[1]) == nil {
		t.Error("Wrong behavior of GetKeys function")
	}

	cache.Remove("1")
	cache.Remove("3")

	keys = cache.GetKeys()
	if len(keys) != 0 {
		t.Error("Wrong behavior of Remove function")
	}
}

func BenchmarkGetKeysMethod(b *testing.B) {
	cache := NewCache()
	for i := 0; i < 1000; i++ {
		cache.Put(strconv.Itoa(i), i)
	}

	for n := 0; n < b.N; n++ {
		cache.GetKeys()
	}
}
