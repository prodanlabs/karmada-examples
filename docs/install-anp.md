##### clone code

```shell
git clone --branch v0.0.30 https://github.com/kubernetes-sigs/apiserver-network-proxy.git
```

##### generate certs

`PROXY_SERVER_HOST` is proxy server IP

```shell
cd apiserver-network-proxy
make certs PROXY_SERVER_IP={PROXY_SERVER_HOST}
```

##### deploy proxy-server

Copy the certificate to the `/etc/anp-server` directory where the `proxy-server` host is `PROXY_SERVER_HOST`

```shell
mkdir  /etc/anp-server
cp certs/frontend/issued/ca.crt  /etc/anp-server/server-ca.crt
cp certs/frontend/issued/proxy-frontend.crt /etc/anp-server/server-proxy-frontend.crt
cp certs/frontend/private/proxy-frontend.key /etc/anp-server/server-proxy-frontend.key
cp certs/agent/issued/ca.crt /etc/anp-server/cluster-ca.crt
cp certs/agent/issued/proxy-frontend.crt /etc/anp-server/cluster-proxy-frontend.crt
cp certs/agent/private/proxy-frontend.key /etc/anp-server/cluster-proxy-frontend.key
```
proxy-server.yaml `PROXY_SERVER_NAME` is `PROXY_SERVER_HOST` Node name

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: proxy-server
  namespace: karmada-system
spec:
  replicas: 1
  selector:
    matchLabels:
      app: proxy-server
  template:
    metadata:
      labels:
        app: proxy-server
    spec:
      containers:
      - args:
          - --health-port=8092
          - --proxy-strategies=destHost
          - --server-ca-cert=/certs/server-ca.crt
          - --server-cert=/certs/server-proxy-frontend.crt
          - --server-key=/certs/server-proxy-frontend.key
          - --cluster-ca-cert=/certs/cluster-ca.crt
          - --cluster-cert=/certs/cluster-proxy-frontend.crt
          - --cluster-key=/certs/cluster-proxy-frontend.key
        image: us.gcr.io/k8s-artifacts-prod/kas-network-proxy/proxy-server:v0.0.30
        imagePullPolicy: IfNotPresent
        livenessProbe:
          failureThreshold: 3
          httpGet:
            path: /healthz
            port: 8092
            scheme: HTTP
          initialDelaySeconds: 10
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 60
        name: proxy-server
        volumeMounts:
        - mountPath: /certs
          name: cert
      restartPolicy: Always
      hostNetwork: true
      nodeSelector:
        kubernetes.io/hostname: {PROXY_SERVER_NAME}
      volumes:
      - name: cert
        hostPath:
          path: /etc/anp-server
```

##### deploy proxy-agent

Packaging certs and upload it to the member cluster in pull mode
```shell
mkdir anp-agent
cp certs/agent/issued/ca.crt anp-agent/ca.crt
cp certs/agent/issued/proxy-agent.crt anp-agent/proxy-agent.crt
cp certs/agent/private/proxy-agent.key anp-agent/proxy-agent.key
tar -zcvf anp-agent.tar.gz anp-agent
```
in member cluster(pull mode) decompress
```shell
tar -zxvf anp-agent.tar.gz -C /etc
```

proxy-agent.yaml  `PROXY_AGENT_NAME` is the host Node name where the certificate is located.
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: proxy-agent
  name: proxy-agent
  namespace: karmada-system
spec:
  replicas: 1
  selector:
    matchLabels:
      app: proxy-agent
  template:
    metadata:
      labels:
        app: proxy-agent
    spec:
      containers:
      - args:
        - --ca-cert=/certs/ca.crt
        - --agent-cert=/certs/proxy-agent.crt
        - --agent-key=/certs/proxy-agent.key
        - --proxy-server-host={PROXY_SERVER_HOST}
        - --proxy-server-port=8091
        - --agent-identifiers=host=${HOST_IP}
        image: us.gcr.io/k8s-artifacts-prod/kas-network-proxy/proxy-agent:v0.0.30
        imagePullPolicy: IfNotPresent
        name: proxy-agent
        env:
          - name: HOST_IP
            valueFrom:
              fieldRef:
                fieldPath: status.hostIP
        livenessProbe:
          httpGet:
            scheme: HTTP
            port: 8093
            path: /healthz
          initialDelaySeconds: 15
          timeoutSeconds: 60
        volumeMounts:
          - mountPath: /certs
            name: cert
      nodeSelector:
        kubernetes.io/hostname: {PROXY_AGENT_NAME}
      volumes:
        - name: cert
          hostPath:
            path: /etc/anp-agent
```