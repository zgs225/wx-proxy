微信公众号授权域名代理
===

在做微信公众号开发的时候，由于平台只支持设置两个域名，所以在很多时候会不够用。
所以使用这个项目做一个代理登录。

## 安装和使用

```
get get -u github.com/zgs225/wx-proxy
wxproxy --config /path/to/config
```

## 配置

``` yaml
# 微信公众号 APP ID
app_id: APP_ID
# WEB 根目录路径
web_root_dir: .
# 重定向域名白名单
allow_hosts:
  - "example.com"
# 代理重定向域名
host: http://example.com
```

## OAuth2

访问链接 `${部署域名}/connect/oauth2/authorize?redirect_to=REDIRECT_TO`

## 变更记录

### v1.0.1

+ 添加 `host` 配置，作为重定向域名前缀

### v1.0.0

+ 完成基本的代理功能
