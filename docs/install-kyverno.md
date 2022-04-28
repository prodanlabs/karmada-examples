#### 1、Create namespace in karmada control plane
```shell
kubectl create -f kyverno-ns.yaml  --kubeconfig /etc/karmada/karmada-apiserver.config
```

#### 2、Create kyverno  configmap in karmada control plane
```shell
kubectl create -f kyverno-cm.yaml  --kubeconfig /etc/karmada/karmada-apiserver.config
```
#### 3、Create kyverno  crds in karmada control plane
```shell
kubectl create -f kyverno-crds.yaml --kubeconfig /etc/karmada/karmada-apiserver.config
```

#### 4、Create secrets in hosts cluster
```shell
kubectl create secret generic karmada-config  --from-file=/etc/karmada/karmada-apiserver.config -n kyverno
```

#### 5、Create  kyverno rbac in hosts cluster
```shell
kubectl create -f kyverno-rbac.yaml
```


#### 6、Create  kyverno  deployment in hosts cluster
```shell
kubectl create -f kyverno.yaml
```

>Need to modify `kyverno.yaml`


* Service changed to `NodePort` type
```yaml
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: kyverno
    app.kubernetes.io/component: kyverno
    app.kubernetes.io/instance: kyverno
    app.kubernetes.io/name: kyverno
    app.kubernetes.io/part-of: kyverno
    app.kubernetes.io/version: latest
  name: kyverno-svc
  namespace: kyverno
spec:
  type: NodePort
  ports:
  - name: https
    port: 443
    targetPort: https
    nodePort: 33443
  selector:
    app: kyverno
    app.kubernetes.io/name: kyverno
```

`--serverIP` needs to be changed to your {NODE_IP}:{NODE_PORT}
```shell
        - --autogenInternals=false
        - --kubeconfig=/etc/karmada/karmada-apiserver.config
        - --serverIP=172.31.6.145:33443
```

yaml see: https://github.com/prodanlabs/karmada-examples/tree/main/manifests/kyverno