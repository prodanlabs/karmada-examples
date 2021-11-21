package karmada

import (
	"encoding/base64"
	"testing"

	"k8s.io/klog/v2"

	"github.com/prodanlabs/kaadm/app/utils"
)

func TestInitKArmadaResources(t *testing.T) {

	caCert, err := utils.FileToBytes("/opt/workspaces/kaadm/", "ca.crt")
	if err != nil {
		klog.Exitln("Failed to get ca cert.", err)
	}

	caBase64 := base64.StdEncoding.EncodeToString(caCert)

	if err = InitKArmadaResources("karmada-system", "/opt/workspaces/kaadm/karmada-apiserver.config", caBase64); err != nil {
		klog.Exit(err)
	}
}
