package storage

import (
	"util"
	"fmt"
	"strings"
	"mylog"
)

type Storage interface {
	GetData() (map[string]interface{}, error)
}

type StorageCreater interface {
	GetTypes() []string
	// 获取相关Storage,必须支持空参数放回一个相关Storage实体
	GetStorage(config map[string]interface{}) (interface{}, error)
}

var storages = map[string]StorageCreater{}

func Register(sc StorageCreater) {
	storage, _ := sc.GetStorage(map[string]interface{}{})
	if storage == nil {
		panic(fmt.Errorf("config 参数为空时必须返回一个不为nil的Storage对象"))
	}
	if _, ok := storage.(Storage); !ok {
		panic(fmt.Errorf("返回的类型必须是Storage类型"))
	}
	for _, name := range sc.GetTypes() {
		name = strings.ToLower(name)
		if _, ok := storages[name]; ok {
			mylog.Warn("已经存在一个该type的sotrage creater, type: " + name)
		}
		storages[name] = sc
	}
}

// GetStorage 根据配置信息获取监控储存
func GetStorage(config map[string]interface{}) (Storage, error) {
	class := util.GetMapString(config, "_type", "default")
	storage, err := storages[class].GetStorage(config)
	if err != nil {
		return nil, err
	}
	return storage.(Storage), nil
}

func init() {
	Register(EsStorage{})
}
