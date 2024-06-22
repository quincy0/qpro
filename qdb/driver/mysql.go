package driver

import (
	"bytes"
	"strconv"

	"github.com/quincy0/qpro/qConfig"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Mysql struct {
}

func (e *Mysql) Open(conn string) (db *gorm.DB, err error) {
	return gorm.Open(mysql.Open(conn), &gorm.Config{})
}

func (e *Mysql) GetConnect() string {

	var conn bytes.Buffer
	conn.WriteString(qConfig.Settings.Database.Username)
	conn.WriteString(":")
	conn.WriteString(qConfig.Settings.Database.Password)
	conn.WriteString("@tcp(")
	conn.WriteString(qConfig.Settings.Database.Host)
	conn.WriteString(":")
	conn.WriteString(strconv.Itoa(qConfig.Settings.Database.Port))
	conn.WriteString(")")
	conn.WriteString("/")
	conn.WriteString(qConfig.Settings.Database.Name)
	conn.WriteString("?charset=utf8&parseTime=True&loc=Local&timeout=10000ms")
	return conn.String()
}
