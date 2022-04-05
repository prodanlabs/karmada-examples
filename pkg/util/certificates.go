package util

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// GetClientTLSConfig returns tlsConfig based on x509 certs
// `apiserver-network-proxy` components are intended to run as standalone binaries and should not be imported as a library.
// https://github.com/kubernetes-sigs/apiserver-network-proxy/blob/master/pkg/util/certificates.go#L42-L68
func GetClientTLSConfig(caFile, certFile, keyFile, serverName string, protos []string) (*tls.Config, error) {
	certPool, err := getCACertPool(caFile)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		RootCAs:    certPool,
		MinVersion: tls.VersionTLS12,
	}
	if len(protos) != 0 {
		tlsConfig.NextProtos = protos
	}
	if certFile == "" && keyFile == "" {
		// return TLS config based on CA only
		return tlsConfig, nil
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load X509 key pair %s and %s: %v", certFile, keyFile, err)
	}

	tlsConfig.ServerName = serverName
	tlsConfig.Certificates = []tls.Certificate{cert}
	return tlsConfig, nil
}

// getCACertPool loads CA certificates to pool
// https://github.com/kubernetes-sigs/apiserver-network-proxy/blob/master/pkg/util/certificates.go#L28-L39
func getCACertPool(caFile string) (*x509.CertPool, error) {
	certPool := x509.NewCertPool()
	caCert, err := ioutil.ReadFile(filepath.Clean(caFile))
	if err != nil {
		return nil, fmt.Errorf("failed to read CA cert %s: %v", caFile, err)
	}
	ok := certPool.AppendCertsFromPEM(caCert)
	if !ok {
		return nil, fmt.Errorf("failed to append CA cert to the cert pool")
	}
	return certPool, nil
}
