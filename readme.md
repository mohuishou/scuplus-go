# scuplus 

> scuplus go 语言版本,微信小程序 后端api

## 目录结构

```
.
├── api                // 接口目录
├── cache              // 缓存接口，例如验证码缓存
├── ChangeLog.md       // 变更文档
├── config             // 配置目录
├── cron               // 定时任务
├── Dockerfile         
├── Gopkg.lock         
├── Gopkg.toml         
├── job                // job 异步任务
├── main.go
├── middleware         // 中间件
├── model              // model 数据层接口
├── readme.md
├── route              // 包含所有的api路由
├── task               // 单独执行任务(待废弃)
├── util               // 包含一些公共函数和单独的服务
└── vendor
```

## 错误码说明

```
1xxxx: 系统相关错误码
2xxxx: 教务系统相关
20xxx: 成绩相关
21xxx: 课程表相关
22xxx: 考表相关
3xxxx: 图书馆相关
4xxxx: 用户相关
5xxxx: 资讯相关
6xxxx: 一卡通相关
```