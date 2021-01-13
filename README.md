# Blogger

一个用 golang 实现的从 git 仓库生成静态博客的 webhook 服务.

## 部署

下载最新的可执行程序。

配置信息 config.toml

```toml
[Git]
# 仓库地址
Repo = ""
# 如果是私有的须配置 
AccessToken = ""
# 仓库下载本地位置
Dir = "./repo"

[Blog]
Host = "0.0.0.0"
Port = 8080
# 站点模板位置
Template = "./repo/template"
# 站点输出位置
Dir = "./out"
# webhook 触发 token
WebHookAccessToken = "abcd"
```

启动服务

```shell
nohup ./main >> log.output &
```

## 配置 Github

在GitHub配置好 webhook 后, 每次推送将触发服务，自动生成静态内容。

```text
Repository > Settings > Webhooks > Add Webhooks 
```
