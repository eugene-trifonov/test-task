package cache

import (
	"TestProject/utils"

	"net"
)

type JsonCacheCommands struct {
	BaseCacheCommands
	conn net.Conn
	log  utils.Logger
}

func NewJsonCacheCommands(cache Cache, conn net.Conn, log utils.Logger) CacheCommands {
	cmds := new(JsonCacheCommands)
	cmds.c = cache
	cmds.conn = conn
	cmds.log = log
	return cmds
}

func (this *JsonCacheCommands) GetValue(params []string) (interface{}, error) {
	value, err := this.BaseCacheCommands.GetValue(params)
	this.writeResponse(value, err)
	return value, err
}

func (this *JsonCacheCommands) SetValue(params []string) (interface{}, error) {
	value, err := this.BaseCacheCommands.SetValue(params)
	this.writeResponse(value, err)
	return value, err
}

func (this *JsonCacheCommands) UpdateValue(params []string) (bool, error) {
	value, err := this.BaseCacheCommands.UpdateValue(params)
	this.writeResponse(value, err)
	return value, err
}

func (this *JsonCacheCommands) RemoveValue(params []string) (interface{}, bool, error) {
	value, removed, err := this.BaseCacheCommands.RemoveValue(params)
	if removed {
		this.writeResponse(value, err)
	} else {
		this.writeResponse(nil, err)
	}
	return value, removed, err
}

func (this *JsonCacheCommands) GetListValue(params []string) (interface{}, error) {
	value, err := this.BaseCacheCommands.GetListValue(params)
	this.writeResponse(value, err)
	return value, err
}

func (this *JsonCacheCommands) AppendListValue(params []string) error {
	err := this.BaseCacheCommands.AppendListValue(params)
	this.writeResponse(nil, err)
	return err
}

func (this *JsonCacheCommands) DeleteListValue(params []string) (interface{}, error) {
	value, err := this.BaseCacheCommands.DeleteListValue(params)
	this.writeResponse(value, err)
	return value, err
}

func (this *JsonCacheCommands) GetListSize(params []string) (int, error) {
	size, err := this.BaseCacheCommands.GetListSize(params)
	this.writeResponse(size, err)
	return size, err
}

func (this *JsonCacheCommands) GetDictValue(params []string) (interface{}, error) {
	value, err := this.BaseCacheCommands.GetDictValue(params)
	this.writeResponse(value, err)
	return value, err
}

func (this *JsonCacheCommands) SetDictValue(params []string) (interface{}, error) {
	value, err := this.BaseCacheCommands.SetDictValue(params)
	this.writeResponse(value, err)
	return value, err
}

func (this *JsonCacheCommands) AppendDictValue(params []string) (bool, error) {
	appended, err := this.BaseCacheCommands.AppendDictValue(params)
	this.writeResponse(appended, err)
	return appended, err
}

func (this *JsonCacheCommands) DeleteDictValue(params []string) (interface{}, error) {
	value, err := this.BaseCacheCommands.DeleteDictValue(params)
	this.writeResponse(value, err)
	return value, err
}

func (this *JsonCacheCommands) GetDictSize(params []string) (int, error) {
	size, err := this.BaseCacheCommands.GetDictSize(params)
	this.writeResponse(size, err)
	return size, err
}

func (this *JsonCacheCommands) GetKeys(params []string) ([]string, error) {
	keys, err := this.BaseCacheCommands.GetKeys(params)
	this.writeResponse(keys, err)
	return keys, err
}

func (this *JsonCacheCommands) UpdateTTL(params []string) (bool, error) {
	updated, err := this.BaseCacheCommands.UpdateTTL(params)
	this.writeResponse(updated, err)
	return updated, err
}

func (this *JsonCacheCommands) GetSize() int {
	size := this.BaseCacheCommands.GetSize()
	this.writeResponse(size, nil)
	return size
}

func (this *JsonCacheCommands) writeResponse(value interface{}, err error) {
	if WriteResponse(this.conn, value, err) != nil {
		this.log.Log(err)
	}
}
