# DStorage-server
- 分布式存储,基于k8s docker grpc实现micro server
### File structure
```
|-- DStorage
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
    |-- model                       // 数据库操作
    |-- mq                          // 消息队列
    |-- release                     // 版本
    |-- service                     // 服务处理
    |   |-- download                    // 文件下载服务
    |   |-- upload                      // 文件上传服务
    |-- store                       // 存储配置
    |   |-- ceph            
    |   |-- oss
    |-- unit                        // 单元测试
    |-- .gitignore
    |-- LICENSE                     
    |-- README.md
```
### start 1.0 version
- [ ] 原生API,不使用框架
