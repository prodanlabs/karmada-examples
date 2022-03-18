#!/usr/bin/env bash
set -o errexit
set -o nounset
set -o pipefail

mkdir ./manifests/webhook/test-certs
cd ./manifests/webhook/test-certs
openssl genrsa -out ca.key 2048
openssl req -x509 -new -nodes -key ca.key -subj "/C=CN/ST=Guangdong/L=Guangzhou/O=karmada/OU=System/CN=karmada" -days 3650 -out ca.crt

openssl genrsa -out tls.key 2048
openssl req -new -nodes -sha256 -subj "/C=CN/ST=Guangdong/L=Guangzhou/O=kubernetes/OU=System/CN=karmada"  -key tls.key -out tls.csr
openssl x509 -req -days 3650 \
  -extfile <(printf "keyUsage=critical,Digital Signature, Key Encipherment\nextendedKeyUsage=serverAuth,clientAuth\nauthorityKeyIdentifier=keyid,issuer\nsubjectAltName=DNS:karmada-custom-webhook.karmada-system.svc.cluster.local,DNS:localhost,IP:172.0.0.1,IP:172.31.6.145") \
  -sha256 -CA ca.crt -CAkey ca.key -set_serial 01 -in tls.csr -out tls.crt