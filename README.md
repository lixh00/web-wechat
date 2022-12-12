## WebWechat

基于`Golang`语言和`Gin`框架的个人微信系统，微信协议基于[openwechat](https://github.com/eatmoreapple/openwechat)

## 阿巴阿巴
2022-12-12，跟个风对接了`ChatGPT`，指令格式(@openai是固定句式，不代表任何一个微信账号)：
```text
@openai
你要说的话
```

## 使用方式

```shell
# 下载代码
git clone https://xxxx.xxx.xx/xx.git
# 更新依赖
go mod download
# 编译
go build main.go
# 清理无用mod引用
go mod tidy
```

## Docker一键运行

```shell
docker compose up -d
# 或者
docker-compose up -d
```

## Thanks

<a href="https://www.jetbrains.com/?from=openwechat"><img src="https://account.jetbrains.com/static/images/jetbrains-logo-inv.svg" height="200" alt="JetBrains"/></a>
