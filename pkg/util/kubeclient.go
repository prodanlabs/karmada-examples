package util

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	KubeQPS            = float32(5.000000)
	KubeBurst          = 10
	AcceptContentTypes = "application/json"
	ContentType        = "application/json"
)

func SetupKubeConfig(config *rest.Config) {
	config.QPS = KubeQPS
	config.Burst = KubeBurst
	config.ContentType = ContentType
	config.AcceptContentTypes = AcceptContentTypes
	config.UserAgent = rest.DefaultKubernetesUserAgent()
}

// NewClientSet ClientSet 客户端
func NewClientSet(c *rest.Config) (*kubernetes.Clientset, error) {
	SetupKubeConfig(c)

	clientSet, err := kubernetes.NewForConfig(c)
	return clientSet, err
}
