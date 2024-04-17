package vgorm

import (
	"fmt"
	"time"

	"github.com/superwhys/venkit/dialer"
	"gorm.io/gorm"
)

type SqliteConfig struct {
	DbFile string
}

func (s *SqliteConfig) GetDBType() dbType {
	return sqlite
}

func (s *SqliteConfig) GetService() string {
	return s.DbFile
}

func (s *SqliteConfig) GetUid() string {
	return fmt.Sprintf("sqlite-%v", s.DbFile)
}

func (s *SqliteConfig) DialGorm() (*gorm.DB, error) {
	logPrefix := fmt.Sprintf("sqlite:%s", s.DbFile)
	return dialer.DialSqlLiteGorm(
		s.DbFile,
		dialer.WithLogger(
			NewGormLogger(
				WithPrefix(logPrefix),
				WithSlowThreshold(time.Millisecond*200),
			),
		),
	)
}
