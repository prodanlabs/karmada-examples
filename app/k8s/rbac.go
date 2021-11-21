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

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

func (i *InstallOptions) ClusterRoleFromSpec(name string, rules []rbacv1.PolicyRule) *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1",
			Kind:       "ClusterRole",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: i.Namespace,
		},
		Rules: rules,
	}
}

func (i *InstallOptions) ClusterRoleBindingFromSpec(clusterRoleBindingName, clusterRoleName, saName string) *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1",
			Kind:       "ClusterRoleBinding",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterRoleBindingName,
			Namespace: i.Namespace,
		},

		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     clusterRoleName,
		},

		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      saName,
				Namespace: i.Namespace,
			},
		},
	}
}

func (i *InstallOptions) CreateClusterRole(clusterRole *rbacv1.ClusterRole) error {

	clusterRoleClient := i.KubeClientSet.RbacV1().ClusterRoles()

	clusterRoleList, err := clusterRoleClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, v := range clusterRoleList.Items {
		if clusterRole.Name == v.Name {
			klog.Warningf("ClusterRole %s already exists.", clusterRole.Name)
			return nil
		}
	}

	_, err = clusterRoleClient.Create(context.TODO(), clusterRole, metav1.CreateOptions{})
	if err != nil {
		return errors.Errorf("Create ClusterRole %s failed: %v\n", clusterRole.Name, err)
	}
	return nil
}

func (i *InstallOptions) CreateClusterRoleBinding(clusterRole *rbacv1.ClusterRoleBinding) error {

	crbClient := i.KubeClientSet.RbacV1().ClusterRoleBindings()

	crbList, err := crbClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, v := range crbList.Items {
		if clusterRole.Name == v.Name {
			klog.Infof("CreateClusterRoleBinding %s already exists.", clusterRole.Name)
			return nil
		}
	}

	_, err = crbClient.Create(context.TODO(), clusterRole, metav1.CreateOptions{})
	if err != nil {
		return errors.Errorf("Create CreateClusterRoleBinding %s failed: %v\n", clusterRole.Name, err)
	}
	return nil
}
