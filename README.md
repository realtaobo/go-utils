## go-utils
| pkg| 说明| 
| :--- | :---: | 
|filewatch |对文件进行检测，当发现配置文件更改时调用自定义方法 |
|gorm|gorm 模块的学习使用

## 规范
### Commit message 的格式
每次提交，Commit message 都包括三个部分：Header，Body 和 Footer。
```bash
<type>(<scope>): <subject>
// 空一行
<body>
// 空一行
<footer>
```
其中，Header 是必需的，Body 和 Footer 可以省略。
不管是哪一个部分，任何一行都不得超过72个字符（或100个字符）。这是为了避免自动换行影响美观。

其中，Header部分只有一行，包括三个字段：`type（必需）`、`scope（可选）`和`subject（必需）`。`type`用于说明 commit 的类别，只允许使用`feat、fix、docs、style、refactor、test、chore`7个标识。具体请参考：[Commit message 和 Change log 编写指南](https://www.ruanyifeng.com/blog/2016/01/commit_message_change_log.html)
### Commitizen
[Commitizen](https://github.com/commitizen/cz-cli)是一个撰写合格 Commit message 的工具。
```bash
npm install -g commitizen cz-conventional-changelog
echo '{ "path": "cz-conventional-changelog" }' >>
git cz #replace git commit
```