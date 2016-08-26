package cache

import (
	"TestProject/utils"
	"strconv"
	"testing"
)

func BenchmarkGetValueCommands(b *testing.B) {
	cache := NewCache()
	cmds := BaseCommands(cache)

	for i := 0; i < 1000; i++ {
		cache.Put(strconv.Itoa(i), i)
	}

	params := []string{"155"}
	for n := 0; n < b.N; n++ {
		cmds.GetValue(params)
	}

}

func BenchmarkSetValueCommands(b *testing.B) {
	cache := NewCache()
	cmds := BaseCommands(cache)

	params := []string{"A", "B"}
	for n := 0; n < b.N; n++ {
		params[0] = strconv.Itoa(n)
		params[1] = params[0]
		cmds.SetValue(params)
	}

}

func BenchmarkSetExpirapbleCommands(b *testing.B) {
	cache := NewCache()
	cmds := BaseCommands(cache)

	params := []string{"A", "B", "100000"}
	for n := 0; n < b.N; n++ {
		params[0] = strconv.Itoa(n)
		params[1] = params[0]
		cmds.SetValue(params)
	}

}

func BenchmarkAppendListValuesCommands(b *testing.B) {
	cache := NewCache()
	cmds := BaseCommands(cache)

	for i := 0; i < 1000; i++ {
		cache.Put(strconv.Itoa(i), i)
	}

	cache.Remove("133")

	params := []string{"133", "value"}
	for n := 0; n < b.N; n++ {
		params[1] = strconv.Itoa(n)
		cmds.AppendListValue(params)
	}

}

func BenchmarkGetListValuesCommands(b *testing.B) {
	cache := NewCache()
	cmds := BaseCommands(cache)

	for i := 0; i < 1000; i++ {
		cache.Put(strconv.Itoa(i), i)
	}

	list := utils.NewSyncList("1", "2", "3", "4", "5", "6", "7")

	cache.Put("133", list)

	params := []string{"133", "3"}
	for n := 0; n < b.N; n++ {
		cmds.GetListValue(params)
	}

}

func BenchmarkSetDictValuesCommands(b *testing.B) {
	cache := NewCache()
	cmds := BaseCommands(cache)

	for i := 0; i < 1000; i++ {
		cache.Put(strconv.Itoa(i), i)
	}

	dict := utils.NewDict()

	cache.Put("133", dict)

	params := []string{"133", "key", "value"}
	for n := 0; n < b.N; n++ {
		params[1] = strconv.Itoa(n)
		params[2] = params[1]
		cmds.SetDictValue(params)
	}

}

func BenchmarkGetDictValuesCommands(b *testing.B) {
	cache := NewCache()
	cmds := BaseCommands(cache)

	for i := 0; i < 1000; i++ {
		cache.Put(strconv.Itoa(i), i)
	}

	dict := utils.NewDict()
	for i := 0; i < 100; i++ {
		dict.Put(strconv.Itoa(i), i)
	}

	cache.Put("133", dict)

	params := []string{"133", "55"}
	for n := 0; n < b.N; n++ {
		cmds.GetDictValue(params)
	}

}
