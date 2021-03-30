# CECO: Cloud-Edge Coordination Operator

## 安装

安装：

> https://sdk.operatorframework.io/docs/installation/

## 试用

参考官方示例：https://sdk.operatorframework.io/docs/building-operators/golang/

环境信息:

* operator-sdk version v1.3.0
* go version go1.15.4 linux/amd64
* mercurial version 5.6.1-1.x86_64
* docker version 19.03.0
* kubectl version v1.18.3
* kubernetes cluster v1.18.3

``` BASH
# 生成crd清单
$ make generate
# 指定权限并生成RBAC清单
$ make manifests
# 构建并运行
$ make install
# 修改baseimage dockerfile以构建镜像
# 构建镜像(如果没有写test,可以注释掉docker-build里的test项,直接build)
$ make docker-build IMG=172.31.0.7:5000/ceco-controller:v0.0.1
# 推到镜像仓库
$ make docker-push IMG=172.31.0.7:5000/ceco-controller:v0.0.1
# 安装部署(默认config/default/manager_auth_proxy_patch.yaml中的镜像可能拉不下来，可改为其他镜像)
$ make deploy IMG=172.31.0.7:5000/ceco-controller:v0.0.1
# 查看 controller
$ kubectl get deploy -n ceco-system -owide
NAME                      READY   UP-TO-DATE   AVAILABLE   AGE    CONTAINERS                IMAGES                                                                    SELECTOR
ceco-controller-manager   1/1     1            1           115s   kube-rbac-proxy,manager   carlosedp/kube-rbac-proxy:v0.5.0,172.31.0.7:5000/ceco-controller:v0.0.1   control-plane=controller-manager
$ kubectl get po -n ceco-system -owide
NAME                                      READY   STATUS    RESTARTS   AGE    IP             NODE            NOMINATED NODE   READINESS GATES
ceco-controller-manager-69b6b6fb8-8qx4k   2/2     Running   0          104s   10.253.0.106   10.110.26.178   <none>           <none>
# 根据定义的type创建yaml,然后创建实例：
#   · 需要先部署好 nats server 集群，参见 https://docs.nats.io/nats-on-kubernetes/minimal-setup
#   · 目前版本 hostname 不支持ip地址的形式：
#       a DNS-1123 label must consist of lower case alphanumeric characters or '-', and must start and
#       end with an alphanumeric character (e.g. 'my-name',  or '123-abc', regex used for validation is 
#       '[a-z0-9]([-a-z0-9]*[a-z0-9])?')
#   · source的文件需要提前准备好
#   · 目前版本的云边镜像是写死在controller里的，需要先手动打好镜像：
#       $ cd deployment/nats/source
#       $ docker build -t 172.31.0.7:5000/source-nats:v0.0.2 .
#       $ docker push 172.31.0.7:5000/source-nats:v0.0.2
#       $ cd ../destination
#       $ docker build -t 172.31.0.7:5000/destination-nats:v0.0.2 .
#       $ docker push 172.31.0.7:5000/destination-nats:v0.0.2
$ kubectl apply -f config/samples/ceco_v1alpha1_natsco.yaml
# 查看部署结果
$ kubectl get deploy -w
# 查看日志
$ kubectl logs ceco-controller-manager-69b6b6fb8-8qx4k -c manager --tail 30 -n ceco-system -f
# 删除实例
$ kubectl delete -f config/samples/ceco_v1alpha1_natsco.yaml
# 删除部署
$ make undeploy
```