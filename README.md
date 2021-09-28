# goauth-server

云原生权限服务

- 完备的rbac模型
 
- 适配原生Oauth2.0

- 基于wolf-rbac apisix插件适配 apisix网关


## 本地开发

需要在根目录

1. 启动通过docker启动mysql,redis
```bash
docker-compose -f config/docker-compose.yml up
```
2. 运行服务
```bash
go run main.go
```

## TODO

1. 抽出配置

2. 规范编写注释

3. 梳理代码，完事基础功能