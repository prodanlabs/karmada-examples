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
package crds

import (
	"context"
	"encoding/json"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"

	"github.com/prodanlabs/kaadm/app/utils"
)

func CreateCRDs(crdClient *clientset.Clientset, staticYaml string) error {

	obj := apiextensionsv1.CustomResourceDefinition{}

	if err := json.Unmarshal(utils.StaticYamlToJsonByte(staticYaml), &obj); err != nil {
		klog.Errorln("Error convert json byte to apiExtensionsV1 CustomResourceDefinition struct.")
		return err
	}

	crd := crdClient.ApiextensionsV1().CustomResourceDefinitions()
	if _, err := crd.Create(context.TODO(), &obj, metav1.CreateOptions{}); err != nil {
		return err
	}
	return nil
}

func PatchCRDs(crdClient *clientset.Clientset, name, staticYaml string) error {

	crd := crdClient.ApiextensionsV1().CustomResourceDefinitions()
	if _, err := crd.Patch(context.TODO(), name, types.StrategicMergePatchType, utils.StaticYamlToJsonByte(staticYaml), metav1.PatchOptions{}); err != nil {
		return err
	}
	return nil
}
