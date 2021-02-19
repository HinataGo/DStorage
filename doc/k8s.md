# Kubernetes 部署 user 服务
- 通过 yaml 文件描述配置过程和使用 kubectl 命令行工具访问 Kubernetes 的接口，user 服务的 Pod 描述如下：
```yaml
apiVersion: v1 
kind: Pod 
metadata: 
  name: user-service 
  labels: 
    name: user-service 
spec: 
  containers:                    #定义user容器，开放10086端口 
    - name: user 
      image: user 
      ports: 
        - containerPort: 10086 
      imagePullPolicy: IfNotPresent 
    - name: mysql                     #定义MySQL容器，开放3306端口 
      image: mysql-for-user 
      ports: 
        - containerPort: 3306 
      env: 
        - name: MYSQL_ROOT_PASSWORD 
          value: "123456" 
      imagePullPolicy: IfNotPresent 
    - name: redis                     #定义Redis容器，开放6379端口 
      image: redis:5.0 
      ports: 
        - containerPort: 6379 
      imagePullPolicy: IfNotPresent
```
- 由于在同一个 Pod 中的多个容器是并发启动的，为了保证 user 服务启动时 Redis 和 MySQL 数据库已经部署启动完成，在 user 服务的 main 函数中增加了 time.Sleep 延迟了 user 服务的启动。 
- 通过 kubectl create 命令和 yaml 描述启动 Pod
```shell
kubectl create -f user-service.yaml 
# 将在 Kubernetes 集群的 Node 节点中创建单个 Pod
```
- 通过以下两个命令可分别查看 user-service Pod 的信息和进入到 Pod 中:
```shell
kubectl get pod user-service  
kubectl exec -ti user-service -n default  -- /bin/bash 

```
- 单个 Pod 不具备自我恢复的能力，当 Pod 所在的 Node 出现问题，Pod 就很可能被删除，这就会导致 Pod 中容器提供的服务被终止。为了避免这种情况的发生，可以使用 Controller 来管理 Pod，Controller 提供创建和管理多个 Pod 的能力，从而使得被管理的 Pod 具备自愈和更新的能力。
```md
Replication Controller，确保用户定义的 Pod 副本数保持不变；
ReplicaSet，是 RC 的升级版，在选择器（Selector）的支持上优于 RC，RC 只支持基于等式的选择器，但 RS 还支持基于集合的选择器；
Deployment，在 RS 的基础上提供了 Pod 的更新能力，在 Deployment 配置文件中 Pod template 发生变化时，它能将现在集群的状态逐步更新成 Deployment 中定义的目标状态；
StatefulSets，其中的 Pod 是有序部署和具备稳定的标识，是一组存在状态的 Pod 副本。
```
- 使用 Deployment Controller 为我们管理 user-service Pod，配置如下：
    - 指定了 kind 的类型为 Deployment、副本的数量为 3 和选择器为匹配标签 name: user-service。可以发现原来 Pod 的配置放到了 template 标签下，并添加 name: user-service 的标签，Deployment Controller 将会使用 template 下的 Pod 配置来创建 Pod 副本，并通过标签选择器来监控 Pod 副本的数量，当副本数不足时，将会根据 template 创建 Pod。
```yml
    apiVersion: apps/v1 
    kind: Deployment 
    metadata: 
      name: user-service 
      labels: 
        name: user-service 
    spec: 
      replicas: 3 
      selector: 
          matchLabels: 
            name: user-service 
      template: 
        metadata: 
          labels: 
            name: user-service 
        spec: 
          containers:                    #定义user容器，开放10086端口 
            - name: user 
              image: user 
              ports: 
                - containerPort: 10086 
              imagePullPolicy: IfNotPresent 
            - name: mysql                     #定义MySQL容器，开放3306端口 
              image: mysql-for-user 
              ports: 
                - containerPort: 3306 
              env: 
                - name: MYSQL_ROOT_PASSWORD 
                  value: "123456" 
              imagePullPolicy: IfNotPresent 
            - name: redis                     #定义Redis容器，开放6379端口 
              image: redis:5.0 
              ports: 
                - containerPort: 6379 
              imagePullPolicy: IfNotPresent 

```


- 通过以下命令即可通过 Deployment Controller 管理 user-service Pod
```shell
kubectl create -f user-service-deployment.yaml 
```
- 可以通过 kubectl get Deployment 命令查看 user-service 的 Pod 副本状态
    - Deployment Controller 默认使用 RollingUpdate 策略更新 Pod，也就是滚动更新的方式；另一种更新策略是 Recreate，创建出新的 Pod 之前会先杀掉所有已存在的 Pod，可以通过 spec.strategy.type 标签指定更新策略。
```shell
kubectl get Deployment user-service 
```

- Deployment 的 rollout 当且仅当 Deployment 的 Pod template 中的 label 更新或者镜像更改时被触发，比如我们希望更新 redis 的版本：
    -  这将触发 user-service Pod 的重新更新部署。
```shell
kubectl set image deployment/user-service redis=redis:6.0 
```
- 当 Pod 被 Deployment Controller 管理时，单独使用 kubectl delete pod 无法删除相关 Pod，Deployment Controller 会维持 Pod 副本数量不变，这时则需要通过 kubectl delete Deployment 删除相关 Deployment 配置，比如删除 user-service 的 Deployment 配置，如下命令所示：
```shell
kubectl delete Deployment user-service 
```