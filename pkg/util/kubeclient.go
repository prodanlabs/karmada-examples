package util

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	KubeQPS            = float32(5.000000)
	KubeBurst          = 10
	AcceptContentTypes = "application/json"
	ContentType        = "application/json"
)

// SetupKubeConfig set parameter
func SetupKubeConfig(config *rest.Config) {
	config.QPS = KubeQPS
	config.Burst = KubeBurst
	config.ContentType = ContentType
	config.AcceptContentTypes = AcceptContentTypes
	config.UserAgent = rest.DefaultKubernetesUserAgent()
}

// RestConfig  Kubernetes kubeconfig
func RestConfig(inCluster bool, kubeConfigPath string) (*rest.Config, error) {
	if inCluster {
		return rest.InClusterConfig()
	}

	return clientcmd.BuildConfigFromFlags("", kubeConfigPath)
}

// NewClientSet Kubernetes ClientSet
func NewClientSet(c *rest.Config) (*kubernetes.Clientset, error) {
	SetupKubeConfig(c)

	return kubernetes.NewForConfig(c)
}
