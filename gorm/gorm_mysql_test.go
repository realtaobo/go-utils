package gorm

import (
	"log"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
var dsn string = "test:test1998@tcp(127.0.0.1:3306)/test"

type User struct {
	// gorm.Model
	Ver string `gorm:"column:ver"`
	Md5 string `gorm:"column:md5;primaryKey"`
	Url string `gorm:"column:url"`
}

func (*User) TableName() string {
	return "user"
}

// 测试自动建表
func TestAutoMigrate(t *testing.T) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	log.Println(err, "......")
	if err = db.AutoMigrate(&User{}); err != nil {
		log.Println(err)
	}
}

// 测试插入表行数据
func TestCreate(t *testing.T) {
	var err error
	var db *gorm.DB
	if db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		t.Log(err)
	}
	v := User{
		Ver: "v0.2",
		Md5: "2",
		Url: "127.0.0.2",
	}
	if err = db.Create(&v).Error; err != nil {
		t.Log(err)
		return
	}
}

// 测试相关方法
func TestMethod(t *testing.T) {
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// 存在
	if db.Migrator().HasTable(&User{}) {
		t.Log("exist")
	} else {
		t.Log("no exist")
	}
	if db.Migrator().HasTable("user") {
		t.Log("exist")
	}
	// 删除表
	if err := db.Migrator().DropTable("user"); err != nil {
		t.Log("drop", err)
	}
	// 创建表
	if err := db.Migrator().CreateTable(&User{}); err != nil {
		t.Log(err)
	}
	if err := db.Migrator().CreateTable(&User{}); err != nil {
		t.Log(err)
	}
	// 创建索引
	if err := db.Migrator().CreateIndex(&User{}, "Ver"); err != nil {
		t.Log(err)
	}
	if err := db.Migrator().CreateIndex(&User{}, "idx_ver"); err != nil {
		t.Log(err)
	}
}

func TestQuery(t *testing.T) {
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// take 查询一条记录
	v := &User{}
	db.Take(v)
	t.Log(v)
	// first 查询一条记录，根据主键ID排序(正序)，返回第一条记录
	v2 := &User{}
	db.First(v2)
	t.Log(v2)
	// last 查询一条记录, 根据主键ID排序(倒序)，返回第一条记录
	v3 := &User{}
	db.Last(v3)
	t.Log(v3)
	// find 查询多条记录，Find函数返回的是一个数组
	var users []User
	if err := db.Find(&users).Error; err != nil {
		t.Log(err)
	}
	t.Log(users)
	// pluck 查询一列值
	var md5s []string
	db.Model(&User{}).Pluck("md5", &md5s)
	t.Log(md5s)
}

func TestWhere(t *testing.T) {
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//例子1:
	//等价于: SELECT * FROM `user`  WHERE (md5 = '1') LIMIT 1
	//这里问号(?), 在执行的时候会被1替代
	v := User{}
	db.Where("md5 = ?", 1).Take(&v)
	t.Log(v)
	//例子2:
	// in 语句
	//等价于: SELECT * FROM `user`  WHERE (md5 in ('1','2','5','6', '19')) LIMIT 1
	var v2 User
	db.Where("md5 in (?)", []string{"1", "2", "5", "6", "19"}).Take(&v2)
	t.Log(v2)
	//例子3:
	//等价于: SELECT * FROM `user`  WHERE (md5 >= '1' and md5 <= '19')
	//这里使用了两个问号(?)占位符，后面传递了两个参数替换两个问号。
	var v3 []User
	db.Where("md5 >= ? and md5 <= ?", "1", "19").Find(&v3)
	t.Log(v3)
	//例子4:
	//like语句
	//等价于: SELECT * FROM `user`  WHERE (md5 like '127%')
	var v4 []User
	db.Where("url like ?", "127%").Find(&v4)
	t.Log(v4)
}

func TestSelect(t *testing.T) {
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//例子1:
	//等价于: SELECT ver,md5 FROM `user`  WHERE `user`.`md5` = '1' LIMIT 1
	v1 := User{}
	db.Select("ver,md5").Where("md5 = ?", "1").Take(&v1)
	v2 := User{}
	db.Select([]string{"ver", "md5"}).Where("md5 = ?", "19").Take(&v2)
	//例子2:
	//等价于: SELECT count(*) as total FROM `user`
	total := []int{}
	db.Model(&User{}).Select("count(*) as total").Pluck("md5", &total)

	t.Log(v1, v2, total)
}

func TestOrder(t *testing.T) {
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//例子:
	//等价于: SELECT * FROM `user`  WHERE (md5 >= '1') ORDER BY md5 desc
	v2 := []User{}
	db.Where("md5 >= ?", "1").Order("md5 desc").Find(&v2)
	t.Log(v2)
	v3 := []User{}
	//等价于: SELECT * FROM `user` ORDER BY md5 desc LIMIT 10 OFFSET 0
	db.Order("md5 desc").Limit(10).Offset(0).Find(&v3)
	t.Log(v3)

	//例子:
	var total int64 = 0
	//等价于: SELECT count(*) FROM `user`
	//这里也需要通过model设置模型，让gorm可以提取模型对应的表名
	db.Model(User{}).Count(&total)
	t.Log(total)
}

func TestUpdate(t *testing.T) {
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// v := User{}
	// db.Where("md5 = 1").Take(&v)
	// t.Log(v)
	// v.Md5 = "4"
	// db.Create(&v)

	// v1 := User{}
	// db.Where("md5 = ?", 2).Take(&v1)
	// db.Delete(&v1)
	db.Where("md5 = ?", "4").Delete(&User{})

}
