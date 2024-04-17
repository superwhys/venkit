package vgorm

import (
	"testing"

	"github.com/superwhys/venkit/lg"
)

type UserModel struct {
	ID   uint `gorm:"primarykey"`
	Name string
	Age  int
}

func (um *UserModel) TableName() string {
	return "user"
}

func TestDialMysqlDB(t *testing.T) {
	mysqlConf := MysqlConfig{
		Instance: "localhost:3306",
		Database: "sql_test",
		Username: "root",
		Password: "yang4869",
	}

	RegisterSqlModel(&mysqlConf, &UserModel{})
	var resp []*UserModel
	if err := GetDbByModel(&UserModel{}).Find(&resp).Error; err != nil {
		lg.Errorf("get user data error: %v", err)
		return
	}

	lg.Info(lg.Jsonify(resp))
}
