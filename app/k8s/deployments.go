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
	"fmt"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	//	//karmada apiserver
	kArmadaDeploymentAPIVersion       = "apps/v1"
	kArmadaDeploymentKind             = "Deployment"
	kArmadaApiServerDeploymentName    = "karmada-apiserver"
	kArmadaApiServerServiceName       = "karmada-apiserver"
	kArmadaApiServerContainerPortName = "server"
	kArmadaApiServerContainerPort     = 5443
	serviceClusterIP                  = "10.96.0.0/12"
	kubeConfigSecretName              = "kubeconfig"
	kArmadaCertSecretName             = "karmada-cert"
	//karmada Kube Controller Manager
	kArmadaControllerManagerClusterRoleName    = "karmada-controller-manager"
	kArmadaKubeControllerManagerDeploymentName = "karmada-kube-controller-manager"
	kArmadaKubeControllerManagerPortName       = "server"
	kArmadaKubeControllerManagerServiceName    = "kube-controller-manager"
	kArmadaKubeControllerManagerPort           = 10257
	kubeConfigVolumeMountName                  = "kubeconfig"
	kubeConfigContainerMountPath               = "/etc/kubeconfig"
	kubeConfigVolumeReadOnly                   = true
	//karmada-scheduler
	KArmadaSchedulerDeploymentName     = "karmada-scheduler"
	KArmadaSchedulerServiceAccountName = "karmada-scheduler"
	//karmada Controller Manager
	kArmadaControllerManagerDeploymentName = "karmada-controller-manager"
	kArmadaControllerManagerSecurePort     = 10357
	kArmadaControllerManagerPortName       = "server"
	kArmadaControllerManagerServiceName    = "karmada-controller-manager"
	//	webhook
	kArmadaWebhookDeploymentName      = "karmada-webhook"
	kArmadaWebhookServiceAccountName  = "karmada-webhook"
	kArmadaWebhookServiceName         = "karmada-webhook"
	kArmadaWebhookCertSecretName      = "karmada-webhook-cert"
	kArmadaWebhookCertVolumeMountPath = "/var/serving-cert"
	kArmadaWebhookCertVolumeReadOnly  = true
	kArmadaWebhookPortName            = "webhook"
	kArmadaWebhookTargetPort          = 8443
	kArmadaWebhookPort                = 443
)

var (
	kArmadaApiServerLabels             = map[string]string{"app": kArmadaApiServerDeploymentName}
	kArmadaKubeControllerManagerLabels = map[string]string{"app": kArmadaKubeControllerManagerDeploymentName}
	kArmadaSchedulerLabels             = map[string]string{"app": KArmadaSchedulerDeploymentName}
	kArmadaControllerManagerLabels     = map[string]string{"app": kArmadaControllerManagerDeploymentName}
	kArmadaWebhookLabels               = map[string]string{"app": kArmadaWebhookDeploymentName}
)

func (i *InstallOptions) kArmadaApiServerContainerCommand() []string {

	etcdClusterConfig := ""
	for v := int32(0); v < EtcdReplicas; v++ {
		etcdClusterConfig += fmt.Sprintf("https://%s-%v.%s.%s.svc.cluster.local:%v", etcdStatefulSetName, v, etcdServiceName, i.Namespace, etcdContainerClientPort) + ","
	}

	command := []string{
		"kube-apiserver",
		"--allow-privileged=true",
		"--authorization-mode=Node,RBAC",
		fmt.Sprintf("--client-ca-file=%s/%s", certsVolumeMountPath, i.CertAndKeyFileName.CACertFileName),
		"--enable-admission-plugins=NodeRestriction",
		"--enable-bootstrap-token-auth=true",
		fmt.Sprintf("--etcd-cafile=%s/%s", certsVolumeMountPath, i.CertAndKeyFileName.CACertFileName),
		fmt.Sprintf("--etcd-certfile=%s/%s", certsVolumeMountPath, i.CertAndKeyFileName.EtcdClientCertFileName),
		fmt.Sprintf("--etcd-keyfile=%s/%s", certsVolumeMountPath, i.CertAndKeyFileName.EtcdClientKeFileName),
		fmt.Sprintf("--etcd-servers=%s", strings.TrimRight(etcdClusterConfig, ",")),
		"--bind-address=0.0.0.0",
		"--insecure-port=0",
		fmt.Sprintf("--kubelet-client-certificate=%s/%s", certsVolumeMountPath, i.CertAndKeyFileName.KArmadaCertFileName),
		fmt.Sprintf("--kubelet-client-key=%s/%s", certsVolumeMountPath, i.CertAndKeyFileName.KArmadaKeyFileName),
		"--kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname",
		"--disable-admission-plugins=StorageObjectInUseProtection,ServiceAccount",
		"--runtime-config=",
		fmt.Sprintf("--apiserver-count=%v", KArmadaApiServerReplicas),
		fmt.Sprintf("--secure-port=%v", kArmadaApiServerContainerPort),
		"--service-account-issuer=https://kubernetes.default.svc.cluster.local",
		fmt.Sprintf("--service-account-key-file=%s/%s", certsVolumeMountPath, i.CertAndKeyFileName.KArmadaKeyFileName),
		fmt.Sprintf("--service-account-signing-key-file=%s/%s", certsVolumeMountPath, i.CertAndKeyFileName.KArmadaKeyFileName),
		fmt.Sprintf("--service-cluster-ip-range=%s", serviceClusterIP),
		fmt.Sprintf("--proxy-client-cert-file=%s/%s", certsVolumeMountPath, i.CertAndKeyFileName.KArmadaCertFileName),
		fmt.Sprintf("--proxy-client-key-file=%s/%s", certsVolumeMountPath, i.CertAndKeyFileName.KArmadaKeyFileName),
		"--requestheader-allowed-names=front-proxy-client",
		fmt.Sprintf("--requestheader-client-ca-file=%s/%s", certsVolumeMountPath, i.CertAndKeyFileName.KArmadaCertFileName),
		"--requestheader-extra-headers-prefix=X-Remote-Extra-",
		"--requestheader-group-headers=X-Remote-Group",
		"--requestheader-username-headers=X-Remote-User",
		fmt.Sprintf("--tls-cert-file=%s/%s", certsVolumeMountPath, i.CertAndKeyFileName.KArmadaCertFileName),
		fmt.Sprintf("--tls-private-key-file=%s/%s", certsVolumeMountPath, i.CertAndKeyFileName.KArmadaKeyFileName),
	}

	return command

}

func (i *InstallOptions) makeKArmadaApiServerDeployment() *appsv1.Deployment {

	apiServer := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: kArmadaDeploymentAPIVersion,
			Kind:       kArmadaDeploymentKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      kArmadaApiServerDeploymentName,
			Namespace: i.Namespace,
			Labels:    kArmadaApiServerLabels,
		},
	}

	// Probes
	livenesProbe := &corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/livez",
				Port: intstr.IntOrString{
					IntVal: kArmadaApiServerContainerPort,
				},
				Scheme: corev1.URISchemeHTTPS,
			},
		},
		InitialDelaySeconds: 15,
		FailureThreshold:    3,
		PeriodSeconds:       30,
		TimeoutSeconds:      5,
	}
	readinesProbe := &corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/readyz",
				Port: intstr.IntOrString{
					IntVal: kArmadaApiServerContainerPort,
				},
				Scheme: corev1.URISchemeHTTPS,
			},
		},
		FailureThreshold: 3,
		PeriodSeconds:    30,
		TimeoutSeconds:   5,
	}

	podSpec := corev1.PodSpec{
		Affinity: &corev1.Affinity{
			PodAntiAffinity: &corev1.PodAntiAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
					{
						TopologyKey: "kubernetes.io/hostname",
						LabelSelector: &metav1.LabelSelector{
							MatchExpressions: []metav1.LabelSelectorRequirement{
								{
									Key:      "app",
									Operator: metav1.LabelSelectorOpIn,
									Values:   []string{kArmadaApiServerDeploymentName},
								},
							},
						},
					},
				},
			},
		},
		Containers: []corev1.Container{
			{
				Name:    kArmadaApiServerDeploymentName,
				Image:   KArmadaApiServerImage,
				Command: i.kArmadaApiServerContainerCommand(),
				Ports: []corev1.ContainerPort{
					{
						Name:          kArmadaApiServerContainerPortName,
						ContainerPort: kArmadaApiServerContainerPort,
						Protocol:      corev1.ProtocolTCP,
						HostPort:      kArmadaApiServerContainerPort,
					},
				},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      kArmadaCertSecretName,
						ReadOnly:  certsVolumeReadOnly,
						MountPath: certsVolumeMountPath,
					},
				},
				LivenessProbe:  livenesProbe,
				ReadinessProbe: readinesProbe,
			},
		},
		Volumes: []corev1.Volume{
			{
				Name: kArmadaCertSecretName,
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: kArmadaCertSecretName,
					},
				},
			},
		},
		NodeSelector: NodeSelectorLabels,
		//HostNetwork:  true,
		Tolerations: []corev1.Toleration{
			{
				Effect:   corev1.TaintEffectNoExecute,
				Operator: corev1.TolerationOpExists,
			},
		},
	}

	// PodTemplateSpec
	podTemplateSpec := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name:      kArmadaApiServerDeploymentName,
			Namespace: i.Namespace,
			Labels:    kArmadaApiServerLabels,
		},
		Spec: podSpec,
	}

	// DeploymentSpec
	apiServer.Spec = appsv1.DeploymentSpec{
		Replicas: &KArmadaApiServerReplicas,
		Template: podTemplateSpec,
		Selector: &metav1.LabelSelector{
			MatchLabels: kArmadaApiServerLabels,
		},
	}
	return apiServer
}

func (i *InstallOptions) makeKArmadaKubeControllerManagerDeployment() *appsv1.Deployment {

	kArmadaKubeControllerManager := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: kArmadaDeploymentAPIVersion,
			Kind:       kArmadaDeploymentKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      kArmadaKubeControllerManagerDeploymentName,
			Namespace: i.Namespace,
			Labels:    kArmadaKubeControllerManagerLabels,
		},
	}

	podSpec := corev1.PodSpec{
		Affinity: &corev1.Affinity{
			PodAntiAffinity: &corev1.PodAntiAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
					{
						TopologyKey: "kubernetes.io/hostname",
						LabelSelector: &metav1.LabelSelector{
							MatchExpressions: []metav1.LabelSelectorRequirement{
								{
									Key:      "app",
									Operator: metav1.LabelSelectorOpIn,
									Values:   []string{kArmadaKubeControllerManagerDeploymentName},
								},
							},
						},
					},
				},
			},
		},
		Containers: []corev1.Container{
			{
				Name:  kArmadaKubeControllerManagerDeploymentName,
				Image: KArmadaKubeControllerManagerImage,
				Command: []string{
					"kube-controller-manager",
					"--allocate-node-cidrs=true",
					"--authentication-kubeconfig=/etc/kubeconfig",
					"--authorization-kubeconfig=/etc/kubeconfig",
					"--bind-address=0.0.0.0",
					fmt.Sprintf("--client-ca-file=%s/%s", certsVolumeMountPath, i.CertAndKeyFileName.CACertFileName),
					"--cluster-cidr=10.244.0.0/16",
					fmt.Sprintf("--cluster-name=%s", ClusterName),
					fmt.Sprintf("--cluster-signing-cert-file=%s/%s", certsVolumeMountPath, i.CertAndKeyFileName.CACertFileName),
					fmt.Sprintf("--cluster-signing-key-file=%s/%s", certsVolumeMountPath, i.CertAndKeyFileName.CAKeyFileName),
					"--controllers=namespace,garbagecollector,serviceaccount-token",
					"--kubeconfig=/etc/kubeconfig",
					"--leader-elect=true",
					fmt.Sprintf("--leader-elect-resource-namespace=%s", i.Namespace),
					"--node-cidr-mask-size=24",
					"--port=0",
					fmt.Sprintf("--root-ca-file=%s/%s", certsVolumeMountPath, i.CertAndKeyFileName.CACertFileName),
					fmt.Sprintf("--service-account-private-key-file=%s/%s", certsVolumeMountPath, i.CertAndKeyFileName.KArmadaKeyFileName),
					fmt.Sprintf("--service-cluster-ip-range=%s", serviceClusterIP),
					"--use-service-account-credentials=true",
					"--v=4",
				},
				Ports: []corev1.ContainerPort{
					{
						Name:          kArmadaKubeControllerManagerPortName,
						ContainerPort: kArmadaKubeControllerManagerPort,
						Protocol:      corev1.ProtocolTCP,
					},
				},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      kubeConfigVolumeMountName,
						ReadOnly:  kubeConfigVolumeReadOnly,
						MountPath: kubeConfigContainerMountPath,
						SubPath:   kubeConfigVolumeMountName,
					},
					{
						Name:      kArmadaCertSecretName,
						ReadOnly:  certsVolumeReadOnly,
						MountPath: certsVolumeMountPath,
					},
				},
			},
		},
		Volumes: []corev1.Volume{
			{
				Name: kubeConfigVolumeMountName,
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: kubeConfigSecretName,
					},
				},
			},
			{
				Name: kArmadaCertSecretName,
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: kArmadaCertSecretName,
					},
				},
			},
		},

		Tolerations: []corev1.Toleration{
			{
				Effect:   corev1.TaintEffectNoExecute,
				Operator: corev1.TolerationOpExists,
			},
		},
	}
	// PodTemplateSpec
	podTemplateSpec := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name:      kArmadaKubeControllerManagerDeploymentName,
			Namespace: i.Namespace,
			Labels:    kArmadaKubeControllerManagerLabels,
		},
		Spec: podSpec,
	}
	// DeploymentSpec
	kArmadaKubeControllerManager.Spec = appsv1.DeploymentSpec{
		Replicas: &KArmadaKubeControllerManagerReplicas,
		Template: podTemplateSpec,
		Selector: &metav1.LabelSelector{
			MatchLabels: kArmadaKubeControllerManagerLabels,
		},
	}

	return kArmadaKubeControllerManager
}

func (i *InstallOptions) makeKArmadaSchedulerDeployment() *appsv1.Deployment {

	scheduler := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: kArmadaDeploymentAPIVersion,
			Kind:       kArmadaDeploymentKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      KArmadaSchedulerDeploymentName,
			Namespace: i.Namespace,
			Labels:    kArmadaSchedulerLabels,
		},
	}

	podSpec := corev1.PodSpec{
		Affinity: &corev1.Affinity{
			PodAntiAffinity: &corev1.PodAntiAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
					{
						TopologyKey: "kubernetes.io/hostname",
						LabelSelector: &metav1.LabelSelector{
							MatchExpressions: []metav1.LabelSelectorRequirement{
								{
									Key:      "app",
									Operator: metav1.LabelSelectorOpIn,
									Values:   []string{KArmadaSchedulerDeploymentName},
								},
							},
						},
					},
				},
			},
		},
		Containers: []corev1.Container{
			{
				Name:  KArmadaSchedulerDeploymentName,
				Image: KArmadaSchedulerImage,
				Command: []string{
					"/bin/karmada-scheduler",
					"--kubeconfig=/etc/kubeconfig",
					"--bind-address=0.0.0.0",
					"--secure-port=10351",
					"--feature-gates=Failover=true",
					"--enable-scheduler-estimator=true",
					"--leader-elect=true",
					fmt.Sprintf("--leader-elect-resource-namespace=%s", i.Namespace),
					"--v=4",
				},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      kubeConfigVolumeMountName,
						ReadOnly:  kubeConfigVolumeReadOnly,
						MountPath: kubeConfigContainerMountPath,
						SubPath:   kubeConfigVolumeMountName,
					},
				},
			},
		},
		Volumes: []corev1.Volume{
			{
				Name: kubeConfigVolumeMountName,
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: kubeConfigSecretName,
					},
				},
			},
		},

		Tolerations: []corev1.Toleration{
			{
				Effect:   corev1.TaintEffectNoExecute,
				Operator: corev1.TolerationOpExists,
			},
		},
		ServiceAccountName: KArmadaSchedulerServiceAccountName,
	}

	// PodTemplateSpec
	podTemplateSpec := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name:      KArmadaSchedulerDeploymentName,
			Namespace: i.Namespace,
			Labels:    kArmadaSchedulerLabels,
		},
		Spec: podSpec,
	}

	// DeploymentSpec
	scheduler.Spec = appsv1.DeploymentSpec{
		Replicas: &KArmadaSchedulerReplicas,
		Template: podTemplateSpec,
		Selector: &metav1.LabelSelector{
			MatchLabels: kArmadaSchedulerLabels,
		},
	}

	return scheduler
}

func (i *InstallOptions) makeKArmadaControllerManagerDeployment() *appsv1.Deployment {
	karmadaControllerManager := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: kArmadaDeploymentAPIVersion,
			Kind:       kArmadaDeploymentKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      kArmadaControllerManagerDeploymentName,
			Namespace: i.Namespace,
			Labels:    kArmadaControllerManagerLabels,
		},
	}

	podSpec := corev1.PodSpec{
		Affinity: &corev1.Affinity{
			PodAntiAffinity: &corev1.PodAntiAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
					{
						TopologyKey: "kubernetes.io/hostname",
						LabelSelector: &metav1.LabelSelector{
							MatchExpressions: []metav1.LabelSelectorRequirement{
								{
									Key:      "app",
									Operator: metav1.LabelSelectorOpIn,
									Values:   []string{kArmadaControllerManagerDeploymentName},
								},
							},
						},
					},
				},
			},
		},
		Containers: []corev1.Container{
			{
				Name:  kArmadaControllerManagerDeploymentName,
				Image: KArmadaControllerManagerImage,
				Command: []string{
					"/bin/karmada-controller-manager",
					"--kubeconfig=/etc/kubeconfig",
					"--bind-address=0.0.0.0",
					"--cluster-status-update-frequency=10s",
					"--secure-port=10357",
					"--v=4",
				},
				Ports: []corev1.ContainerPort{
					{
						Name:          kArmadaControllerManagerPortName,
						ContainerPort: kArmadaControllerManagerSecurePort,
						Protocol:      corev1.ProtocolTCP,
					},
				},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      kubeConfigVolumeMountName,
						ReadOnly:  kubeConfigVolumeReadOnly,
						MountPath: kubeConfigContainerMountPath,
						SubPath:   kubeConfigVolumeMountName,
					},
				},
			},
		},
		Volumes: []corev1.Volume{
			{
				Name: kubeConfigVolumeMountName,
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: kubeConfigSecretName,
					},
				},
			},
		},

		Tolerations: []corev1.Toleration{
			{
				Effect:   corev1.TaintEffectNoExecute,
				Operator: corev1.TolerationOpExists,
			},
		},
		ServiceAccountName: kArmadaControllerManagerServiceName,
	}

	// PodTemplateSpec
	podTemplateSpec := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name:      kArmadaControllerManagerDeploymentName,
			Namespace: i.Namespace,
			Labels:    kArmadaControllerManagerLabels,
		},
		Spec: podSpec,
	}
	// DeploymentSpec
	karmadaControllerManager.Spec = appsv1.DeploymentSpec{
		Replicas: &KArmadaControllerManagerReplicas,
		Template: podTemplateSpec,
		Selector: &metav1.LabelSelector{
			MatchLabels: kArmadaControllerManagerLabels,
		},
	}

	return karmadaControllerManager
}

func (i *InstallOptions) makeKArmadaWebhookDeployment() *appsv1.Deployment {

	webhook := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: kArmadaDeploymentAPIVersion,
			Kind:       kArmadaDeploymentKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      kArmadaWebhookDeploymentName,
			Namespace: i.Namespace,
			Labels:    kArmadaWebhookLabels,
		},
	}

	podSpec := corev1.PodSpec{
		Affinity: &corev1.Affinity{
			PodAntiAffinity: &corev1.PodAntiAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
					{
						TopologyKey: "kubernetes.io/hostname",
						LabelSelector: &metav1.LabelSelector{
							MatchExpressions: []metav1.LabelSelectorRequirement{
								{
									Key:      "app",
									Operator: metav1.LabelSelectorOpIn,
									Values:   []string{kArmadaWebhookDeploymentName},
								},
							},
						},
					},
				},
			},
		},
		Containers: []corev1.Container{
			{
				Name:  kArmadaWebhookDeploymentName,
				Image: KArmadaWebhookImage,
				Command: []string{
					"/bin/karmada-webhook",
					"--kubeconfig=/etc/kubeconfig",
					"--bind-address=0.0.0.0",
					fmt.Sprintf("--secure-port=%v", kArmadaWebhookTargetPort),
					fmt.Sprintf("--cert-dir=%s", kArmadaWebhookCertVolumeMountPath),
					"--v=4",
				},
				Ports: []corev1.ContainerPort{
					{
						Name:          kArmadaWebhookPortName,
						ContainerPort: kArmadaWebhookTargetPort,
						Protocol:      corev1.ProtocolTCP,
					},
				},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      kubeConfigVolumeMountName,
						ReadOnly:  kubeConfigVolumeReadOnly,
						MountPath: kubeConfigContainerMountPath,
						SubPath:   kubeConfigVolumeMountName,
					},
					{
						Name:      kArmadaWebhookCertSecretName,
						ReadOnly:  kArmadaWebhookCertVolumeReadOnly,
						MountPath: kArmadaWebhookCertVolumeMountPath,
					},
				},
				ReadinessProbe: &corev1.Probe{
					Handler: corev1.Handler{
						HTTPGet: &corev1.HTTPGetAction{
							Path: "/readyz",
							Port: intstr.IntOrString{
								IntVal: kArmadaWebhookTargetPort,
							},
							Scheme: corev1.URISchemeHTTPS,
						},
					},
				},
			},
		},
		Volumes: []corev1.Volume{
			{
				Name: kubeConfigVolumeMountName,
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: kubeConfigSecretName,
					},
				},
			},
			{
				Name: kArmadaWebhookCertSecretName,
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: kArmadaWebhookCertSecretName,
					},
				},
			},
		},

		Tolerations: []corev1.Toleration{
			{
				Effect:   corev1.TaintEffectNoExecute,
				Operator: corev1.TolerationOpExists,
			},
		},
	}

	// PodTemplateSpec
	podTemplateSpec := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name:      kArmadaWebhookDeploymentName,
			Namespace: i.Namespace,
			Labels:    kArmadaWebhookLabels,
		},
		Spec: podSpec,
	}
	// DeploymentSpec
	webhook.Spec = appsv1.DeploymentSpec{
		Replicas: &KArmadaWebhookReplicas,
		Template: podTemplateSpec,
		Selector: &metav1.LabelSelector{
			MatchLabels: kArmadaWebhookLabels,
		},
	}

	return webhook
}
