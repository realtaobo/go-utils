# go-utils

[![Go Report Card](https://goreportcard.com/badge/github.com/tryturned/go-utils)](https://goreportcard.com/report/github.com/tryturned/go-utils)

## 🔮 Vision

The go-utils project aims to build a stupid golang pkg for myself.

## 💌 Features

| Feature                               | Status | Notes |
|---------------------------------------|--------|-------|
| **common**                            |
| json                                  | ✔ | convert the json string to the specified structure |
| yaml                                  | ✔ | convert the yaml string to the specified structure |
| **gorm**                              | ✔ | [gorm](https://github.com/go-gorm/gorm) Expansion pack
| **log**                               | ✔ | [logrus](https://github.com/sirupsen/logrus) Expansion pack
| **cron**                              | ✔ | [cron/v3](https://github.com/robfig/cron/v3) Expansion pack
| **sync**                             |
| Mutex                                 | ✔ | [sync.Mutex](https://pkg.go.dev/sync#Mutex) Expansion pack
| SpinLock                              | ✔ | Implementation of spin lock
| RecursiveMutex                        | ✔ | Implementation of reentrant lock
| RecursiveMutexByToken                 | ✔ | Implementation of token-based reentrant lock

## 👋 Contributors

- 所有的 git 提交使用 [Commitizen](https://github.com/commitizen/cz-cli) 工具进行格式化提交信息

- 保证包不会阻塞主程序的正常运行, 除系统错误以外, 包产生的错误需要返回给业务方，由业务进行错误处理

- 如非必要, 包尽量不要产生日志输出
