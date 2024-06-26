package vgorm

import (
	"fmt"
	
	"github.com/superwhys/venkit/lg/v2"
	"gorm.io/gorm"
)

type dbType int

const (
	mysql = iota
	sqlite
)

type Config interface {
	GetDBType() dbType
	GetUid() string
	GetService() string
	DialGorm() (*gorm.DB, error)
}

type client struct {
	db     *gorm.DB
	config Config
}

func NewClient(conf Config) *client {
	if conf.GetService() == "" {
		panic(fmt.Sprintf("vgorm: %v db service name can not be empty", conf.GetDBType()))
	}
	
	c := &client{config: conf}
	c.dial()
	return c
}

func (c *client) dial() {
	db, err := c.config.DialGorm()
	lg.PanicError(err, "mqlClient: new client error")
	
	c.db = db
}

func (c *client) DB() *gorm.DB {
	return c.db
}
