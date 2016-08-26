package utils

import (
	"fmt"
	"strings"
	"sync"
)

//Represents a dictionary interface.
type Dict interface {
	SyncMap
	//Represents a dictionary as a one string value where necessary.
	String() string
}

//Synchronized implementation of dictionary interface
type baseDict struct {
	sync.RWMutex
	m map[string]interface{}
}

//Creates a new instance of dictionary
func NewDict() Dict {
	dict := new(baseDict)
	dict.m = make(map[string]interface{})
	return dict
}

//Returns string interpretation of the dictionary.
func (dict *baseDict) String() string {
	list := dict.copyKeyValues()
	return fmt.Sprint("Map [", strings.Join(list, " | "), "]")
}

func (dict *baseDict) copyKeyValues() []string {
	dict.Lock()
	defer dict.Unlock()

	list := make([]string, len(dict.m))
	i := 0
	for key := range dict.m {
		list[i] = fmt.Sprint(key, "-", dict.m[key])
		i++
	}

	return list
}

//Returns a word(value) for a passed key.
func (dict *baseDict) Get(key string) interface{} {
	dict.Lock()
	defer dict.Unlock()

	return dict.m[key]
}

//Sets key-value pair into a dictionary.
func (dict *baseDict) Put(key string, value interface{}) interface{} {
	dict.Lock()
	defer dict.Unlock()

	oldValue := dict.m[key]
	dict.m[key] = value

	return oldValue
}

//Inserts a new key-value pair into a dictionary.
//If the dictionary contains any other value for a passed key,
//then value won't be inserted.
func (dict *baseDict) PutIfAbsent(key string, value interface{}) interface{} {
	dict.Lock()
	defer dict.Unlock()

	oldValue := dict.m[key]
	if oldValue != nil {
		return oldValue
	}

	dict.m[key] = value
	return nil
}

//Returns a set of dictionary keys
func (dict *baseDict) GetKeys() []string {
	dict.Lock()
	defer dict.Unlock()

	keys := make([]string, len(dict.m))
	i := 0
	for key := range dict.m {
		keys[i] = key
		i++
	}

	return keys
}

//Replaces existing value for a passed key.
//Returns <nil> when replacement was unsuccessful and previous value otherwise.
func (dict *baseDict) Replace(key string, value interface{}) interface{} {
	dict.Lock()
	defer dict.Unlock()

	oldValue := dict.m[key]
	if oldValue == nil {
		return nil
	}

	dict.m[key] = value
	return oldValue
}

//Replaces an exact key-value pair if it exists.
//Returns <true> if the replacement was successful, otherwise returns <false>.
func (dict *baseDict) ReplaceValue(key string, oldValue, newValue interface{}) bool {
	dict.Lock()
	defer dict.Unlock()

	if oldValue != dict.m[key] {
		return false
	}

	dict.m[key] = newValue
	return true
}

//Removes a key-value pairby passed key.
//Returns <nil> if Dictionary doesn't contain the key or contains <nil>-value, otherwise returns a removed value.
func (dict *baseDict) Remove(key string) interface{} {
	dict.Lock()
	defer dict.Unlock()

	oldValue := dict.m[key]
	delete(dict.m, key)

	return oldValue
}

//Removes an exact key-value pair if it exists.
//Returns <true> if removal was successful, and <false> otherwise.
func (dict *baseDict) RemovePair(key string, value interface{}) bool {
	dict.Lock()
	defer dict.Unlock()

	oldValue := dict.m[key]
	if oldValue != value {
		return false
	}

	delete(dict.m, key)
	return true
}

//Returns size of dictionary, count of key-value pair in it.
func (dict *baseDict) Size() int {
	dict.Lock()
	defer dict.Unlock()

	return len(dict.m)
}
