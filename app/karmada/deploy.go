package karmada

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"

	"github.com/prodanlabs/kaadm/app/karmada/crds"
	"github.com/prodanlabs/kaadm/app/utils"
)

func InitKArmadaResources(namespace, kubeconfig, caBase64 string) error {
	restConfig, err := utils.RestConfig(kubeconfig)
	if err != nil {
		return err
	}

	clientSet, err := utils.NewClientSet(restConfig)
	if err != nil {
		return err
	}
	clientSet.CoreV1().Namespaces().Create(context.TODO(), &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}, metav1.CreateOptions{})

	// create webhook configuration
	// https://github.com/karmada-io/karmada/blob/master/artifacts/deploy/webhook-configuration.yaml

	klog.Info("Crate MutatingWebhookConfiguration mutating-config.")

	if err = createMutatingWebhookConfiguration(clientSet, mutatingConfig(namespace, caBase64)); err != nil {
		klog.Exitln(err)
	}
	klog.Info("Crate ValidatingWebhookConfiguration validating-config.")
	if err = createValidatingWebhookConfiguration(clientSet, validatingConfig(namespace, caBase64)); err != nil {
		klog.Exitln(err)
	}

	crdClient, err := utils.NewCRDsClient(restConfig)
	if err != nil {
		return err
	}

	crdSlice := []string{
		crds.Config,
		crds.ClustersCluster,
		crds.ServiceExports,
		crds.ServiceImports,
		crds.ClusterOverridePolicies,
		crds.ClusterPropagationPolicies,
		crds.OverridePolicies,
		crds.PropagationPolicies,
		crds.ReplicaSchedulingPolicies,
		crds.ClusterResourceBindings,
		crds.ResourceBindings,
		crds.Works,
	}

	klog.Info("Create karmada CustomResourceDefinition.")
	for _, v := range crdSlice {
		if err = crds.CreateCRDs(crdClient, v); err != nil {
			return err
		}
	}

	pathClusterResourceBindings := crds.PathClusterResourceBindings(caBase64)
	pathResourceBindings := crds.PathResourceBindings(caBase64)

	klog.Info("Patch CustomResourceDefinition clusterresourcebindings.work.karmada.io.")
	if err = crds.PatchCRDs(crdClient, "clusterresourcebindings.work.karmada.io", pathClusterResourceBindings); err != nil {
		return err
	}
	klog.Info("Patch CustomResourceDefinition resourcebindings.work.karmada.io.")
	if err = crds.PatchCRDs(crdClient, "resourcebindings.work.karmada.io", pathResourceBindings); err != nil {
		return err
	}

	return nil
}
