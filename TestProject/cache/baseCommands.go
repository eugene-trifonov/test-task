package cache

import (
	"TestProject/utils"
	"errors"
	"fmt"
	"strconv"
	"time"
)

const (
	TIMEOUT     = 60000000000 //1 minute
	DEFAULT_TTL = -1
)

type BaseCacheCommands struct {
	c Cache
}

func BaseCommands(cache Cache) CacheCommands {
	cmds := new(BaseCacheCommands)
	cmds.c = cache
	return cmds
}

func (this *BaseCacheCommands) RemoveValue(params []string) (interface{}, bool, error) {
	length := len(params)
	if length < 1 || length > 2 {
		return nil, false, errors.New("Wrong params count")
	}

	if length == 1 {
		oldValue := this.c.Remove(params[0])
		return oldValue, oldValue != nil, nil
	} else {
		return params[1], this.c.RemovePair(params[0], params[1]), nil
	}
}

func (this *BaseCacheCommands) UpdateValue(params []string) (bool, error) {
	length := len(params)
	if length < 3 || length > 4 {
		return false, errors.New("Wrong params count")
	}

	ttl := DEFAULT_TTL
	if length == 4 {
		value, err := strconv.Atoi(params[3])
		if err != nil {
			return false, err
		} else {
			ttl = value * int(time.Second)
		}
	}

	return this.c.ReplaceValueExpirable(params[0], params[1], params[2], int64(ttl)), nil
}

func (this *BaseCacheCommands) GetKeys(params []string) ([]string, error) {
	length := len(params)
	if length > 2 {
		return nil, errors.New("Wrong params count")
	}

	keys := this.c.GetKeys()

	if length == 0 {
		return keys, nil
	}

	startIndex, err := strconv.Atoi(params[0])
	if err != nil {
		return nil, err
	}

	if startIndex > length {
		return nil, errors.New(fmt.Sprintf("Too big integer [%v]", params[0]))
	}

	if length == 1 {
		return keys[startIndex:], nil
	}

	stopIndex, err := strconv.Atoi(params[1])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Invalid integer [%v]", params[1]))
	}

	if startIndex > length {
		return nil, errors.New(fmt.Sprintf("Too big integer [%v]", params[1]))
	}

	return keys[startIndex:stopIndex], nil
}

func (this *BaseCacheCommands) GetValue(params []string) (interface{}, error) {
	if len(params) != 1 {
		return nil, errors.New("Wrong params count")
	}

	return this.c.Get(params[0]), nil
}

func (this *BaseCacheCommands) SetValue(params []string) (interface{}, error) {
	length := len(params)
	if length < 2 || length > 3 {
		return nil, errors.New("Wrong params count")
	}

	ttl := DEFAULT_TTL
	if length == 3 {
		value, err := strconv.Atoi(params[2])
		if err != nil {
			return nil, errors.New("Invalid \"ttl\" value")
		} else {
			ttl = value * int(time.Second)
		}
	}

	return this.c.PutExpirable(params[0], params[1], int64(ttl)), nil
}

func (this *BaseCacheCommands) GetListValue(params []string) (interface{}, error) {
	if len(params) != 2 {
		return nil, errors.New("Wrong params count")
	}

	index, err := strconv.Atoi(params[1])
	if err != nil {
		return nil, err
	}

	value := this.c.Get(params[0])
	if value == nil {
		return nil, nil
	}

	listValue, ok := value.(utils.List)
	if !ok {
		return nil, errors.New(fmt.Sprintf("The value for the key [%v] is not a list.", params[0]))
	}

	value, err = listValue.Get(index)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (this *BaseCacheCommands) AppendListValue(params []string) error {
	length := len(params)
	if length < 2 || length > 3 {
		return errors.New("Wrong params count")
	}

	ttl := DEFAULT_TTL
	if length == 3 {
		ttlValue, err := strconv.Atoi(params[2])
		if err != nil {
			return errors.New("Wrong \"ttl\" value")
		}
		ttl = ttlValue * int(time.Second)
	}

	for {
		value := this.c.Get(params[0])
		if value == nil {
			list := utils.NewSyncList(params[1])
			existingValue := this.c.PutExpirableIfAbsent(params[0], list, int64(ttl))
			if existingValue == nil {
				break
			}
		} else {
			list, ok := value.(utils.List)
			if !ok {
				return errors.New(fmt.Sprintf("The value for the key [%v] is not a list.", params[0]))
			}
			list.Append(params[1])
			break
		}
	}

	return nil
}

func (this *BaseCacheCommands) DeleteListValue(params []string) (interface{}, error) {
	if len(params) != 2 {
		return nil, errors.New("Wrong params count")
	}

	index, err := strconv.Atoi(params[1])
	if err != nil {
		return nil, errors.New("Invalid \"index\" parameter")
	}

	for {
		value := this.c.Get(params[0])
		if value == nil {
			return nil, errors.New("Index out of bounds")
		}

		list, ok := value.(utils.List)
		if !ok {
			return nil, errors.New(fmt.Sprintf("The value for the key [%v] is not a list.", params[0]))
		}

		value, err := list.Remove(index)
		if err != nil {
			return nil, err
		}

		return value, nil
	}
}

func (this *BaseCacheCommands) GetListSize(params []string) (int, error) {
	if len(params) != 1 {
		return -1, errors.New("Wrong params count")
	}

	value := this.c.Get(params[0])
	if value == nil {
		return 0, nil
	}

	list, ok := value.(utils.List)
	if !ok {
		return -1, errors.New(fmt.Sprintf("The value for the key [%v] is not a list.", params[0]))
	}

	return list.Size(), nil
}

func (this *BaseCacheCommands) GetDictValue(params []string) (interface{}, error) {
	if len(params) != 2 {
		return nil, errors.New("Wrong params count")
	}

	value := this.c.Get(params[0])
	if value == nil {
		return nil, nil
	}

	dict, ok := value.(utils.Dict)
	if !ok {
		return nil, errors.New(fmt.Sprintf("The value for the key [%v] is not a dictionary.", params[0]))
	}

	return dict.Get(params[1]), nil
}

func (this *BaseCacheCommands) SetDictValue(params []string) (interface{}, error) {
	if len(params) != 3 {
		return nil, errors.New("Wrong params count")
	}

	for {
		value := this.c.Get(params[0])
		if value == nil {
			dict := utils.NewDict()
			dict.Put(params[1], params[2])
			existingValue := this.c.PutIfAbsent(params[0], dict)
			if existingValue == nil {
				return nil, nil
			}
		} else {
			dict, ok := value.(utils.Dict)
			if !ok {
				return nil, errors.New(fmt.Sprintf("The value for the key [%v] is not a dictionary.", params[0]))
			}
			return dict.Put(params[1], params[2]), nil
		}
	}
}

func (this *BaseCacheCommands) AppendDictValue(params []string) (bool, error) {
	if len(params) != 3 {
		return false, errors.New("Wrong params count")
	}

	for {
		value := this.c.Get(params[0])
		if value == nil {
			dict := utils.NewDict()
			dict.Put(params[1], params[2])
			existingValue := this.c.PutIfAbsent(params[0], dict)
			if existingValue == nil {
				return true, nil
			}
		} else {
			dict, ok := value.(utils.Dict)
			if !ok {
				return false, errors.New(fmt.Sprintf("The value for the key [%v] is not a dictionary.", params[0]))
			}
			existingValue := dict.PutIfAbsent(params[1], params[2])
			if existingValue == nil {
				return true, nil
			} else {
				return false, nil
			}
		}
	}
}

func (this *BaseCacheCommands) DeleteDictValue(params []string) (interface{}, error) {
	if len(params) != 2 {
		return nil, errors.New("Wrong params count")
	}

	value := this.c.Get(params[0])
	if value == nil {
		return nil, nil
	} else {
		dict, ok := value.(utils.Dict)
		if !ok {
			return nil, errors.New(fmt.Sprintf("The value for the key [%v] is not a dictionary.", params[0]))
		}

		return dict.Remove(params[1]), nil
	}
}

func (this *BaseCacheCommands) GetDictSize(params []string) (int, error) {
	if len(params) != 1 {
		return 0, errors.New("Wrong params count")
	}

	value := this.c.Get(params[0])
	if value == nil {
		return 0, nil
	} else {
		dict, ok := value.(utils.Dict)
		if !ok {
			return 0, errors.New(fmt.Sprintf("The value for the key [%v] is not a dictionary.", params[0]))
		}

		return dict.Size(), nil
	}
}

func (this *BaseCacheCommands) UpdateTTL(params []string) (bool, error) {
	if len(params) != 2 {
		return false, errors.New("Wrong params count")
	}

	ttl, err := strconv.Atoi(params[1])
	if err != nil {
		return false, errors.New(fmt.Sprintf("Invalid integer [%v]", params[1]))
	}

	return this.c.UpdateTTL(params[0], int64(ttl*int(time.Second))), nil
}

func (this *BaseCacheCommands) GetSize() int {
	return this.c.Size()
}

type UserFriendlyCacheCommands struct {
	BaseCacheCommands
	log utils.Logger
}
