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
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
)

const (
	// Query the IP address of the current host accessing the Internet
	getInternetIPUrl = "https://myexternalip.com/raw"
	// A split symbol that receives multiple values from a command flag
	separator = ","
)

func HomeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func PathIsExist(path string) bool {
	_, err := os.Stat(path)

	if err != nil && os.IsNotExist(err) {
		if err = os.MkdirAll(path, 0755); err != nil {
			return false
		}
	}
	return true
}

func IsIP(s string) bool {
	address := net.ParseIP(s)
	if address == nil {
		return false
	}
	return true
}

// flagsExternalIP Receive external IP from command flags
func FlagsExternalIP(externalIPs string) []net.IP {

	var ips []net.IP

	arr := strings.Split(externalIPs, separator)
	for _, v := range arr {
		ips = append(ips, StringToNetIP(v))
	}
	return ips
}

// internetIP Current host Internet IP.
func InternetIP() (net.IP, error) {

	resp, err := http.Get(getInternetIPUrl)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return StringToNetIP(string(content)), nil
}
