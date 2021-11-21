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

	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

func (i *InstallOptions) ServiceAccountFromSpec(name []string) *[]corev1.ServiceAccount {

	var sa []corev1.ServiceAccount

	for _, v := range name {

		sa = append(sa, corev1.ServiceAccount{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "V1",
				Kind:       "ServiceAccount",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      v,
				Namespace: i.Namespace,
			},
		})

	}

	return &sa
}

func (i *InstallOptions) CreateServiceAccount(sa *[]corev1.ServiceAccount) error {
	saClient := i.KubeClientSet.CoreV1().ServiceAccounts(i.Namespace)

	for _, v := range *sa {

		if _, err := saClient.Get(context.TODO(), v.Name, metav1.GetOptions{}); err == nil {
			klog.Warningf("ServiceAccount %s already exists. ", v.Name)
			continue
		}

		if _, err := saClient.Create(context.TODO(), &v, metav1.CreateOptions{}); err != nil {
			return errors.Errorf("Create secret %s failed: %v\n", v.Name, err)
		}

	}

	return nil

}
