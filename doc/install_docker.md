# docker 安装配置操作
### 参考官网

### 免sudo 使用
- 临时
```shell
 sudo gpasswd -a $USER docker
 newgrp docker
```
- 永久
```shell
sudo chmod a+rw /var/run/docker.sock
```

### docker 常用命令
```shell
docker ps   #  查看本地docker运行情况
docker images     #  查看本地所有镜像
docker rm xx(id)  # 删除执行容器 记住这里后面用 镜像id
docker stop # 命令用于停止容器
docker kill # 命令用于杀死容器
docker attach # 命令用于连接运行中容器
docker rmi     # 删除镜像
```
- 使用 docker exec 命令在容器中执行命令，比如我们可以在容器中启动一个 /bin/bash 交互式终端，从而进入到容器中执行各种命令
- 进入到 redis-5.0 容器之后，就可以通过 redis-cli 命令访问容器内的 Redis 数据库
```shell
docker exec -it redis-5.0 /bin/bash 
```
### none 镜像怎么删？
- 原因： 进行拉取未命名 ，创建未命名， 打包中途失败，虽然失败，但是会留下容器，执行到一半退出， stop，没用

```shell
docker ps -a
# 找到无关的容器
docker rm id
# 随后 找到none 的id
docker images -a
# 删除镜像
docker rmi id

```
### Error response from daemon: conflict: unable to delete f78e2d2133a2 (cannot be forced) - image has dependent child images
- 找到相关的依赖容器， 删除主的， 与其相关的none 也会自动被删除
```shell
# 另外的 image FROM 了这个 image，可以使用下面的命令列出所有在指定 image 之后创建的 image 的父 image
docker image inspect --format='{{.RepoTags}} {{.Id}} {{.Parent}}' $(docker image ls -q --filter since=xxxxxx)
# 删除关联的依赖镜像，关联的none镜像也会被删除
docker rmi id
#  这时查看镜像就不会有之前的问题了
docker images -a
```
### 其他docker操作
```shell
# 停止所有容器
docker ps -a | grep "Exited" | awk '{print $1 }'|xargs docker stop

# 删除所有容器
docker ps -a | grep "Exited" | awk '{print $1 }'|xargs docker rm

# 删除所有none容器
docker images|grep none|awk '{print $3 }'|xargs docker rmi
```