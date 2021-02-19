#!/bin/bash

kubectl apply -f svc_user.yaml
kubectl apply -f svc_auth.yaml
kubectl apply -f svc_proxy.yaml
kubectl apply -f svc_upload.yaml 
kubectl apply -f svc_download.yaml 
kubectl apply -f svc_transfer.yaml 
# 通知更新配置
kubectl apply -f service-ingress.yaml 
