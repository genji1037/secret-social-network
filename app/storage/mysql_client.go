package storage

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	"reflect"
	"secret-social-network/app/config"
)

// 数据库实例
var gormDb *gorm.DB

// 数据表接口
type Table interface {
	Create() error
}

// 数据库表集合
var tables = []Table{
	new(ConsensusOrder),
}

// 排序方式
type OrderBy string

const (
	OrderByAsc  OrderBy = "asc"
	OrderByDesc OrderBy = "desc"
)

// 打开数据库
func Open(conn Connection) error {
	var err error
	args := fmt.Sprintf("%s:%s@tcp(%s)/?charset=%s&parseTime=True&loc=%s",
		conn.User, conn.Password, conn.Host, conn.Charset, "Asia%2FShanghai")
	gormDb, err = gorm.Open("mysql", args)
	if err != nil {
		return err
	}

	var result []struct{}
	err = gormDb.Raw(fmt.Sprintf("SHOW DATABASES LIKE '%s'", conn.Database)).Scan(&result).Error
	if err != nil {
		return err
	}

	if len(result) == 0 {
		gormDb.Exec(fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARACTER SET %s;", conn.Database, conn.Charset))
	}

	args = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local",
		conn.User, conn.Password, conn.Host, conn.Database, conn.Charset)
	gormDb, err = gorm.Open("mysql", args)
	if err != nil {
		return err
	}

	gormDb.DB().SetMaxIdleConns(conn.MaxIdleConns)
	gormDb.DB().SetMaxOpenConns(conn.MaxOpenConns)

	for i := 0; i < len(tables); i++ {
		ensureTable(tables[i])
	}
	return nil
}

// 开始事务
func Transaction() (tx *gorm.DB, err error) {
	tx = gormDb.Begin()
	if err = tx.Error; err != nil {
		return nil, err
	}
	return tx, nil
}

// 确保表存在
func ensureTable(table interface{}) {
	typ := reflect.TypeOf(table)
	tablename := typ.String()[1:]
	if !gormDb.HasTable(table) {
		log.Infof("Creating table: %s", tablename)
		if err := gormDb.CreateTable(table).Error; err != nil {
			log.Warnf("Failed to create table %s, %v", tablename, err)
		}
	} else {
		serverCfg := config.GetServe()
		if serverCfg.MySQL.AlwayMigrate {
			log.Infof("Auto migrate table: %s", tablename)
			if err := gormDb.AutoMigrate(table).Error; err != nil {
				log.Warnf("Failed to auto migrate table %s, %v", tablename, err)
			}
		}
	}
}

func TxBegin() *gorm.DB {
	return gormDb.Begin()
}
