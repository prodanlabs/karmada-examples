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
package utils

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
)

func FileToBytes(path, name string) ([]byte, error) {

	filename := filepath.Join(path, name)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats, err := file.Stat()
	if err != nil {
		return nil, err
	}

	data := make([]byte, stats.Size())

	_, err = file.Read(data)
	if err != nil {
		return nil, err
	}

	return data, nil

}

func BytesToFile(path, name string, data []byte) error {

	filename := filepath.Join(path, name)
	_, err := os.Stat(filename)
	if err == nil {
		return nil
	}

	// Create kubeconfig
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func StringToNetIP(ip string) net.IP {

	if IsIP(ip) {
		return net.ParseIP(ip)
	}
	return net.ParseIP("127.0.0.1")
}

// returns the paths for the certificate and key given the path and basename.
func PathForKey(pkiPath, name string) string {
	return filepath.Join(pkiPath, fmt.Sprintf("%s.key", name))
}

func PathForCert(pkiPath, name string) string {
	return filepath.Join(pkiPath, fmt.Sprintf("%s.crt", name))
}

// MapToString  labels to string
func MapToString(labels map[string]string) string {
	v := new(bytes.Buffer)
	for key, value := range labels {
		fmt.Fprintf(v, "%s=%s,", key, value)
	}
	return strings.TrimRight(v.String(), ",")

}

func StaticYamlToJsonByte(staticYaml string) []byte {

	jsonByte, err := yaml.YAMLToJSON([]byte(staticYaml))
	if err != nil {
		fmt.Println("Error convert string to json byte.")
		os.Exit(1)
	}
	return jsonByte
}
