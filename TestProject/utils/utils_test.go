package utils

import (
	"testing"
)

func TestMinFunction(t *testing.T) {
	result := Min(1, 2)
	if result != 1 {
		t.Error("Wrong behavior of Min function")
	}

	result = Min(2, 1)
	if result != 1 {
		t.Error("Wrong behavior of Min function")
	}

	result = Min(1, -1)
	if result != -1 {
		t.Error("Wrong behavior of Min function")
	}

	result = Min(1, 1)
	if result != 1 {
		t.Error("Wrong behavior of Min function")
	}
}

func BenchmarkMinFunction(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Min(-100, 100)
	}
}

func TestCheckPortFunction(t *testing.T) {
	err := CheckPort("A123")
	if err == nil {
		t.Error("Wrong behavior of CheckPort function")
	}

	err = CheckPort("2343")
	if err != nil {
		t.Error("Wrong behavior of CheckPort function")
	}
}

func BenchmarkIncorrectCheckPort(b *testing.B) {
	benchmarkCheckPort("123A", b)
}

func BenchmarkCorrectCheckPort(b *testing.B) {
	benchmarkCheckPort("1123", b)
}

func benchmarkCheckPort(port string, b *testing.B) {
	for n := 0; n < b.N; n++ {
		CheckPort(port)
	}
}
