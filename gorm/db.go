package gorm

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var db *gorm.DB

// mysql config
type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DBName   string `yaml:"db"`
	UserName string `yaml:"user"`
	Passwd   string `yaml:"passwd"`
}

// 获取 DB 链接
func GetDBConn(conf Database) *gorm.DB {
	if db != nil {
		return db
	} else {
		db, err := InitConnection(conf)
		if err != nil {
			panic(err)
		}
		return db
	}
}

// InitConnection 初始化mysql数据库链接 -- 方法一
func InitConnection(conf Database) (*gorm.DB, error) {
	path := strings.Join([]string{conf.UserName, ":", conf.Passwd, "@tcp(", conf.Host, ":", conf.Port, ")/", conf.DBName, "?charset=utf8&parseTime=true"}, "")
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       path,  // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}))
	if err != nil {
		return nil, err
	}
	return db, nil
}

// OpenConnection 初始化mysql数据库链接 -- 方法二
func OpenConnection(dbConf Database) (db *gorm.DB, err error) {
	dbDSN := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConf.UserName,
		dbConf.Passwd,
		dbConf.Host,
		dbConf.Port,
		dbConf.DBName,
	)
	for i := 0; i <= 30; i++ {
		db, err = gorm.Open(mysql.Open(dbDSN), &gorm.Config{})
		if err == nil {
			break
		} else {
			log.Printf("Retry | failed to connect database, got error %v", err)
			time.Sleep(time.Second)
		}
	}
	log.Println("success to connect database", dbDSN[strings.Index(dbDSN, `@`):])
	return
}

// AutoMigrateTable 自动初始化表
// table 类型应该为 &struct{}
func AutoMigrateTable(db *gorm.DB, table interface{}) {
	var err error
	// 自动初始化数据表
	if err = db.AutoMigrate(table); err != nil {
		log.Fatalf("Failed to auto migrate ProjectMember, but got error %v", err)
	}
}

// CreateOrUpdateTable 不存在时插入，存在时更新即可
// table 类型应该为 &struct{}
func CreateOrUpdateTable(db *gorm.DB, table interface{}) error {
	if sqlErr := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(table).Error; sqlErr != nil {
		return fmt.Errorf("db upsert table failed, due to %v", sqlErr)
	}
	return nil
}

// CreateTable 插入数据
// table 类型应该为 &struct{}
func CreateTable(db *gorm.DB, table interface{}) error {
	if sqlErr := db.Create(table).Error; sqlErr != nil {
		return fmt.Errorf("db insert table failed, due to %v", sqlErr)
	}
	return nil
}
