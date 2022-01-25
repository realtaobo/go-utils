package process

import (
	"fmt"

	cron "github.com/robfig/cron/v3"
)

// 需要运行的任务
var (
	tasks      = make(map[string]taskInfo)
	minuteCron = cron.New()
	secondCron = cron.New(cron.WithSeconds())
)

// taskInfo
type taskInfo struct {
	ID     cron.EntryID
	second bool
}

/*
RegisterTask 注册并启动定时任务
默认支持分钟级定时任务, 注册秒级定时任务请设置 second 参数为 true
taskName: 任务名称, 不允许重复
schedule: schedule, 调度规则, 分钟级别是5位, 秒级是6位, 详见https://pkg.go.dev/github.com/robfig/cron/v3#hdr-Usage
如 "0/10 * * * *" 表示每小时的0,20,30,40,50分执行一次定时任务
taskFunc: 用户指定的任务
*/
func RegisterTask(taskName string, schedule string, taskFunc func(), second bool) error {
	if _, ok := tasks[taskName]; ok {
		return fmt.Errorf("task already registered")
	}
	var entryId cron.EntryID
	var err error
	if second {
		entryId, err = secondCron.AddFunc(schedule, taskFunc)
		secondCron.Start()
	} else {
		entryId, err = minuteCron.AddFunc(schedule, taskFunc)
		minuteCron.Start()
	}
	if err != nil {
		return err
	}
	tasks[taskName] = taskInfo{
		ID:     entryId,
		second: second,
	}
	return nil
}

// 取消定时任务
func CancelTask(taskName string) error {
	if _, ok := tasks[taskName]; !ok {
		return fmt.Errorf("no such task")
	}
	taskID := tasks[taskName].ID
	second := tasks[taskName].second
	delete(tasks, taskName)
	if second {
		secondCron.Remove(taskID)
	} else {
		minuteCron.Remove(taskID)
	}
	return nil
}
