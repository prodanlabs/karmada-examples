/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package k8s

import (
	"context"
	"encoding/base64"
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	"github.com/prodanlabs/kaadm/app/certs"
	"github.com/prodanlabs/kaadm/app/karmada"
	"github.com/prodanlabs/kaadm/app/utils"
)

const (
	ClusterName           = "karmada"
	UserName              = "admin"
	KArmadaKubeConfigName = "karmada-apiserver.config"
)

// install flags
var (
	Namespace          string
	KubeConfig         string
	KArmadaMasterIP    string
	KArmadaMasterPort  int32
	DataPath           string
	ExternalIP         string
	StorageClassesName string
	NodeSelectorLabels = map[string]string{}

	//etcd
	EtcdStorageMode string
	EtcdImage       string
	EtcdInitImage   string
	EtcdReplicas    int32
	EtcdDataPath    string
	EtcdStorageSize string

	//	karmada
	KArmadaApiServerImage                string
	KArmadaApiServerReplicas             int32
	KArmadaSchedulerImage                string
	KArmadaKubeControllerManagerImage    string
	KArmadaControllerManagerImage        string
	KArmadaControllerManagerReplicas     int32
	KArmadaKubeControllerManagerReplicas int32
	KArmadaSchedulerReplicas             int32
	KArmadaWebhookImage                  string
	KArmadaWebhookReplicas               int32
)

type InstallOptions struct {
	Namespace          string
	CertAndKeyFileName certs.CertAndKeyFileName
	CertAndKeyFileData map[string][]byte
	KubeClientSet      *kubernetes.Clientset
	RestConfig         *rest.Config
}

type InstallOptionsController interface {
	initialization() error
	CreateNamespace() error
	CreateServiceAccount(sa *[]corev1.ServiceAccount) error
	CreateClusterRole(clusterRole *rbacv1.ClusterRole) error
	CreateClusterRoleBinding(clusterRole *rbacv1.ClusterRoleBinding) error
	CreateSecret(secret *corev1.Secret) error
	CreateService(service *corev1.Service) error
}

func Deploy() {

	verifying()

	cert := &certs.Config{
		Namespace:                   Namespace,
		PkiPath:                     DataPath,
		KArmadaMasterIP:             KArmadaMasterIP,
		FlagsExternalIP:             ExternalIP,
		EtcdReplicas:                EtcdReplicas,
		EtcdStatefulSetName:         etcdStatefulSetName,
		EtcdServiceName:             etcdServiceName,
		KArmadaApiServerServiceName: kArmadaApiServerServiceName,
		KArmadaWebhookServiceName:   kArmadaWebhookServiceName,
	}

	certsName, err := cert.CertificateGeneration()
	if err != nil {
		klog.Exitln("Error generating certificate.", err)
	}

	i := &InstallOptions{}
	i.CertAndKeyFileName = *certsName

	if err = i.initialization(); err != nil {
		klog.Exitln(err)
	}

	// Create karmada kubeconfig
	serverURL := fmt.Sprintf("https://%s:%v", KArmadaMasterIP, KArmadaMasterPort)
	if err = utils.WriteKubeConfigFromSpec(serverURL, UserName, ClusterName, DataPath, KArmadaKubeConfigName, i.CertAndKeyFileData[i.CertAndKeyFileName.CACertFileName],
		i.CertAndKeyFileData[i.CertAndKeyFileName.KArmadaKeyFileName], i.CertAndKeyFileData[i.CertAndKeyFileName.KArmadaCertFileName]); err != nil {
		klog.Exitln("Failed to create karmada kubeconfig file.", err)
	}
	klog.Info("Create karmada kubeconfig success.")

	//	create ns
	if err = i.CreateNamespace(); err != nil {
		klog.Exitln(err)
	}

	// Create sa
	saSpec := i.ServiceAccountFromSpec([]string{kArmadaControllerManagerServiceName, KArmadaSchedulerServiceAccountName, kArmadaWebhookServiceAccountName})
	if err = i.CreateServiceAccount(saSpec); err != nil {
		klog.Exitln(err)
	}

	// Create karmada-controller-manager ClusterRole and ClusterRoleBinding
	prSpec := i.ClusterRoleFromSpec(kArmadaControllerManagerClusterRoleName, []rbacv1.PolicyRule{
		{
			APIGroups: []string{"*"},
			Resources: []string{"*"},
			Verbs:     []string{"get", "watch", "list", "create", "update", "delete"},
		},
		{
			NonResourceURLs: []string{"*"},
			Verbs:           []string{"get"},
		},
	})
	if err = i.CreateClusterRole(prSpec); err != nil {
		klog.Exitln(err)
	}

	// Create kubeconfig Secret
	kArmadaServerURL := fmt.Sprintf("https://%s.%s.svc.cluster.local:%v", kArmadaApiServerServiceName, i.Namespace, KArmadaMasterPort)
	config := utils.CreateWithCerts(kArmadaServerURL, UserName, UserName, i.CertAndKeyFileData[i.CertAndKeyFileName.CACertFileName],
		i.CertAndKeyFileData[i.CertAndKeyFileName.KArmadaKeyFileName], i.CertAndKeyFileData[i.CertAndKeyFileName.KArmadaCertFileName])
	configBytes, err := clientcmd.Write(*config)
	if err != nil {
		klog.Exitln("failure while serializing admin kubeConfig", err)
	}

	kubeConfigSecret := i.SecretFromSpec(kubeConfigSecretName, corev1.SecretTypeOpaque, map[string]string{kubeConfigSecretName: string(configBytes)})
	if err = i.CreateSecret(kubeConfigSecret); err != nil {
		klog.Exitln(err)
	}

	// Create cert Secret
	etcdCert := map[string]string{
		i.CertAndKeyFileName.CACertFileName:         string(i.CertAndKeyFileData[i.CertAndKeyFileName.CACertFileName]),
		i.CertAndKeyFileName.CAKeyFileName:          string(i.CertAndKeyFileData[i.CertAndKeyFileName.CAKeyFileName]),
		i.CertAndKeyFileName.EtcdServerCertFileName: string(i.CertAndKeyFileData[i.CertAndKeyFileName.EtcdServerCertFileName]),
		i.CertAndKeyFileName.EtcdServerKeyFileName:  string(i.CertAndKeyFileData[i.CertAndKeyFileName.EtcdServerKeyFileName]),
	}
	etcdSecret := i.SecretFromSpec(etcdCertSecretName, corev1.SecretTypeOpaque, etcdCert)
	if err = i.CreateSecret(etcdSecret); err != nil {
		klog.Exitln(err)
	}

	kArmadaCert := map[string]string{
		i.CertAndKeyFileName.CACertFileName:         string(i.CertAndKeyFileData[i.CertAndKeyFileName.CACertFileName]),
		i.CertAndKeyFileName.CAKeyFileName:          string(i.CertAndKeyFileData[i.CertAndKeyFileName.CAKeyFileName]),
		i.CertAndKeyFileName.KArmadaCertFileName:    string(i.CertAndKeyFileData[i.CertAndKeyFileName.KArmadaCertFileName]),
		i.CertAndKeyFileName.KArmadaKeyFileName:     string(i.CertAndKeyFileData[i.CertAndKeyFileName.KArmadaKeyFileName]),
		i.CertAndKeyFileName.EtcdClientCertFileName: string(i.CertAndKeyFileData[i.CertAndKeyFileName.EtcdClientCertFileName]),
		i.CertAndKeyFileName.EtcdClientKeFileName:   string(i.CertAndKeyFileData[i.CertAndKeyFileName.EtcdClientKeFileName]),
	}

	kArmadaSecret := i.SecretFromSpec(kArmadaCertSecretName, corev1.SecretTypeOpaque, kArmadaCert)
	if err = i.CreateSecret(kArmadaSecret); err != nil {
		klog.Exitln(kArmadaSecret)
	}

	kArmadaWebhookCert := map[string]string{
		"tls.crt": string(i.CertAndKeyFileData[i.CertAndKeyFileName.KArmadaCertFileName]),
		"tls.key": string(i.CertAndKeyFileData[i.CertAndKeyFileName.KArmadaKeyFileName]),
	}
	kArmadaWebhookSecret := i.SecretFromSpec(kArmadaWebhookCertSecretName, corev1.SecretTypeOpaque, kArmadaWebhookCert)
	if err = i.CreateSecret(kArmadaWebhookSecret); err != nil {
		klog.Exitln(kArmadaSecret)
	}

	// add node labels
	if err = AddNodeSelectorLabels(i.KubeClientSet); err != nil {
		klog.Exitf("Node failed to add '%s' label. %v", NodeSelectorLabels, err)
	}

	// Create etcd
	if err = i.CreateService(i.makeEtcdService(etcdServiceName)); err != nil {
		klog.Exitln(err)
	}

	klog.Info("create etcd StatefulSets")
	if _, err = i.KubeClientSet.AppsV1().StatefulSets(i.Namespace).Create(context.TODO(), i.makeETCDStatefulSet(), metav1.CreateOptions{}); err != nil {
		klog.Warning(err)
	}
	if WaitEtcdReplicasetInDesired(i.KubeClientSet, i.Namespace, utils.MapToString(etcdLabels), 30) != nil {
		klog.Warning(err)
	}
	if WaitPodReady(i.KubeClientSet, i.Namespace, utils.MapToString(etcdLabels), 30) != nil {
		klog.Warning(err)
	}

	// Create karmada-apiserver
	klog.Info("create karmada ApiServer Deployment")
	if err = i.CreateService(i.makeKArmadaApiServerService()); err != nil {
		klog.Exitln(err)
	}

	if _, err = i.KubeClientSet.AppsV1().Deployments(i.Namespace).Create(context.TODO(), i.makeKArmadaApiServerDeployment(), metav1.CreateOptions{}); err != nil {
		klog.Warning(err)
	}
	if WaitPodReady(i.KubeClientSet, i.Namespace, utils.MapToString(kArmadaApiServerLabels), 60) != nil {
		klog.Exitln(err)
	}

	//Create CRDS.  in karmada
	caBase64 := base64.StdEncoding.EncodeToString(i.CertAndKeyFileData[i.CertAndKeyFileName.CACertFileName])
	if err = karmada.InitKArmadaResources(i.Namespace, filepath.Join(DataPath, KArmadaKubeConfigName), caBase64); err != nil {
		klog.Exitln(err)
	}

	// Create karmada-kube-controller-manager
	// https://github.com/karmada-io/karmada/blob/master/artifacts/deploy/kube-controller-manager.yaml
	klog.Info("create karmada kube controller manager Deployment")
	if err = i.CreateService(i.kArmadaKubeControllerManagerService()); err != nil {
		klog.Exitln(err)
	}

	if _, err = i.KubeClientSet.AppsV1().Deployments(i.Namespace).Create(context.TODO(), i.makeKArmadaKubeControllerManagerDeployment(), metav1.CreateOptions{}); err != nil {
		klog.Warning(err)
	}
	if WaitPodReady(i.KubeClientSet, i.Namespace, utils.MapToString(kArmadaKubeControllerManagerLabels), 30) != nil {
		klog.Warning(err)
	}

	// Create karmada-scheduler
	// https://github.com/karmada-io/karmada/blob/master/artifacts/deploy/karmada-scheduler.yaml
	klog.Info("create karmada scheduler Deployment")
	if _, err = i.KubeClientSet.AppsV1().Deployments(i.Namespace).Create(context.TODO(), i.makeKArmadaSchedulerDeployment(), metav1.CreateOptions{}); err != nil {
		klog.Warning(err)
	}
	if WaitPodReady(i.KubeClientSet, i.Namespace, utils.MapToString(kArmadaSchedulerLabels), 30) != nil {
		klog.Warning(err)
	}

	// Create karmada-controller-manager
	// https://github.com/karmada-io/karmada/blob/master/artifacts/deploy/controller-manager.yaml
	klog.Info("create karmada controller manager Deployment")
	if _, err = i.KubeClientSet.AppsV1().Deployments(i.Namespace).Create(context.TODO(), i.makeKArmadaControllerManagerDeployment(), metav1.CreateOptions{}); err != nil {
		klog.Warning(err)
	}
	if WaitPodReady(i.KubeClientSet, i.Namespace, utils.MapToString(kArmadaControllerManagerLabels), 30) != nil {
		klog.Warning(err)
	}

	// Create karmada-webhook
	// https://github.com/karmada-io/karmada/blob/master/artifacts/deploy/karmada-webhook.yaml
	klog.Info("create karmada webhook Deployment")
	if err = i.CreateService(i.kArmadaWebhookService()); err != nil {
		klog.Exitln(err)
	}

	if _, err = i.KubeClientSet.AppsV1().Deployments(i.Namespace).Create(context.TODO(), i.makeKArmadaWebhookDeployment(), metav1.CreateOptions{}); err != nil {
		klog.Warning(err)
	}
	if WaitPodReady(i.KubeClientSet, i.Namespace, utils.MapToString(kArmadaWebhookLabels), 30) != nil {
		klog.Warning(err)
	}

	utils.GenExamples(DataPath)
}

func (i *InstallOptions) initialization() error {

	i.CertAndKeyFileData = map[string][]byte{}
	for _, v := range i.CertAndKeyFileName.ALLCertFileName {
		cert, err := utils.FileToBytes(DataPath, v)
		if err != nil {
			return errors.Errorf("Failed to get cert '%s'. %v", v, err)
		}
		i.CertAndKeyFileData[v] = cert
	}
	for _, v := range i.CertAndKeyFileName.ALLKeyFileName {
		key, err := utils.FileToBytes(DataPath, v)
		if err != nil {
			return errors.Errorf("Failed to get Key '%s'. %v", v, err)
		}
		i.CertAndKeyFileData[v] = key
	}

	i.Namespace = Namespace

	if !utils.PathIsExist(KubeConfig) {
		klog.Exitln("kubeconfig file does not exist.")
	}
	klog.Infof("kubeconfig file: %s", KubeConfig)

	restConfig, err := utils.RestConfig(KubeConfig)
	if err != nil {
		return err
	}
	i.RestConfig = restConfig

	clientSet, err := utils.NewClientSet(restConfig)
	if err != nil {
		return err
	}
	i.KubeClientSet = clientSet

	return nil
}

func verifying() {

	if KArmadaMasterIP == "" {
		klog.Exitln("error verifying flag, master is missing. See 'kaadm install --help'.")
	}
	if !utils.IsIP(KArmadaMasterIP) {
		klog.Exitf("error verifying flag value, '%s' is not a valid flag value. See 'kaadm install --help'.", KArmadaMasterIP)
	}

	if EtcdStorageMode == "hostPath" && EtcdReplicas != 1 {
		klog.Exitln("When etcd storage mode is hostPath, replicas is 1. See 'kaadm install --help'.")
	}

	if EtcdStorageMode == "hostPath" && EtcdDataPath == "" {
		klog.Exitln("When etcd storage mode is hostPath, dataPath is not empty. See 'kaadm install --help'.")
	}

	if EtcdStorageMode == "PVC" && StorageClassesName == "" {
		klog.Exitln("When etcd storage mode is PVC, storageClassesName is not empty. See 'kaadm install --help'.")
	}
}
