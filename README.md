# go-utils

[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/tryturned/go-utils)

## 代码结构说明

```bash
# 各项详情详见子文件夹 README
.
├── README.md
├── common  # 通用公共文件夹, 需完全独立
├── gorm    # gorm 操作 mysql 的相关接口封装
├── log     # 对logrus的简单封装
└── process # golang 操作系统相关接口实现
```

## 仓库规范

- 所有的 git 提交使用 [Commitizen](https://github.com/commitizen/cz-cli) 工具进行格式化提交信息

- 保证包不会阻塞主程序的正常运行, 除系统错误以外, 包产生的错误需要返回给业务方，由业务进行错误处理

- 如非必要, 包尽量不要产生日志输出

- 尽量包含单元测试
