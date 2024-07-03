package main

import (
	"context"

	"github.com/superwhys/venkit/vgorm/v2"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name string
}

func (u *User) TableName() string {
	return "users"
}

func main() {
	conf := &vgorm.MysqlConfig{
		Instance: "localhost:3306",
		Database: "test",
		Username: "root",
		Password: "password",
	}

	_ = vgorm.RegisterSqlModelWithConf(conf, &User{})

	db := vgorm.GetDbByConf(conf)
	// vgorm.GetDbByModel(&User{})

	db.WithContext(context.Background())
}
