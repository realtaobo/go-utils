package gorm

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 连接数据库的基本配置
type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DBName   string `yaml:"db"`
	UserName string `yaml:"user"`
	Passwd   string `yaml:"passwd"`
}

// 获取数据库连接句柄
func GetDBConn(conf Database) *gorm.DB {
	db, err := OpenConnection(conf)
	if err != nil {
		panic(err)
	}
	return db
}

// 初始化mysql数据库链接
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
	return
}

// 执行原生sql语句
func ExecRawSQL(db *gorm.DB, sql string, result interface{}, value ...interface{}) error {
	if err := db.Raw(sql, value...).Scan(result).Error; err != nil {
		return err
	}
	return nil
}

// 插入数据。
// table 类型应该为 &struct{}。
func CreateTable(db *gorm.DB, table interface{}) error {
	if sqlErr := db.Create(table).Error; sqlErr != nil {
		return fmt.Errorf("db insert table failed, due to %v", sqlErr)
	}
	return nil
}

// 获取某个查询语句的查询结果总数，适用于只获取结果数的场景
func GetTableQueryTotal(db *gorm.DB, subquery string, value ...interface{}) (int64, error) {
	sql := "select count(*) from (" + subquery + ") T"
	var result int64
	if err := db.Raw(sql, value...).Scan(&result).Error; err != nil {
		return result, err
	}
	return result, nil
}

// 向表中批量插入数据
func CreateTableInBatches(db *gorm.DB, table interface{}, size int) error {
	if sqlErr := db.CreateInBatches(table, size).Error; sqlErr != nil {
		return fmt.Errorf("db insert table failed, due to %v", sqlErr)
	}
	return nil
}

// 表行不存在时插入行，存在时更新。
// table 类型应该为 &struct{}。
func CreateOrUpdateTable(db *gorm.DB, table interface{}) error {
	if sqlErr := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).CreateInBatches(table, 1000).Error; sqlErr != nil {
		return fmt.Errorf("db upsert table failed, due to %v", sqlErr)
	}
	return nil
}

// 带where条件的查询语句。
// result 为 *[]TableStruct结构
// limit = -1表示不进行分页查询
func GetTableBySpec(db *gorm.DB, limit, offset int, result, query interface{}, args ...interface{}) error {
	if err := db.Where(query, args...).Limit(limit).Offset(offset).Find(result).Error; err != nil {
		return err
	}
	return nil
}

// 不带where条件的查询语句。
// result 为 *[]TableStruct结构
// limit = -1表示不进行分页查询
func GetTable(db *gorm.DB, limit, offset int, result interface{}) error {
	if err := db.Limit(limit).Offset(offset).Find(result).Error; err != nil {
		return err
	}
	return nil
}

// 自动创建指定结构表。
// table 类型应该为 &struct{}。
func AutoMigrateTable(db *gorm.DB, table interface{}) error {
	if err := db.AutoMigrate(table); err != nil {
		return err
	}
	return nil
}

// 删除并重新创建表
func ReCreateTable(db *gorm.DB, table interface{}) error {
	// 若存在这样的表，则先删除
	if db.Migrator().HasTable(table) {
		if err := db.Migrator().DropTable(table); err != nil {
			return err
		}
	}
	return AutoMigrateTable(db, table)
}
