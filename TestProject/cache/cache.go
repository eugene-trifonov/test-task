package cache

import (
	"TestProject/utils"
	"sync"
	"time"
)

//Represents an interface for a Cache
type Cache interface {
	utils.SyncMap

	//Provides an ability to put a key-value pair with a specific time to live.
	//Returns a value that was replaced during the putting a new key-value pair.
	PutExpirable(key string, value interface{}, ttl int64) interface{}

	//Provides an ability to put a key-value pair with a specific time to live in case there is no value in a Cache yet.
	//Returns <nil> if a key-value pair was put into a Cache, or a value of an existing key-value pair.
	PutExpirableIfAbsent(key string, value interface{}, ttl int64) interface{}

	//Provides an ability to replace an existing key-value pair with a new value with a specific time to live.
	//Returns <nil> in case a key-value pair was not put into a Cache, otherwise returns a value that was replaced
	ReplaceExpirable(key string, value interface{}, ttl int64) interface{}

	//Provides an ability to replace an existing key-value pair with a new value with a specific time to live.
	//Replacement happens in case an existig value is equal to a passed <oldValue>.
	//Returns <true> if the replacement was successful, otherwise returns <false>.
	ReplaceValueExpirable(key string, oldValue, newValue interface{}, ttl int64) bool

	//Provides an ability to update time to live for any key in a Cache.
	//Returns <true> in case of replacement was successful, otherwise returns <false>.
	UpdateTTL(key string, ttl int64) bool
}

type cacheValue struct {
	storedTime time.Time
	ttl        time.Duration
	value      interface{}
}

type syncMap struct {
	sync.RWMutex
	m map[string]*cacheValue
}

func NewCache() Cache {
	cache := new(syncMap)
	cache.m = make(map[string]*cacheValue)
	return cache
}

func (this *cacheValue) isExpired(currentTime time.Time) bool {
	if this.ttl <= 0 {
		return false
	}
	return currentTime.After(this.storedTime.Add(this.ttl))
}

func (cache *syncMap) Get(key string) interface{} {
	cache.Lock()
	defer cache.Unlock()
	value := cache.m[key]
	if value == nil {
		return nil
	}

	if value.isExpired(time.Now()) {
		cache.remove(key)
		return nil
	}

	return value.value
}

func (cache *syncMap) Put(key string, value interface{}) interface{} {
	return cache.PutExpirable(key, value, -1)
}

func (cache *syncMap) PutExpirable(key string, value interface{}, ttl int64) interface{} {
	cache.Lock()
	defer cache.Unlock()
	return cache.put(key, value, ttl)
}

func (cache *syncMap) put(key string, value interface{}, ttl int64) interface{} {
	currValue := cache.m[key]
	cache.m[key] = &cacheValue{time.Now(), time.Duration(ttl), value}
	if currValue == nil || currValue.isExpired(time.Now()) {
		return nil
	} else {
		return currValue.value
	}
}

func (cache *syncMap) Remove(key string) interface{} {
	cache.Lock()
	defer cache.Unlock()
	value := cache.remove(key)
	if value == nil || value.isExpired(time.Now()) {
		return nil
	}

	return value.value
}

func (cache *syncMap) remove(key string) *cacheValue {
	oldValue := cache.m[key]
	delete(cache.m, key)
	return oldValue
}

func (cache *syncMap) PutIfAbsent(key string, value interface{}) interface{} {
	return cache.PutExpirableIfAbsent(key, value, -1)
}

func (cache *syncMap) PutExpirableIfAbsent(key string, value interface{}, ttl int64) interface{} {
	cache.Lock()
	defer cache.Unlock()
	t := time.Now()
	oldValue := cache.m[key]
	if oldValue == nil {
		cache.m[key] = &cacheValue{t, time.Duration(ttl), value}
		return nil
	}

	if oldValue.isExpired(t) {
		cache.remove(key)
		cache.m[key] = &cacheValue{t, time.Duration(ttl), value}
		return nil
	}

	return oldValue.value
}

func (cache *syncMap) GetKeys() []string {
	cache.Lock()
	defer cache.Unlock()

	keys := []string{}

	t := time.Now()
	for k := range cache.m {
		value := cache.m[k]
		if !value.isExpired(t) {
			keys = append(keys, k)
		}
	}

	return keys
}

func (cache *syncMap) RemovePair(key string, value interface{}) bool {
	cache.Lock()
	defer cache.Unlock()

	existingValue := cache.m[key]
	if existingValue == nil {
		return false
	}

	if existingValue.isExpired(time.Now()) {
		cache.remove(key)
		return false
	}

	if existingValue.value == value {
		cache.remove(key)
		return true
	}

	return false
}

func (cache *syncMap) Replace(key string, value interface{}) interface{} {
	return cache.ReplaceExpirable(key, value, -1)
}

func (cache *syncMap) ReplaceExpirable(key string, value interface{}, ttl int64) interface{} {
	cache.Lock()
	defer cache.Unlock()

	existingValue := cache.m[key]
	if existingValue == nil {
		return nil
	}

	if existingValue.isExpired(time.Now()) {
		cache.remove(key)
		return nil
	}

	cache.put(key, value, ttl)

	return existingValue.value
}

func (cache *syncMap) ReplaceValue(key string, oldValue, newValue interface{}) bool {
	cache.Lock()
	defer cache.Unlock()

	existingValue := cache.m[key]
	if existingValue == nil {
		return false
	}

	if existingValue.isExpired(time.Now()) {
		cache.remove(key)
		return false
	}

	if existingValue.value != oldValue {
		return false
	}

	cache.put(key, newValue, int64(existingValue.ttl))

	return true
}

func (cache *syncMap) ReplaceValueExpirable(key string, oldValue, newValue interface{}, ttl int64) bool {
	cache.Lock()
	defer cache.Unlock()

	existingValue := cache.m[key]
	if existingValue == nil {
		return false
	}

	if existingValue.isExpired(time.Now()) {
		cache.remove(key)
		return false
	}

	if existingValue.value != oldValue {
		return false
	}

	cache.put(key, newValue, ttl)

	return true
}

func (cache *syncMap) Size() int {
	cache.Lock()
	defer cache.Unlock()

	return len(cache.m)
}

func (cache *syncMap) UpdateTTL(key string, ttl int64) bool {
	cache.Lock()
	defer cache.Unlock()

	value := cache.m[key]
	if value == nil {
		return false
	}

	currentTime := time.Now()
	if value.isExpired(currentTime) {
		cache.remove(key)
		return false
	}

	value.ttl = time.Duration(ttl)
	if value.isExpired(currentTime) {
		cache.remove(key)
	}
	return true
}
