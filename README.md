# kaadm
karmada cli tool

```shell
 ~/kaadm-linux-amd64 # ./kaadm  install -h
Installation options.

Usage:
  kaadm install [flags]

Examples:
kaadm install --master=xxx.xxx.xxx.xxx

Flags:
      --cert-external-ip string                          the external IP of Karmada certificate (e.q 192.168.1.2,172.16.1.2)
      --etcd-data string                                 etcd data path,valid in hostPath mode. (default "/var/lib/karmada-etcd")
      --etcd-image string                                etcd image (default "k8s.gcr.io/etcd:3.5.1-0")
      --etcd-init-image string                           etcd init container image (default "docker.io/alpine:3.14.3")
      --etcd-replicas int32                              etcd replica set, cluster 3,5...singular (default 1)
      --etcd-storage-mode string                         etcd data storage mode(emptyDir,hostPath,PVC). value is PVC, specify --storage-classes-name;value is hostPath,--etcd-replicas is 1 (default "emptyDir")
      --etcd-storage-size string                         etcd data path,valid in pvc mode. (default "1Gi")
  -h, --help                                             help for install
      --karmada-apiserver-image string                   Kubernetes apiserver image (default "k8s.gcr.io/kube-apiserver:v1.20.11")
      --karmada-apiserver-replicas int32                 karmada apiserver replica set (default 1)
      --karmada-controller-manager-image string          karmada controller manager  image (default "swr.ap-southeast-1.myhuaweicloud.com/karmada/karmada-controller-manager:latest")
      --karmada-controller-manager-replicas int32        karmada controller manager replica set (default 1)
  -d, --karmada-data string                              karmada data path. kubeconfig and cert files (default "/var/lib/karmada")
      --karmada-kube-controller-manager-image string     Kubernetes controller manager image (default "k8s.gcr.io/kube-controller-manager:v1.20.11")
      --karmada-kube-controller-manager-replicas int32   karmada kube controller manager replica set (default 1)
      --karmada-scheduler-image string                   karmada scheduler image (default "swr.ap-southeast-1.myhuaweicloud.com/karmada/karmada-scheduler:latest")
      --karmada-scheduler-replicas int32                 karmada scheduler replica set (default 1)
      --karmada-webhook-image string                     karmada webhook image (default "swr.ap-southeast-1.myhuaweicloud.com/karmada/karmada-webhook:latest")
      --karmada-webhook-replicas int32                   karmada webhook replica set (default 1)
      --kubeconfig string                                absolute path to the kubeconfig file (default "/root/.kube/config")
      --master string                                    Karmada master ip. (e.g. --master 192.168.1.2)
  -n, --namespace string                                 Kubernetes namespace (default "karmada-system")
  -p, --port int32                                       Karmada apiserver port (default 5443)
      --storage-classes-name string                      Kubernetes StorageClasses Name
```

### Examples
```shell
kaadm install --master=192.168.1.2 --etcd-replicas 3
```
* 192.168.1.2 is karmada apiserver node,kaadm finds node by IP and adds `karmada.io/master=` tag,Binding by label.
* `etcd-replicas=3` Is a 3-node etcd cluster. 

If nothing happens, print the following information.
```shell
┌─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
| Push mode                                                                                                                                                           |
|                                                                                                                                                                     |
| Step 1: Member kubernetes join karmada control plane                                                                                                                |
|                                                                                                                                                                     |
| (In karmada)~#  cat ~/.kube/config  | grep current-context | sed 's/: /\n/g'| sed '1d' #MEMBER_CLUSTER_NAME                                                         |
| (In karmada)~# kubectl-karmada  --kubeconfig /var/lib/karmada/karmada-apiserver.config  join ${MEMBER_CLUSTER_NAME} --cluster-kubeconfig=$HOME/.kube/config         |
|                                                                                                                                                                     |
| Step 2: Create member kubernetes kubeconfig secret                                                                                                                  |
|                                                                                                                                                                     |
| (In member kubernetes)~# kubectl create ns karmada-system                                                                                                           |
| (In member kubernetes)~# kubectl create secret generic ${MEMBER_CLUSTER_NAME}-kubeconfig --from-file=${MEMBER_CLUSTER_NAME}-kubeconfig=$HOME/.kube/config  -n karmada-system              |
|                                                                                                                                                                     |
| Step 3: Create karmada scheduler estimator                                                                                                                          |
|                                                                                                                                                                     |
| (In member kubernetes)~# sed -i "s/{{member_cluster_name}}/${MEMBER_CLUSTER_NAME}/g" /var/lib/karmada/karmada-scheduler-estimator.yaml                              |
| (In member kubernetes)~# kubectl create -f  /var/lib/karmada/karmada-scheduler-estimator.yaml                                                                       |
|                                                                                                                                                                     |
| Step 4: Show members of karmada                                                                                                                                     |
|                                                                                                                                                                     |
| (In karmada)~# kubectl  --kubeconfig /var/lib/karmada/karmada-apiserver.config get clusters                                                                         |
|                                                                                                                                                                     |
├── —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —──┤
| Pull mode                                                                                                                                                           |
|                                                                                                                                                                     |
| Step 1:  Send karmada kubeconfig and karmada-agent.yaml to member kubernetes                                                                                        |
|                                                                                                                                                                     |
| (In karmada)~# scp /var/lib/karmada/karmada-apiserver.config /var/lib/karmada/karmada-agent.yaml {member kubernetes}:~                                              |
|                                                                                                                                                                     |
| Step 2:  Create karmada kubeconfig secret                                                                                                                           |
|  Notice:                                                                                                                                                            |
|    Cross-network, need to change the config server address.                                                                                                         |
|                                                                                                                                                                     |
| (In member kubernetes)~#  kubectl create ns karmada-system                                                                                                          |
| (In member kubernetes)~#  kubectl create secret generic karmada-kubeconfig --from-file=karmada-kubeconfig=/root/karmada-apiserver.config  -n karmada-system         |
|                                                                                                                                                                     |
| Step 3: Create karmada agent                                                                                                                                        |
|                                                                                                                                                                     |
| (In member kubernetes)~#  MEMBER_CLUSTER_NAME="demo"                                                                                                                |
| (In member kubernetes)~#  sed -i "s/{member_cluster_name}/${MEMBER_CLUSTER_NAME}/g" karmada-agent.yaml                                                              |
| (In member kubernetes)~#  kubectl create -f karmada-agent.yaml                                                                                                      |
|                                                                                                                                                                     |
| Step 4: Show members of karmada                                                                                                                                     |
|                                                                                                                                                                     |
| (In karmada)~# kubectl  --kubeconfig /var/lib/karmada/karmada-apiserver.config get clusters                                                                         |
└─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘

```

### plan
- [ √ ]  karmada in kubernetes.
- [ ] karmada in linux.