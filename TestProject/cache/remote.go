package cache

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
)

var MAGIC_CMD = []byte{0, 0, 0}
var buffer [1024]byte

//Represents an API for remote Cache instance usage.
//Will be useful on any client side.
type RemoteCache interface {

	//Returns a value for a passed <key> that remove Cache contains, or <nil> in case there is no such <key> in a Cache.
	//Can return error, e.g. when a remove Cache is unavailable.
	Get(key string) (interface{}, error)

	//Puts a key-value pair into a remote Cache.
	//Returns a replaced value from a Cache, or <nil> if Cache doesn't contain any value for passed <key>.
	//Can return error, e.g. when a remove Cache is unavailable.
	Put(key string, value interface{}) (interface{}, error)

	//Provides an ability to put a key-value pair into a remote Cache with a specific time to live.
	//Returns a replaced value from a Cache, or <nil> if Cache doesn't contain any value for passed <key>.
	//Can return error, e.g. when a remove Cache is unavailable.
	PutExpirable(key string, value interface{}, ttl int64) (interface{}, error)

	//Removes a key-value pair from a remote Cache by provided <key>.
	//Returns a removed value from a Cache, or <nil> if Cache doesn't contain any value for passed <key>.
	//Can return error, e.g. when a remove Cache is unavailable.
	Remove(key string) (interface{}, error)

	//Removes a key-value pair from a remote Cache by provided <key> and <value,
	//the removal happens only in case both <key> and <value> matches an existing key-value pair in a remote Cache.
	//Returns <true> if the removal was successful, otherwise returns <false>.
	//Can return error, e.g. when a remove Cache is unavailable.
	RemovePair(key string, value interface{}) (bool, error)

	//Replaces value in a remote Cache if it exists.
	//Returns <true> if the replacement was successful, otherwise returns <false>.
	//Can return error, e.g. when a remove Cache is unavailable.
	ReplaceValue(key string, oldValue, newValue interface{}) (bool, error)

	//Replaces value in a remote Cache if it exists with a specific time to live.
	//Returns <true> if the replacement was successful, otherwise returns <false>.
	//Can return error, e.g. when a remove Cache is unavailable.
	ReplaceValueExpirable(key string, oldValue, newValue interface{}, ttl int64) (bool, error)

	//Updates time to live of any value stored in a remote Cache by <key>.
	//Returns <true> if update was successful, otherwise returns <false>.
	//Can return error, e.g. when a remove Cache is unavailable.
	UpdateTTL(key string, ttl int64) (bool, error)

	//Returns size of a remote Cache, the nmber of its key-value pairs at current time.
	//Can return error, e.g. when a remove Cache is unavailable.
	Size() (int, error)
}

type BaseRemoteCache struct {
	conn   net.Conn
	reader *bufio.Reader
}

func NewRemoteCache(conn net.Conn) RemoteCache {
	cache := new(BaseRemoteCache)
	cache.conn = conn
	cache.reader = bufio.NewReader(conn)
	return cache
}

func (this *BaseRemoteCache) Get(key string) (interface{}, error) {
	return this.execCmd(assembleCmd("get", key))
}

func (this *BaseRemoteCache) Put(key string, value interface{}) (interface{}, error) {
	return this.execCmd(assembleCmd("set", key, value))
}

func (this *BaseRemoteCache) PutExpirable(key string, value interface{}, ttl int64) (interface{}, error) {
	return this.execCmd(assembleCmd("set", key, value, ttl))
}

func (this *BaseRemoteCache) Remove(key string) (interface{}, error) {
	return this.execCmd(assembleCmd("delete", key))
}

func (this *BaseRemoteCache) RemovePair(key string, value interface{}) (bool, error) {
	result, err := this.execCmd(assembleCmd("delete", key, value))
	if err != nil {
		return false, err
	}

	removed, ok := result.(bool)
	if !ok {
		return false, errors.New("Unexpected non-bool value")
	}
	return removed, err
}

func (this *BaseRemoteCache) ReplaceValue(key string, oldValue, newValue interface{}) (bool, error) {
	result, err := this.execCmd(assembleCmd("update", key, oldValue, newValue))
	if err != nil {
		return false, err
	}

	updated, ok := result.(bool)
	if !ok {
		return false, errors.New("Unexpected non-bool value")
	}
	return updated, err
}

func (this *BaseRemoteCache) ReplaceValueExpirable(key string, oldValue, newValue interface{}, ttl int64) (bool, error) {
	result, err := this.execCmd(assembleCmd("update", key, oldValue, newValue, ttl))
	if err != nil {
		return false, err
	}

	updated, ok := result.(bool)
	if !ok {
		return false, errors.New("Unexpected non-bool value")
	}
	return updated, err
}

func (this *BaseRemoteCache) UpdateTTL(key string, ttl int64) (bool, error) {
	result, err := this.execCmd(assembleCmd("ttl", key, ttl))
	if err != nil {
		return false, err
	}

	updated, ok := result.(bool)
	if !ok {
		return false, errors.New("Unexpected non-bool value")
	}
	return updated, err
}

func (this *BaseRemoteCache) Size() (int, error) {
	result, err := this.execCmd("size")
	if err != nil {
		return -1, err
	}

	size, ok := result.(int)
	if !ok {
		return -1, errors.New("Unexpected non-bool value")
	}
	return size, err
}

func assembleCmd(values ...interface{}) string {
	strValues := make([]string, len(values))
	for i := 0; i < len(values); i++ {
		strValues[i] = fmt.Sprint(values[i])
	}
	return strings.Join(strValues, " ")
}

func createFirstCmd(cmd string) []byte {
	command := []byte{}
	command = append(command, MAGIC_CMD...)
	command = append(command, []byte(cmd)...)
	command = append(command, '\n')
	return command
}

func (this *BaseRemoteCache) execCmd(cmd string) (interface{}, error) {

	this.conn.Write([]byte(cmd + "\n"))

	result, _, err := this.reader.ReadLine()
	if err != nil {
		return nil, err
	}
	response := new(JsonResponse)
	err = json.Unmarshal(result, response)
	if err != nil {
		return nil, err
	}
	return response.Value, nil
}
