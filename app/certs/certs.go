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
package certs

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/pkg/errors"

	certutil "k8s.io/client-go/util/cert"
	"k8s.io/client-go/util/keyutil"
	"k8s.io/klog/v2"
	"k8s.io/kube-openapi/pkg/util/sets"

	"github.com/prodanlabs/kaadm/app/utils"
)

const (
	// CertificateBlockType is a possible value for pem.Block.Type.
	CertificateBlockType     = "CERTIFICATE"
	rsaKeySize               = 2048
	duration365d             = time.Hour * 24 * 365
	caCertAndKeyName         = "ca"
	etcdServerCertAndKeyName = "etcd-server"
	etcdClientCertAndKeyName = "etcd-client"
	kArmadaCertAndKeyName    = "karmada"
)

// NewPrivateKey returns a new private key.
var NewPrivateKey = GeneratePrivateKey

func GeneratePrivateKey(keyType x509.PublicKeyAlgorithm) (crypto.Signer, error) {
	if keyType == x509.ECDSA {
		return ecdsa.GenerateKey(elliptic.P256(), cryptorand.Reader)
	}

	return rsa.GenerateKey(cryptorand.Reader, rsaKeySize)
}

// CertConfig is a wrapper around certutil.Config extending it with PublicKeyAlgorithm.
type CertConfig struct {
	certutil.Config
	NotAfter           *time.Time
	PublicKeyAlgorithm x509.PublicKeyAlgorithm
}

// CertAndKeyFileName is generate certificate and key file name
type CertAndKeyFileName struct {
	CACertFileName         string
	CAKeyFileName          string
	EtcdServerCertFileName string
	EtcdServerKeyFileName  string
	EtcdClientCertFileName string
	EtcdClientKeFileName   string
	KArmadaCertFileName    string
	KArmadaKeyFileName     string
	ALLCertFileName        []string
	ALLKeyFileName         []string
}

// Config certificate information
type Config struct {
	PkiPath                     string
	Namespace                   string
	EtcdStatefulSetName         string
	EtcdServiceName             string
	EtcdReplicas                int32
	KArmadaMasterIP             string
	KArmadaApiServerServiceName string
	KArmadaWebhookServiceName   string
	FlagsExternalIP             string
}

// EncodeCertPEM returns PEM-endcoded certificate data
func EncodeCertPEM(cert *x509.Certificate) []byte {
	block := pem.Block{
		Type:  CertificateBlockType,
		Bytes: cert.Raw,
	}
	return pem.EncodeToMemory(&block)
}

// NewCertificateAuthority creates new certificate and private key for the certificate authority
func NewCertificateAuthority(config *CertConfig) (*x509.Certificate, crypto.Signer, error) {
	key, err := NewPrivateKey(config.PublicKeyAlgorithm)
	if err != nil {
		return nil, nil, errors.Errorf("unable to create private key while generating CA certificate %v", err)
	}

	cert, err := certutil.NewSelfSignedCACert(config.Config, key)
	if err != nil {
		return nil, nil, errors.Errorf("unable to create self-signed CA certificate %v", err)
	}

	return cert, key, nil
}

// NewCACertAndKey The public and private keys of the root certificate are returned
func NewCACertAndKey() (*x509.Certificate, *crypto.Signer, error) {

	certCfg := &CertConfig{Config: certutil.Config{
		CommonName:   "ca",
		Organization: []string{"karmada"},
	},
	}
	caCert, caKey, err := NewCertificateAuthority(certCfg)
	if err != nil {
		return nil, nil, errors.Errorf("failure while generating CA certificate and key: %v", err)
	}

	return caCert, &caKey, nil
}

// NewSignedCert creates a signed certificate using the given CA certificate and key
func NewSignedCert(cfg *CertConfig, key crypto.Signer, caCert *x509.Certificate, caKey crypto.Signer, isCA bool) (*x509.Certificate, error) {
	serial, err := cryptorand.Int(cryptorand.Reader, new(big.Int).SetInt64(math.MaxInt64))
	if err != nil {
		return nil, err
	}
	if len(cfg.CommonName) == 0 {
		return nil, errors.New("must specify a CommonName")
	}

	keyUsage := x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature
	if isCA {
		keyUsage |= x509.KeyUsageCertSign
	}

	RemoveDuplicateAltNames(&cfg.AltNames)

	notAfter := time.Now().Add(duration365d).UTC()
	if cfg.NotAfter != nil {
		notAfter = *cfg.NotAfter
	}

	certTmpl := x509.Certificate{
		Subject: pkix.Name{
			CommonName:   cfg.CommonName,
			Organization: cfg.Organization,
		},
		DNSNames:              cfg.AltNames.DNSNames,
		IPAddresses:           cfg.AltNames.IPs,
		SerialNumber:          serial,
		NotBefore:             caCert.NotBefore,
		NotAfter:              notAfter,
		KeyUsage:              keyUsage,
		ExtKeyUsage:           cfg.Usages,
		BasicConstraintsValid: true,
		IsCA:                  isCA,
	}
	certDERBytes, err := x509.CreateCertificate(cryptorand.Reader, &certTmpl, caCert, key.Public(), caKey)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(certDERBytes)
}

// RemoveDuplicateAltNames removes duplicate items in altNames.
func RemoveDuplicateAltNames(altNames *certutil.AltNames) {
	if altNames == nil {
		return
	}

	if altNames.DNSNames != nil {
		altNames.DNSNames = sets.NewString(altNames.DNSNames...).List()
	}

	ipsKeys := make(map[string]struct{})
	var ips []net.IP
	for _, one := range altNames.IPs {
		if _, ok := ipsKeys[one.String()]; !ok {
			ipsKeys[one.String()] = struct{}{}
			ips = append(ips, one)
		}
	}
	altNames.IPs = ips
}

// NewCertAndKey creates new certificate and key by passing the certificate authority certificate and key
func NewCertAndKey(caCert *x509.Certificate, caKey crypto.Signer, config *CertConfig) (*x509.Certificate, crypto.Signer, error) {
	if len(config.Usages) == 0 {
		return nil, nil, errors.New("must specify at least one ExtKeyUsage")
	}

	key, err := NewPrivateKey(config.PublicKeyAlgorithm)
	if err != nil {
		return nil, nil, errors.Errorf("unable to create private key %v", err)
	}

	cert, err := NewSignedCert(config, key, caCert, caKey, false)
	if err != nil {
		return nil, nil, errors.Errorf("unable to sign certificate. %v", err)
	}

	return cert, key, nil
}

// WriteCert stores the given certificate at the given location
func WriteCert(pkiPath, name string, cert *x509.Certificate) error {
	if cert == nil {
		return errors.New("certificate cannot be nil when writing to file")
	}

	certificatePath := utils.PathForCert(pkiPath, name)
	if err := certutil.WriteCert(certificatePath, EncodeCertPEM(cert)); err != nil {
		return errors.Errorf("unable to write certificate to file %v", err)
	}

	return nil
}

// WriteKey stores the given key at the given location
func WriteKey(pkiPath, name string, key crypto.Signer) error {
	if key == nil {
		return errors.New("private key cannot be nil when writing to file")
	}

	privateKeyPath := utils.PathForKey(pkiPath, name)
	encoded, err := keyutil.MarshalPrivateKeyToPEM(key)
	if err != nil {
		return errors.Errorf("unable to marshal private key to PEM %v", err)
	}
	if err := keyutil.WriteKey(privateKeyPath, encoded); err != nil {
		return errors.Errorf("unable to write private key to file %v", err)
	}

	return nil
}

// WriteCertAndKey Write certificate and key to file.
func WriteCertAndKey(pkiPath, pkiName string, ca *x509.Certificate, key *crypto.Signer) error {

	if err := WriteKey(pkiPath, pkiName, *key); err != nil {
		return err
	}

	if err := WriteCert(pkiPath, pkiName, ca); err != nil {
		return err
	}

	klog.Infof("Generate %s certificate success.", pkiName)
	return nil
}

// Create CA certificate and sign etcd karma certificate.
func (c *Config) CertificateGeneration() (*CertAndKeyFileName, error) {

	caCert, caKey, err := NewCACertAndKey()
	if err != nil {
		return nil, err
	}

	if err = WriteCertAndKey(c.PkiPath, caCertAndKeyName, caCert, caKey); err != nil {
		return nil, err
	}

	notAfter := time.Now().Add(duration365d * 10).UTC()
	etcdServerCert, etcdServerKey, err := NewCertAndKey(caCert, *caKey, c.etcdServerCertCfg(&notAfter))
	if err != nil {
		return nil, err
	}
	if err = WriteCertAndKey(c.PkiPath, etcdServerCertAndKeyName, etcdServerCert, &etcdServerKey); err != nil {
		return nil, err
	}

	etcdClientCertCfg := &CertConfig{Config: certutil.Config{
		CommonName:   "karmada-etcd-client",
		Organization: []string{"karmada"},
		Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
	},
		NotAfter: &notAfter,
	}
	etcdClientCert, etcdClientKey, err := NewCertAndKey(caCert, *caKey, etcdClientCertCfg)
	if err != nil {
		return nil, err
	}
	if err = WriteCertAndKey(c.PkiPath, etcdClientCertAndKeyName, etcdClientCert, &etcdClientKey); err != nil {
		return nil, err
	}

	kArmadaCert, kArmadaKeyerr, err := NewCertAndKey(caCert, *caKey, c.kArmadaCertCfg(&notAfter))
	if err != nil {
		return nil, err
	}

	if err = WriteCertAndKey(c.PkiPath, kArmadaCertAndKeyName, kArmadaCert, &kArmadaKeyerr); err != nil {
		return nil, err
	}

	FilesName := &CertAndKeyFileName{
		CACertFileName:         fmt.Sprintf("%s.crt", caCertAndKeyName),
		CAKeyFileName:          fmt.Sprintf("%s.key", caCertAndKeyName),
		EtcdServerCertFileName: fmt.Sprintf("%s.crt", etcdServerCertAndKeyName),
		EtcdServerKeyFileName:  fmt.Sprintf("%s.key", etcdServerCertAndKeyName),
		EtcdClientCertFileName: fmt.Sprintf("%s.crt", etcdClientCertAndKeyName),
		EtcdClientKeFileName:   fmt.Sprintf("%s.key", etcdClientCertAndKeyName),
		KArmadaCertFileName:    fmt.Sprintf("%s.crt", kArmadaCertAndKeyName),
		KArmadaKeyFileName:     fmt.Sprintf("%s.key", kArmadaCertAndKeyName),
	}
	FilesName.ALLCertFileName = append(FilesName.ALLCertFileName, fmt.Sprintf("%s.crt", caCertAndKeyName), fmt.Sprintf("%s.crt", etcdServerCertAndKeyName),
		fmt.Sprintf("%s.crt", etcdClientCertAndKeyName), fmt.Sprintf("%s.crt", kArmadaCertAndKeyName))
	FilesName.ALLKeyFileName = append(FilesName.ALLKeyFileName, fmt.Sprintf("%s.key", caCertAndKeyName), fmt.Sprintf("%s.key", etcdServerCertAndKeyName),
		fmt.Sprintf("%s.key", etcdClientCertAndKeyName), fmt.Sprintf("%s.key", kArmadaCertAndKeyName))
	return FilesName, nil
}

func (c *Config) etcdServerCertCfg(notAfter *time.Time) *CertConfig {

	var dns = []string{
		"localhost",
	}

	for i := int32(0); i < c.EtcdReplicas; i++ {
		dns = append(dns, fmt.Sprintf("%s-%v.%s.%s.svc.cluster.local", c.EtcdStatefulSetName, i, c.EtcdServiceName, c.Namespace))
	}

	return &CertConfig{
		Config: certutil.Config{
			CommonName:   "karmada-etcd-server",
			Organization: []string{"karmada"},
			Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
			AltNames: certutil.AltNames{
				IPs:      []net.IP{utils.StringToNetIP("127.0.0.1")},
				DNSNames: dns,
			},
		},
		NotAfter: notAfter,
	}
}

func (c *Config) kArmadaCertCfg(notAfter *time.Time) *CertConfig {

	var dns = []string{
		"localhost",
		"kubernetes",
		"kubernetes.default",
		"kubernetes.default.svc",
		c.KArmadaApiServerServiceName,
		c.KArmadaWebhookServiceName,
		fmt.Sprintf("%s.%s.svc.cluster.local", c.KArmadaApiServerServiceName, c.Namespace),
		fmt.Sprintf("%s.%s.svc.cluster.local", c.KArmadaWebhookServiceName, c.Namespace),
		fmt.Sprintf("*.%s.svc.cluster.local", c.Namespace),
		fmt.Sprintf("*.%s.svc", c.Namespace),
	}

	if hostName, err := os.Hostname(); err != nil {
		klog.Errorf("%v, Failed to get the current hostname.", err)
	} else {
		dns = append(dns, hostName)
	}

	ips := utils.FlagsExternalIP(c.FlagsExternalIP)

	internetIP, err := utils.InternetIP()
	if err != nil {
		klog.Errorf("%v, Failed to obtain internet IP. ", err)
	} else {
		ips = append(ips, internetIP)
	}

	ips = append(
		ips,
		utils.StringToNetIP("127.0.0.1"),
		utils.StringToNetIP("10.254.0.1"),
		utils.StringToNetIP(c.KArmadaMasterIP),
	)

	return &CertConfig{Config: certutil.Config{
		CommonName:   "system:admin",
		Organization: []string{"system:masters"},
		Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		AltNames: certutil.AltNames{
			IPs:      ips,
			DNSNames: dns,
		},
	},
		NotAfter: notAfter,
	}
}
