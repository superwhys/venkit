package vgorm

import (
	"fmt"
	"reflect"
	"sync"

	"gorm.io/gorm"
)

type SqlModel interface {
	TableName() string
}

type getClientFunc func() *client

var (
	dbInstanceClientFuncMap = make(map[string]getClientFunc)
	modelDbMap              = make(map[string]string)
)

func getMysqlDB(m SqlModel) *gorm.DB {
	rt := reflect.TypeOf(m)
	modelKey := fmt.Sprintf("%v-%v", rt.String(), m.TableName())

	key, exists := modelDbMap[modelKey]
	if !exists {
		panic(fmt.Sprintf("model %v has not been register", rt.String()))
	}

	clientFunc, ok := getInstanceClientFunc(key)
	if !ok {
		panic(fmt.Sprintf("db instance %v not found", key))
	}

	return clientFunc().DB()
}

func getInstanceClientFunc(key string) (getClientFunc, bool) {
	f, exists := dbInstanceClientFuncMap[key]
	return f, exists
}

func GetMysqlDByModel(m SqlModel) *gorm.DB {
	db := getMysqlDB(m).Model(m)
	return db
}

func registerInstance(conf Config) {
	key := conf.GetUid()
	dbInstanceClientFuncMap[key] = func() getClientFunc {
		var cli *client
		var once sync.Once

		f := func() *client {
			once.Do(func() {
				cli = NewClient(conf)
			})
			return cli
		}

		return f
	}()
}

func RegisterSqlModel(conf Config, ms ...SqlModel) {
	if _, exists := getInstanceClientFunc(conf.GetUid()); exists {
		panic(fmt.Sprintf("%v has been register", conf.GetUid()))
	}

	registerInstance(conf)

	for _, m := range ms {
		rt := reflect.TypeOf(m)
		modelKey := fmt.Sprintf("%v-%v", rt.String(), m.TableName())
		if key, exists := modelDbMap[modelKey]; exists {
			panic(fmt.Sprintf("model %v has been register into %v", modelKey, key))
		}
		modelDbMap[modelKey] = conf.GetUid()
	}
}
