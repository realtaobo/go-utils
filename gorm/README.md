# 使用gorm操作mysql

```go
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
```

## 1 连接mysql

```go
// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
```

## 2 Migration特性

GORM支持Migration特性，支持根据Go Struct结构自动生成对应的表结构。**GORM 的AutoMigrate函数，仅支持建表，不支持修改字段和删除字段，避免意外导致丢失数据。**

### 2.1 AutoMigrate方法

通过AutoMigrate函数可以快速建表，如果表已经存在不会重复创建。

```go
  // 根据User结构体，自动创建表结构.
  db.AutoMigrate(&User{})

  // 一次创建User、Product、Order三个结构体对应的表结构
  db.AutoMigrate(&User{}, &Product{}, &Order{})

  // 可以通过Set设置附加参数，下面设置表的存储引擎为InnoDB
  db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{})
```

### 2.2 Migrator().HasTable

检测表是否存在，如：

```go
  // 检测User结构体对应的表是否存在
  db.Migrator().HasTable(&User{})

  // 检测表名users是否存在
  db.Migrator().HasTable("users")
```

### 2.3 Migrator().CreateTable

```go
  // 根据User结构体建表
  db.Migrator().CreateTable(&User{})
```

### 2.4 Migrator().DropTable

```go
  // 删除User结构体对应的表
  db.Migrator().DropTable(&User{})

  // 删除表名为users的表
  db.Migrator().DropTable("users")
```

### 2.5 Migrator().DropColumn

```go
  // 删除User结构体对应表中的description字段
  db.Migrator().DropColumn(&User{}, "Name")
```

### 2.6 Migrator().CreateIndex

```go
  type User struct {
    gorm.Model
    Name string `gorm:"size:255;index:idx_name,unique"`
  }

  // 为 Name 字段创建索引
  db.Migrator().CreateIndex(&User{}, "Name")
  db.Migrator().CreateIndex(&User{}, "idx_name")

  // 为 Name 字段删除索引
  db.Migrator().DropIndex(&User{}, "Name")
  db.Migrator().DropIndex(&User{}, "idx_name")

  // 检查索引是否存在
  db.Migrator().HasIndex(&User{}, "Name")
  db.Migrator().HasIndex(&User{}, "idx_name")

  type User struct {
    gorm.Model
    Name  string `gorm:"size:255;index:idx_name,unique"`
    Name2 string `gorm:"size:255;index:idx_name_2,unique"`
  }
  // 修改索引名
  db.Migrator().RenameIndex(&User{}, "Name", "Name2")
  db.Migrator().RenameIndex(&User{}, "idx_name", "idx_name_2")
```

## 3 查询数据

<font color=red>注意：gorm库是协程安全的，gorm提供的函数可以并发的在多个协程安全的执行。</font>

表结构定义如下:

```go
  type User struct {
  // gorm.Model
  Ver string `gorm:"column:ver"`
  Md5 string `gorm:"column:md5;primaryKey"`
  Url string `gorm:"column:url"`
  }
  func (*User) TableName() string {
  return "user"
  }
```

### 3.1 简单查询(无where子句)

```go
  // take 查询一条记录
  v := &User{}
  db.Take(v)
  // first 查询一条记录，根据主键ID排序(正序)，返回第一条记录
  v2 := &User{}
  db.First(v2)
  // last 查询一条记录, 根据主键ID排序(倒序)，返回第一条记录
  v3 := &User{}
  db.Last(v3)
  // find 查询多条记录，Find函数返回的是一个数组
  var users []User
  if err := db.Find(&users).Error; err != nil {
    t.Log(err)
  }
  // pluck 查询一列值
  var md5s []string
  db.Model(&User{}).Pluck("md5", &md5s)
```

### 3.2 where子句

上面的例子都没有指定where条件，这里介绍下如何设置where条件，主要通过db.Where函数设置条件.
函数说明：  
`db.Where(query interface{}, args ...interface{})`

参数说明:

| 参数名| 说明|
| :--- | :---: |
| query | sql语句的where子句, where子句中使用问号(?)代替参数值，则表示通过args参数绑定参数
| args | where子句绑定的参数，可以绑定多个参数

```go
  // 例子1:
  // 等价于: SELECT * FROM `user`  WHERE (md5 = '1') LIMIT 1
  // 这里问号(?), 在执行的时候会被1替代
  v := User{}
  db.Where("md5 = ?", 1).Take(&v)
  t.Log(v)
  // 例子2:
  // in 语句
  // 等价于: SELECT * FROM `user`  WHERE (md5 in ('1','2','5','6', '19')) LIMIT 1
  var v2 User
  db.Where("md5 in (?)", []string{"1", "2", "5", "6", "19"}).Take(&v2)
  t.Log(v2)
  // 例子3:
  // 等价于: SELECT * FROM `user`  WHERE (md5 >= '1' and md5 <= '19')
  // 这里使用了两个问号(?)占位符，后面传递了两个参数替换两个问号。
  var v3 []User
  db.Where("md5 >= ? and md5 <= ?", "1", "19").Find(&v3)
  t.Log(v3)
  // 例子4:
  // like语句
  // 等价于: SELECT * FROM `user`  WHERE (md5 like '127%')
  var v4 []User
  db.Where("url like ?", "127%").Find(&v4)
  t.Log(v4)
```

### 3.3 select

设置select子句, 指定返回的字段

```go
  // 例子1:
  // 等价于: SELECT ver,md5 FROM `user`  WHERE `user`.`md5` = '1' LIMIT 1
  v1 := User{}
  db.Select("ver,md5").Where("md5 = ?", "1").Take(&v1)
  v2 := User{}
  // 这种写法是直接往Select函数传递数组，数组元素代表需要选择的字段名
  db.Select([]string{"ver", "md5"}).Where("md5 = ?", "19").Take(&v2)

  // 例子2:
  // 可以直接书写聚合语句
  // 等价于: SELECT count(*) as total FROM `user`
  total := []int{}
  // Model函数，用于指定绑定的模型，这里生成了一个Food{}变量。目的是从模型变量里面提取表名，Pluck函数我们没有直接传递绑定表名的结构体变量，gorm库不知道表名是什么，所以这里需要指定表名
  db.Model(&User{}).Select("count(*) as total").Pluck("md5", &total)
```

### 3.4 order

设置排序语句，order by子句

```go
  v2 := []User{}
  db.Where("md5 >= ?", "1").Order("md5 desc").Find(&v2)
```

### 3.5 limit & Offset

设置limit和Offset子句，分页的时候常用语句。

```go
  v3 := []User{}
  //等价于: SELECT * FROM `user` ORDER BY md5 desc LIMIT 10 OFFSET 0
  db.Order("md5 desc").Limit(10).Offset(0).Find(&v3)
```

### 3.6 count

Count函数，直接返回查询匹配的行数。

```go
  db.Model(User{}).Count(&total)
```

### 3.7 分组

设置group by子句

```go
  db.Model(*).Select("type, count(*) as  total").Group("type").Having("total > 0").Scan(&results)
```

## 4 直接执行sql语句

对于复杂的查询，例如多表连接查询，我们可以直接编写sql语句，然后执行sql语句。
gorm通过db.Raw设置sql语句，通过Scan执行查询。如:

```go
  sql := "SELECT count(*) as  total FROM `user` where create_time > ? HAVING (total > 0)"
  //因为sql语句使用了一个问号(?)作为绑定参数, 所以需要传递一个绑定参数(Raw第二个参数).
  //Raw函数支持绑定多个参数
  db.Raw(sql, "2018-11-06 00:00:00").Scan(&results)
```

<font color=red>注意：scan类似Find都是用于执行查询语句，然后把查询结果赋值给结构体变量，区别在于scan不会从传递进来的结构体变量提取表名。所以此处使用scan</font>

## 5 更新数据

### 5.1  Save

用于保存模型变量的值。相当于根据主键id，更新所有模型字段值。

```go
  v := User{}
  db.Where("md5 = 2").Take(&v)
  v.Ver = "v0.3"
  db.Save(&v)
```

## 5.2  Update

更新单个字段值

```go
  //例子1:
  //更新food模型对应的表记录
  //等价于: UPDATE `foods` SET `price` = '25'  WHERE `foods`.`id` = '2'
  db.Model(&food).Update("price", 25)
  //通过food模型的主键id的值作为where条件，更新price字段值。


  //例子2:
  //上面的例子只是更新一条记录，如果我们要更全部记录怎么办？
  //等价于: UPDATE `foods` SET `price` = '25'
  db.Model(&Food{}).Update("price", 25)
  //注意这里的Model参数，使用的是Food{}，新生成一个空白的模型变量，没有绑定任何记录。
  //因为Food{}的id为空，gorm库就不会以id作为条件，where语句就是空的

  //例子3:
  //根据自定义条件更新记录，而不是根据主键id
  //等价于: UPDATE `foods` SET `price` = '25'  WHERE (create_time > '2018-11-06 20:00:00') 
  db.Model(&Food{}).Where("create_time > ?", "2018-11-06 20:00:00").Update("price", 25)
```

### 5.3 Updates

更新多个字段值

```go
  //例子1：
  //通过结构体变量设置更新字段
  updataFood := Food{
    Price:120,
    Title:"柠檬雪碧",
  }

  //根据food模型更新数据库记录
  //等价于: UPDATE `foods` SET `price` = '120', `title` = '柠檬雪碧'  WHERE `foods`.`id` = '2'
  //Updates会忽略掉updataFood结构体变量的零值字段, 所以生成的sql语句只有price和title字段。
  db.Model(&food).Updates(&updataFood)

  //例子2:
  //根据自定义条件更新记录，而不是根据模型id
  updataFood := Food{
    Stock:120,
    Title:"柠檬雪碧",
  }
  
  //设置Where条件，Model参数绑定一个空的模型变量
  //等价于: UPDATE `foods` SET `stock` = '120', `title` = '柠檬雪碧'  WHERE (price > '10') 
  db.Model(&Food{}).Where("price > ?", 10).Updates(&updataFood)

  //例子3:
  //如果想更新所有字段值，包括零值，就是不想忽略掉空值字段怎么办？
  //使用map类型，替代上面的结构体变量

  //定义map类型，key为字符串，value为interface{}类型，方便保存任意值
  data := make(map[string]interface{})
  data["stock"] = 0 //零值字段
  data["price"] = 35

  //等价于: UPDATE `foods` SET `price` = '35', `stock` = '0'  WHERE (id = '2')
  db.Model(&Food{}).Where("id = ?", 2).Updates(data)
```

<font color=red>通过结构体变量更新字段值, gorm库会忽略零值字段。就是字段值等于0, nil, "", false这些值会被忽略掉，不会更新。如果想更新零值，可以使用map类型替代结构体。</font>

### 5.4  更新表达式

UPDATE foods SET stock = stock + 1 WHERE id = '2'
这样的带计算表达式的更新语句gorm怎么写？

gorm提供了Expr函数用于设置表达式

```go
  //等价于: UPDATE `foods` SET `stock` = stock + 1  WHERE `foods`.`id` = '2'
  db.Model(&food).Update("stock", gorm.Expr("stock + 1"))
```

## 6 删除数据

### 6.1 删除模型数据

删除模型数据一般用于删除之前查询出来的模型变量绑定的记录。
用法：db.Delete(模型变量)

```go
  // 例子：
  v1 := User{}
  db.Where("md5 = ?", 2).Take(&v1)
  db.Delete(&v1)
```

### 6.2 根据Where条件删除数据

用法：db.Where(条件表达式).Delete(空模型变量指针)

```go
  // 等价于：DELETE from `user` where (`md5` = 5);
  db.Where("md5 = ?", "5").Delete(&User{})
```

## 7 事务处理

### 7.1 自动事务

通过db.Transaction函数实现事务，如果闭包函数返回错误，则回滚事务。

```go
  db.Transaction(func(tx *gorm.DB) error {
    // 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
    if err := tx.Create(&Animal{Name: "Giraffe"}).Error; err != nil {
      // 返回任何错误都会回滚事务
      return err
    }

    if err := tx.Create(&Animal{Name: "Lion"}).Error; err != nil {
      return err
    }
    // 返回 nil 提交事务
    return nil
  })
```

### 7.2 手动事务

在开发中经常需要数据库事务来保证多个数据库写操作的原子性。通常使用以下函数进行事务操作:

```go
// 在事务中执行数据库操作，使用的是tx变量，不是db。
// 开启事务
tx := db.Begin()
// 回滚事务
tx.Rollback()
// 提交事务
tx.Commit()
```

## 8 参考

[GORM快速入门教程](https://www.tizi365.com/archives/6.html)  

[GORM 指南](https://gorm.io/zh_CN/docs/)
