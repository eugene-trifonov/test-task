package utils

//Represents an interface for any synchronized map.
type SyncMap interface {
	//Provides an array of keys that are currently used in a map.
	GetKeys() []string

	//Provides a value that is stored against a passed <key>.
	//<nil> can be returned in case of key-value pair is not in a map or value is <nil>.
	Get(key string) interface{}

	//Provides an ability to set up a key-value pair into a map.
	//Returns a value that was replaced, if any, otherwise returns <nil>.
	Put(key string, value interface{}) interface{}

	//Puts key-value pair into a map in case of a map doesn't contain another value for a passed <key>.
	//Returns <nil> in case the key-value pair was put into a map,
	//otherwise returns an existing value for a passed <key>
	PutIfAbsent(key string, value interface{}) interface{}

	//Replaces an existing value in a map with a passed one.
	//Returns replaced value for a passed <key>, or <nil> if a map doesn't contain any value against a passed <key>.
	Replace(key string, value interface{}) interface{}

	//Replaces an existing value in a map, in case it is the same as passed <oldValue>.
	//Returns <true> in case the replacement was successful, otherwise returns <false>.
	ReplaceValue(key string, oldValue, newValue interface{}) bool

	//Removes a key-value pair from a map by a passed <key>.
	//Returns a value that was stored against a passed <key>, or <nil> if there was no such <key>.
	Remove(key string) interface{}

	//Removes a key-value pair from a map in case both <key> and <value> are the same as a map contains currently.
	//Returns <true> in case of removing was successful, otherwise return <false>.
	RemovePair(key string, value interface{}) bool

	//Returns a size of a map, the count of key-value pairs.
	Size() int
}
