# Blogger

一个用 golang 实现的从 git 仓库生成静态博客的 webhook 服务, 让博客搭建, 发布, 更新更加便捷.

## 快速部署

在 [这里](https://github.com/dengzii/blogger/releases/) 找到最新的版本连接

1. 下载
```shell
wget https://github.com/dengzii/blogger/releases/download/v1.1.0/blogger-v1.1.0.rar
```
2. 解压
```shell
tar -C /blogger -xzvf blogger-v1.1.0.tar.gz
```
3. 在 config.toml 中配置你的博客仓库及 webhook 服务信息
```shell
cd blogger
vim config.toml
```
4. 配置 nginx, 代理 webhook 服务
```shell
server {
    listen 8082;
    server_name _;
    access_log /srv/blog/nginx.log;
    location / {
        proxy_pass http://127.0.0.1:8088;
    }
}
```
5. 运行
```shell
chmod -R 777 blogger run.sh
./run.sh
```
6. 在 GitHub 配置 webhook, 每次推送将自动更新博客内容

https://github.com/用户名称/仓库名称/settings/hooks/new

```text
Repository > Settings > Webhooks > Add Webhooks 
```
