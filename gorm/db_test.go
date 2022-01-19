package gorm

import (
	"testing"

	"gorm.io/gorm"
)

var (
	dbConf = Database{
		Host:     "127.0.0.1",
		Port:     "3306",
		DBName:   "test",
		UserName: "root",
		Passwd:   "abc123",
	}
)

func TestGetDBConn(t *testing.T) {
	type args struct {
		conf Database
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "TestGetDBConn1",
			args: args{
				conf: dbConf,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetDBConn(tt.args.conf)
		})
	}
}

func TestAutoMigrateTable(t *testing.T) {
	user := &User{}
	AutoMigrateTable(GetDBConn(dbConf), user)
}

func TestCreateOrUpdateTable(t *testing.T) {
	user := &User{
		Model: gorm.Model{
			ID: 1,
		},
		Ver: "1",
		Md5: "xxx",
		Url: "baidu.com",
	}
	CreateOrUpdateTable(GetDBConn(dbConf), user)
}
