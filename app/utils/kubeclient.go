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
package utils

import (
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// RestConfig  Kubernetes kubeconfig
func RestConfig(kubeconfigPath string) (*rest.Config, error) {

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}
	config.QPS = float32(5.000000)
	config.Burst = 10
	config.ContentType = "application/json"
	config.AcceptContentTypes = "application/json"
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	return config, err
}

// NewClientSet Kubernetes ClientSet
func NewClientSet(c *rest.Config) (*kubernetes.Clientset, error) {
	return kubernetes.NewForConfig(c)
}

//NewCRDsClient clientset ClientSet
func NewCRDsClient(c *rest.Config) (*clientset.Clientset, error) {

	return clientset.NewForConfig(c)
}
