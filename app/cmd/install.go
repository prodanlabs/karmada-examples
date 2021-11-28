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

package cmd

import (
	"github.com/spf13/cobra"
	"path/filepath"

	"github.com/prodanlabs/kaadm/app/k8s"
	"github.com/prodanlabs/kaadm/app/utils"
)

// newCmdInstall install karmada.
func newCmdInstall() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "install",
		Short: "bootstrap install karmada (default in kubernetes)",
		Long:  `Installation options.`,
		Run: func(cmd *cobra.Command, args []string) {
			k8s.Deploy()
		},
		Example: "kaadm install --master=xxx.xxx.xxx.xxx",
	}

	// cert
	cmd.PersistentFlags().StringVar(&k8s.ExternalIP, "cert-external-ip", "", "the external IP of Karmada certificate (e.g 192.168.1.2,172.16.1.2)")

	// Kubernetes
	cmd.PersistentFlags().StringVarP(&k8s.KubeConfig, "kubeconfig", "", filepath.Join(utils.HomeDir(), ".kube", "config"), "absolute path to the kubeconfig file")
	cmd.PersistentFlags().StringVarP(&k8s.Namespace, "namespace", "n", "karmada-system", "Kubernetes namespace")
	cmd.PersistentFlags().StringVar(&k8s.StorageClassesName, "storage-classes-name", "", "Kubernetes StorageClasses Name")

	// etcd
	cmd.PersistentFlags().StringVarP(&k8s.EtcdStorageMode, "etcd-storage-mode", "", "emptyDir",
		"etcd data storage mode(emptyDir,hostPath,PVC). value is PVC, specify --storage-classes-name;value is hostPath,--etcd-replicas is 1")
	cmd.PersistentFlags().StringVarP(&k8s.EtcdImage, "etcd-image", "", "k8s.gcr.io/etcd:3.5.1-0", "etcd image")
	cmd.PersistentFlags().StringVarP(&k8s.EtcdInitImage, "etcd-init-image", "", "docker.io/alpine:3.14.3", "etcd init container image")
	cmd.PersistentFlags().Int32VarP(&k8s.EtcdReplicas, "etcd-replicas", "", 1, "etcd replica set, cluster 3,5...singular")
	cmd.PersistentFlags().StringVarP(&k8s.EtcdDataPath, "etcd-data", "", "/var/lib/karmada-etcd", "etcd data path,valid in hostPath mode.")
	cmd.PersistentFlags().StringVarP(&k8s.EtcdStorageSize, "etcd-storage-size", "", "1Gi", "etcd data path,valid in pvc mode.")

	// karmada
	cmd.PersistentFlags().StringVar(&k8s.KArmadaMasterIP, "master", "", "Karmada master ip. (e.g. --master 192.168.1.2,192.168.1.3)")
	cmd.PersistentFlags().Int32VarP(&k8s.KArmadaMasterPort, "port", "p", 5443, "Karmada apiserver port")
	cmd.PersistentFlags().StringVarP(&k8s.DataPath, "karmada-data", "d", "/var/lib/karmada", "karmada data path. kubeconfig and cert files")
	cmd.PersistentFlags().StringVarP(&k8s.KArmadaApiServerImage, "karmada-apiserver-image", "", "k8s.gcr.io/kube-apiserver:v1.20.11", "Kubernetes apiserver image")
	cmd.PersistentFlags().Int32VarP(&k8s.KArmadaApiServerReplicas, "karmada-apiserver-replicas", "", 1, "karmada apiserver replica set")
	cmd.PersistentFlags().StringVarP(&k8s.KArmadaSchedulerImage, "karmada-scheduler-image", "", "swr.ap-southeast-1.myhuaweicloud.com/karmada/karmada-scheduler:latest", "karmada scheduler image")
	cmd.PersistentFlags().Int32VarP(&k8s.KArmadaSchedulerReplicas, "karmada-scheduler-replicas", "", 1, "karmada scheduler replica set")
	cmd.PersistentFlags().StringVarP(&k8s.KArmadaKubeControllerManagerImage, "karmada-kube-controller-manager-image", "", "k8s.gcr.io/kube-controller-manager:v1.20.11", "Kubernetes controller manager image")
	cmd.PersistentFlags().Int32VarP(&k8s.KArmadaKubeControllerManagerReplicas, "karmada-kube-controller-manager-replicas", "", 1, "karmada kube controller manager replica set")
	cmd.PersistentFlags().StringVarP(&k8s.KArmadaControllerManagerImage, "karmada-controller-manager-image", "", "swr.ap-southeast-1.myhuaweicloud.com/karmada/karmada-controller-manager:latest", "karmada controller manager  image")
	cmd.PersistentFlags().Int32VarP(&k8s.KArmadaControllerManagerReplicas, "karmada-controller-manager-replicas", "", 1, "karmada controller manager replica set")
	cmd.PersistentFlags().StringVarP(&k8s.KArmadaWebhookImage, "karmada-webhook-image", "", "swr.ap-southeast-1.myhuaweicloud.com/karmada/karmada-webhook:latest", "karmada webhook image")
	cmd.PersistentFlags().Int32VarP(&k8s.KArmadaWebhookReplicas, "karmada-webhook-replicas", "", 1, "karmada webhook replica set")

	return cmd
}
