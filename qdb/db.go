package qdb

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/quincy0/qpro/qConfig"
	"github.com/quincy0/qpro/qLog"
	"github.com/quincy0/qpro/qdb/driver"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

var (
	sqlDB *sql.DB
	Db    *gorm.DB
)

type Driver interface {
	Open(conn string) (db *gorm.DB, err error)

	GetConnect() string
}

func InitMysql() {
	var err error
	d := new(driver.Mysql)
	conn := d.GetConnect()
	Db, err = d.Open(conn)

	if err = Db.Use(tracing.NewPlugin()); err != nil {
		panic(err)
	}

	if err != nil {
		qLog.Fatal("db connect error", zap.Error(err))
	} else {
		qLog.Info("db connect success!")
	}

	if Db.Error != nil {
		qLog.Fatal("database error", zap.Error(Db.Error))
	}

	sqlDB, err = Db.DB()
	if err != nil {
		qLog.Fatal("db connect error", zap.Error(err))
	}

	// 设置最大打开连接数
	sqlDB.SetMaxOpenConns(qConfig.Settings.Database.MaxOpenConn)

	// 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用
	sqlDB.SetMaxIdleConns(qConfig.Settings.Database.MaxIdleConn)
}

func Close() error {
	return sqlDB.Close()
}

func InitData() error {
	filePath := "config/db.sql"
	contents, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	//因为contents是[]byte类型，直接转换成string类型后会多一行空格,需要使用strings.Replace替换换行符
	sql := strings.Replace(string(contents), "\n", "", 1)
	if err != nil {
		fmt.Println("数据库基础数据初始化脚本读取失败！原因:", err.Error())
		return err
	}
	sqlList := strings.Split(sql, ";")
	for _, sql := range sqlList {
		if strings.Contains(sql, "--") {
			fmt.Println(sql)
			continue
		}
		sqlValue := strings.Replace(sql+";", "\n", "", 1)
		if err = Db.Exec(sqlValue).Error; err != nil {
			if !strings.Contains(err.Error(), "Query was empty") {
				return err
			}
		}
	}
	return nil
}
