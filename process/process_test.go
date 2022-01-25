package process

import (
	"log"
	"testing"
	"time"
)

// 测试定时任务
func TestRegisterTask(t *testing.T) {
	// 每两秒执行一次, 10s后取消任务执行
	// 该注册执行5次
	err := RegisterTask("Hello", "*/2 * * * * *", func() { log.Println("Hello World!") }, true)
	if err != nil {
		t.Log(err)
	}
	time.Sleep(2 * time.Second)
	// 继续注册,该注册执行次数小于等于5即正常
	err = RegisterTask("Second", "*/2 * * * * *", func() { log.Println("Second World!") }, true)
	if err != nil {
		t.Log(err)
	}
	time.Sleep(10 * time.Second)
	CancelTask("Hello")
}
