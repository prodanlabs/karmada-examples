package crds

import (
	"fmt"
)

func PathClusterResourceBindings(caBundle string) string {
	return fmt.Sprintf(`apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: clusterresourcebindings.work.karmada.io
spec:
  conversion:
    strategy: Webhook
    webhook:
      clientConfig:
        url: https://karmada-webhook.karmada-system.svc:443/convert
        caBundle: "%s"
      conversionReviewVersions: ["v1"]`, caBundle)
}

func PathResourceBindings(caBundle string) string {
	return fmt.Sprintf(`apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: resourcebindings.work.karmada.io
spec:
  conversion:
    strategy: Webhook
    webhook:
      clientConfig:
        url: https://karmada-webhook.karmada-system.svc:443/convert
        caBundle: %s
      conversionReviewVersions: ["v1"]`, caBundle)
}
