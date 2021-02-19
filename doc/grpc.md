## 安装工具
### grpc
[官网](https://www.grpc.io/docs/languages/go/quickstart/)
- 推荐看GitHub
  [github](https://github.com/grpc/grpc-go)
  
- 生成对应的 grpc pb 文件
```shell
# 工具依赖
go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
go get -u google.golang.org/grpc
# 生成pb 文件命令，先cd 到指定目录
protoc --go_out=plugins=grpc:. *.proto
```
### 利用protoc 生成 pb 文件
[官网](https://developers.google.com/protocol-buffers/docs/gotutorial)
```shell
protoc -I=$SRC_DIR --go_out=$DST_DIR $SRC_DIR/addressbook.proto
```
- proto文件记得加上 go package