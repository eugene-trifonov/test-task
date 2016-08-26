package utils

import (
	"testing"
)

func TestAppendMethods(t *testing.T) {
	list := NewEmptySyncList()

	if list.Size() != 0 {
		t.Error("Wrong result of Append function")
	}

	list.Append("A").Append("B").Append("C")

	if list.Size() != 3 {
		t.Error("Wrong result of Append function")
	}
}

func BenchmarkAppendMethod(b *testing.B) {
	list := NewEmptySyncList()

	for n := 0; n < b.N; n++ {
		list.Append("A")
	}
}

func TestGetMethods(t *testing.T) {
	list := NewEmptySyncList()

	list.Append("D").Append("E").Append("F")

	result, err := list.Get(-1)
	if err == nil {
		t.Error("Wrong bahavior of Get function")
	}

	result, err = list.Get(5)
	if err == nil {
		t.Error("Wrong bahavior of Get function")
	}

	result, err = list.Get(0)
	if err != nil || result != "D" {
		t.Error("Wrong bahavior of Get function")
	}

	result, err = list.Get(2)
	if err != nil || result != "F" {
		t.Error("Wrong bahavior of Get function")
	}
}

func BenchmarkGet0Method(b *testing.B) {
	benchTestGetMethod(0, b)
}

func BenchmarkGet8Method(b *testing.B) {
	benchTestGetMethod(8, b)
}

func benchTestGetMethod(i int, b *testing.B) {
	list := NewSyncList("A", "B", "C", "D", "E", "F", "G", "H", "I")

	for n := 0; n < b.N; n++ {
		list.Get(i)
	}
}

func TestRemoveMethods(t *testing.T) {
	list := NewEmptySyncList()

	list.Append("G").Append("H").Append("I").Append("J")

	result, err := list.Remove(-1)
	if err == nil {
		t.Error("Wrong bahavior of Remove function")
	}

	result, err = list.Remove(4)
	if err == nil {
		t.Error("Wrong bahavior of Remove function")
	}

	result, err = list.Remove(0)
	if err != nil || result != "G" {
		t.Error("Wrong bahavior of Remove function")
	}
	if list.Size() != 3 {
		t.Error("Wrong bahavior of Remove function")
	}

	result, err = list.Remove(2)
	if err != nil || result != "J" {
		t.Error("Wrong bahavior of Remove function")
	}
	if list.Size() != 2 {
		t.Error("Wrong bahavior of Remove function")
	}

	list.Remove(0)
	list.Remove(0)
	if list.Size() != 0 {
		t.Error("Wrong bahavior of Remove function")
	}
}

func TestCreation(t *testing.T) {
	list := NewEmptySyncList()

	if list.Size() != 0 {
		t.Error("Wrong behavior of Creation")
	}

	list = NewSyncList("A", "B", "C", "D", "E")
	if list.Size() != 5 {
		t.Error("Wrong behavior of Creation")
	}
}

func BenchmarkEmptyCreation(b *testing.B) {
	for n := 0; n < b.N; n++ {
		NewEmptySyncList()
	}
}

func BenchmarkCreation(b *testing.B) {
	for n := 0; n < b.N; n++ {
		NewSyncList("A", "B", "C", "D", "E")
	}
}
