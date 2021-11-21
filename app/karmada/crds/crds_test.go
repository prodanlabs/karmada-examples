package crds

import (
	"encoding/base64"
	"testing"

	"k8s.io/klog/v2"

	"github.com/prodanlabs/kaadm/app/utils"
)

func TestPathCRDs(t *testing.T) {

	restConfig, err := utils.RestConfig("./tmp/karmada-apiserver.config")
	if err != nil {
		klog.Exitln(err)
	}

	crdClient, err := utils.NewCRDsClient(restConfig)
	if err != nil {
		klog.Exitln(err)
	}

	caCert, err := utils.FileToBytes("./tmp/", "ca.crt")
	if err != nil {
		klog.Exitln("Failed to get ca cert.", err)
	}

	caBase64 := base64.StdEncoding.EncodeToString(caCert)
	pathClusterResourceBindings := PathClusterResourceBindings("", caBase64)

	if err = PatchCRDs(crdClient, "clusterresourcebindings.work.karmada.io", pathClusterResourceBindings); err != nil {
		klog.Exit(err)
	}
}
