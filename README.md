# go-utils

[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/tryturned/go-utils)

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

## 👋 Contributors

- 所有的 git 提交使用 [Commitizen](https://github.com/commitizen/cz-cli) 工具进行格式化提交信息

- 保证包不会阻塞主程序的正常运行, 除系统错误以外, 包产生的错误需要返回给业务方，由业务进行错误处理

- 如非必要, 包尽量不要产生日志输出

