package vgorm

import (
	"fmt"
	"strings"
	"time"

	"github.com/superwhys/venkit/dialer"
	"gorm.io/gorm"
)

type MysqlConfig struct {
	Instance string
	Database string
	Username string
	Password string
}

func (m *MysqlConfig) GetDBType() dbType {
	return mysql
}

func (m *MysqlConfig) GetService() string {
	return m.Instance
}

func (m *MysqlConfig) GetUid() string {
	return fmt.Sprintf("mysql-%v-%v", m.Instance, m.Database)
}

func (m *MysqlConfig) DialGorm() (*gorm.DB, error) {
	m.TrimSpace()
	logPrefix := fmt.Sprintf("mysql:%s", m.Database)

	return dialer.DialMysqlGorm(
		m.Instance,
		dialer.WithAuth(m.Username, m.Password),
		dialer.WithDBName(m.Database),
		dialer.WithLogger(
			NewGormLogger(
				WithPrefix(logPrefix),
				WithSlowThreshold(time.Millisecond*200),
			),
		),
	)
}

func (m *MysqlConfig) TrimSpace() {
	m.Username = strings.TrimSpace(m.Username)
	m.Password = strings.TrimSpace(m.Password)
	m.Instance = strings.TrimSpace(m.Instance)
	m.Database = strings.TrimSpace(m.Database)
}
