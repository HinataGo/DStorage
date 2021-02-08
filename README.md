# DStorage-server
- 分布式存储,基于k8s docker grpc实现micro server
### File structure
```
|-- DStorage
    |-- api                         // 服务接口
    |   |-- download                    // 文件下载服务
    |   |-- upload                      // 文件上传服务
    |-- cmd                         // 执行脚本
    |-- config                      // 文件配置
    |-- router                      // 服务路由
    |-- db                          // 数据库连接
    |   |-- mysql                       // mysql数据存储
    |-- doc                         // 文档数据 & sql文件
    |-- encrypt                     // 数据加密
    |--format                       // 格式化数据
    |   |--json                         // json数据格式化 
    |-- handler                     // 路由处理服务  
    |-- log                         // 日志服务       
    |-- meta                        // 文件元数据
    |-- service                     // 业务逻辑
    |-- mq                          // 消息队列
    |-- release                     // 版本
    |-- store                       // 存储配置           
    |   |-- oss                         // 阿里云OSS
    |-- unit                        // 单元测试
    |-- .gitignore
    |-- LICENSE                     
    |-- README.md
```
### start 1.0 version
- [x] 基于Gin 快速开发第一版demo
- [x] 用户登录注册
- [x] 文件上传 ,分块上传,快速上传 ,断点续传
- [x] 文件存储
- [x] 文件下载功能
- [x] 第一版demo完成

### 2.0 version 微服务划分
- [ ] 服务划分
- [ ] 加入消息队列
- [ ] 引入grpc 进行rpc模块通信
- [ ] 服务注册
- [ ] 重构设计
- [ ] 添加单元测试
- [ ] docker部署, k8s
- [ ] 完成CI/CD持续集成