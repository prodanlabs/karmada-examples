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
	"encoding/json"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func getNodeName(nodes *corev1.NodeList, nodeIP string) string {

	for _, v := range nodes.Items {
		for _, ip := range v.Status.Addresses {
			if nodeIP == ip.Address {
				return v.GetName()
			}
		}
	}
	return ""
}

func (i *InstallOptions) AddNodeSelectorLabels() error {
	nodes, err := i.KubeClientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, v := range i.MasterIP {
		nodeName := getNodeName(nodes, v.String())
		node, err := i.KubeClientSet.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
		if err != nil {
			return err
		}
		NodeSelectorLabels = map[string]string{"karmada.io/master": ""}
		labels := node.Labels
		labels["karmada.io/master"] = ""
		patchData := map[string]interface{}{"metadata": map[string]map[string]string{"labels": labels}}

		playLoadBytes, _ := json.Marshal(patchData)

		if _, err = i.KubeClientSet.CoreV1().Nodes().Patch(context.TODO(), nodeName, types.StrategicMergePatchType, playLoadBytes, metav1.PatchOptions{}); err != nil {
			return err
		}
	}

	return nil

}
