package main

import (
	"context"

	"github.com/superwhys/venkit/vgorm/v2"
)

func main() {
	conf := &vgorm.MysqlConfig{
		Instance: "localhost:3306",
		Database: "test",
		Username: "root",
		Password: "password",
	}

	db, err := conf.DialGorm()
	if err != nil {
		panic(err)
	}

	db.WithContext(context.Background())
}
