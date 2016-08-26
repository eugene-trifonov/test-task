package cache

import (
	"TestProject/utils"
)

const (
	MAX_LINE_NUMBER_FOR_COMMAND = 20
)

func NewUserFriendlyCommands(cache Cache, logger utils.Logger) CacheCommands {
	cmds := new(UserFriendlyCacheCommands)
	cmds.c = cache
	cmds.log = logger
	return cmds
}

func (this *UserFriendlyCacheCommands) GetValue(params []string) (interface{}, error) {
	value, err := this.BaseCacheCommands.GetValue(params)
	if err != nil {
		this.log.Logln("Cannot get value.", err)
	} else {
		if value == nil {
			this.log.Logf("No value for the key [%v]", params[0])
		} else {
			this.log.Log(value)
		}
	}
	return value, err
}

func (this *UserFriendlyCacheCommands) SetValue(params []string) (interface{}, error) {
	value, err := this.BaseCacheCommands.SetValue(params)
	if err != nil {
		this.log.Logln("Cannot set value.", err)
	} else {
		if value == nil {
			this.log.Logf("New value [%v] was set for a key [%v]", params[1], params[0])
		} else {
			this.log.Logf("The value [%v] for the key [%v] was replaced with a new one [%v]", value, params[0], params[1])
		}
	}
	return value, err
}

func (this *UserFriendlyCacheCommands) UpdateValue(params []string) (bool, error) {
	updated, err := this.BaseCacheCommands.UpdateValue(params)
	if err != nil {
		this.log.Logln("Cannot update value.", err)
	} else {
		if updated {
			this.log.Logf("The value [%v] for the key [%v] was updated with passed value [%v]", params[1], params[0], params[2])
		} else {
			this.log.Logf("Cannot update value [%v] for the key [%v], possibly the cached value was updated already", params[1], params[0])
		}
	}
	return updated, err
}

func (this *UserFriendlyCacheCommands) RemoveValue(params []string) (interface{}, bool, error) {
	oldValue, removed, err := this.BaseCacheCommands.RemoveValue(params)
	if err != nil {
		this.log.Logln("Cannot remove value.", err)
	} else {
		if removed {
			this.log.Logf("The value [%v] was deleted for the key [%v]", oldValue, params[0])
		} else if oldValue != nil {
			this.log.Logf("The value [%v] cannot be removed for the key [%v]. Probably the value was already updated.", oldValue, params[0])
		} else {
			this.log.Logf("There is no value for the key [%v]", params[0])
		}
	}
	return oldValue, removed, err
}

func (this *UserFriendlyCacheCommands) GetKeys(params []string) ([]string, error) {
	keys, err := this.BaseCacheCommands.GetKeys(params)
	if err != nil {
		this.log.Logln("Cannot get keys.", err)
	} else {
		displayKeys := keys[:utils.Min(MAX_LINE_NUMBER_FOR_COMMAND, len(keys))]
		this.log.Logln(utils.ConvertStrings(displayKeys)...)
		this.log.Logf("Count: [%d/%d]", len(displayKeys), len(keys))
	}
	return keys, err
}

func (this *UserFriendlyCacheCommands) GetListValue(params []string) (interface{}, error) {
	value, err := this.BaseCacheCommands.GetListValue(params)
	if err != nil {
		this.log.Log("Cannot get list value.", err)
	} else {
		this.log.Log(value)
	}
	return value, err
}

func (this *UserFriendlyCacheCommands) AppendListValue(params []string) error {
	err := this.BaseCacheCommands.AppendListValue(params)
	if err != nil {
		this.log.Logln("Append operation has been failed.", err)
	} else {
		this.log.Logf("The value [%v] is appended.", params[1])
	}
	return err
}

func (this *UserFriendlyCacheCommands) DeleteListValue(params []string) (interface{}, error) {
	value, err := this.BaseCacheCommands.DeleteListValue(params)
	if err != nil {
		this.log.Logln("Cannot delete value.", err)
	} else {
		this.log.Logf("The value [%v] was deleted.", value)
	}
	return value, err
}

func (this *UserFriendlyCacheCommands) GetListSize(params []string) (int, error) {
	value, err := this.BaseCacheCommands.GetListSize(params)
	if err != nil {
		this.log.Logln("Cannot get size of a list.", err)
	} else {
		this.log.Log(value)
	}
	return value, err
}

func (this *UserFriendlyCacheCommands) GetDictValue(params []string) (interface{}, error) {
	value, err := this.BaseCacheCommands.GetDictValue(params)
	if err != nil {
		this.log.Logln("Cannot get dictionary value.", err)
	} else {
		this.log.Log(value)
	}
	return value, err
}

func (this *UserFriendlyCacheCommands) SetDictValue(params []string) (interface{}, error) {
	value, err := this.BaseCacheCommands.SetDictValue(params)
	if err != nil {
		this.log.Logln("Cannot set dictionary value.", err)
	} else if value == nil {
		this.log.Logf("The dictionary pair was added sucessfully.")
	} else {
		this.log.Logf("The value [%v] was replaced.", value)
	}
	return value, err
}

func (this *UserFriendlyCacheCommands) DeleteDictValue(params []string) (interface{}, error) {
	value, err := this.BaseCacheCommands.DeleteDictValue(params)
	if err != nil {
		this.log.Logln("Cannot delete dictionary value.", err)
	} else if value == nil {
		this.log.Log("Dictionary does not contain a value for this key")
	} else {
		this.log.Logf("The value [%v] was deleted.", value)
	}
	return value, err
}

func (this *UserFriendlyCacheCommands) AppendDictValue(params []string) (bool, error) {
	appended, err := this.BaseCacheCommands.AppendDictValue(params)
	if err != nil {
		this.log.Logln("Cannot append dictionary value.", err)
	} else if appended {
		this.log.Log("The value was appended sucessfully.")
	} else {
		this.log.Log("Cannot append dictionary value, potentially the value exists.")
	}
	return appended, err
}

func (this *UserFriendlyCacheCommands) GetDictSize(params []string) (int, error) {
	size, err := this.BaseCacheCommands.GetDictSize(params)
	if err != nil {
		this.log.Logln("Cannot get size of the dictionary.", err)
	} else {
		this.log.Log(size)
	}
	return size, err
}

func (this *UserFriendlyCacheCommands) UpdateTTL(params []string) (bool, error) {
	updated, err := this.BaseCacheCommands.UpdateTTL(params)
	if err != nil {
		this.log.Logln("Cannot update ttl.", err)
	} else if !updated {
		this.log.Log("Ttl was not updated, potentially there is no value in cache anymore.")
	} else {
		this.log.Log("Ttl was updated sucessfully.")
	}
	return updated, err
}

func (this *UserFriendlyCacheCommands) GetSize() int {
	size := this.BaseCacheCommands.GetSize()
	this.log.Log(size)
	return size
}
