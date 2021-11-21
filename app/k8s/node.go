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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

func AddNodeSelectorLabels(c *kubernetes.Clientset) error {
	nodes, err := c.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	var nodeName string

	for _, v := range nodes.Items {
		for _, ip := range v.Status.Addresses {
			if KArmadaMasterIP == ip.Address {
				nodeName = v.GetName()
			}
		}
	}

	node, err := c.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	NodeSelectorLabels = map[string]string{"karmada.io/master": nodeName}

	labels := node.Labels
	labels["karmada.io/master"] = nodeName
	patchData := map[string]interface{}{"metadata": map[string]map[string]string{"labels": labels}}

	playLoadBytes, _ := json.Marshal(patchData)

	if _, err = c.CoreV1().Nodes().Patch(context.TODO(), nodeName, types.StrategicMergePatchType, playLoadBytes, metav1.PatchOptions{}); err != nil {
		return err
	}
	return nil
}
