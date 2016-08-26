package utils

import (
	"errors"
	"fmt"
	"sync"
)

//Presents a List interface.
type List interface {
	//Returns a value from a list by provided index, or error if the index is invalid.
	Get(i int) (interface{}, error)

	//Removes a value from a list by provided index.
	//Returns removed value, or error if the index is invalid.
	Remove(i int) (interface{}, error)

	//Appends a new value to a list.
	//Returns the same list for a comfortabl usage.
	Append(value interface{}) List

	//Returns a size of a list, the count of its values.
	Size() int

	//Allows a list to be presented as a string where necessary.
	String() string
}

//Synchronized implementation of a List interface
type SyncList struct {
	sync.RWMutex
	slice []interface{}
}

//Creates an empty synchronized list
func NewEmptySyncList() List {
	list := new(SyncList)
	list.slice = []interface{}{}
	return list
}

//Creates synchronized list with populated values
func NewSyncList(values ...interface{}) List {
	list := new(SyncList)
	list.slice = values[:]
	return list
}

//Returns a value by passed index.
//Or error, e.g. int case of out of bound
func (list *SyncList) Get(i int) (interface{}, error) {
	list.Lock()
	defer list.Unlock()

	if i < 0 || i >= len(list.slice) {
		return nil, errors.New("Index out of bound")
	}

	return list.slice[i], nil
}

//Removes a value at passed index.
//Or err, e.g. in case of out of bound
func (list *SyncList) Remove(i int) (interface{}, error) {
	list.Lock()
	defer list.Unlock()

	length := len(list.slice)
	if isInvalidIndex(i, length) {
		return nil, errors.New("Index out of bound")
	}

	value := list.slice[i]
	if i == 0 {
		list.slice = list.slice[1:]
	} else if i == length-1 {
		list.slice = list.slice[:length-1]
	} else {
		list.slice = append(list.slice[:i], list.slice[i+1:]...)
	}

	return value, nil
}

func (list *SyncList) isInvalidIndex(i int) bool {
	return isInvalidIndex(i, len(list.slice))
}

//Appends a passed value to the end of the list.
func (list *SyncList) Append(value interface{}) List {
	list.Lock()
	defer list.Unlock()

	list.slice = append(list.slice, value)
	return list
}

//Returns a size of the list.
func (list *SyncList) Size() int {
	list.Lock()
	defer list.Unlock()

	return len(list.slice)
}

//Returns string interpretation of the list.
func (list *SyncList) String() string {
	slice := list.copySlice()

	length := len(slice)
	if length == 0 {
		return "List []"
	}

	return fmt.Sprint("List ", slice)
}

func (list *SyncList) copySlice() []interface{} {
	list.Lock()
	defer list.Unlock()

	copiedSlice := make([]interface{}, len(list.slice))
	copy(copiedSlice, list.slice)

	return copiedSlice
}

func isInvalidIndex(i, length int) bool {
	return i < 0 || i >= length
}
