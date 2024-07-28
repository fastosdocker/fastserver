package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"
)

func CreateTlsFile(tlsDir string) (error, interface{}, interface{}) {
	// 生成公私钥对
	caPrivkey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err, nil, nil
	}

	template := NewTemplate()
	rootCertDer, err := x509.CreateCertificate(rand.Reader, template, template, &caPrivkey.PublicKey, caPrivkey) //DER 格式
	if err != nil {
		return err, nil, nil
	}

	caPrivBytes, err := x509.MarshalPKCS8PrivateKey(caPrivkey)
	if err != nil {
		return err, nil, nil
	}

	keyPath := fmt.Sprintf("%s/cert.key", tlsDir)
	rootKeyFile, err := os.Create(keyPath)
	if err != nil {
		return err, nil, nil
	}

	if err = pem.Encode(rootKeyFile, &pem.Block{Type: "EC PRIVATE KEY", Bytes: caPrivBytes}); err != nil {
		return err, nil, nil
	}

	rootKeyFile.Close()
	crtPath := fmt.Sprintf(fmt.Sprintf("%s/cert.crt", tlsDir))
	rootCertFile, err := os.Create(crtPath)
	if err != nil {
		return err, nil, nil
	}

	if err = pem.Encode(rootCertFile, &pem.Block{Type: "CERTIFICATE", Bytes: rootCertDer}); err != nil {
		return err, nil, nil
	}

	rootCertFile.Close()

	return nil, keyPath, crtPath
}

func NewTemplate() *x509.Certificate {
	max := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, max)

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{ // 证书的主题信息
			Country:            []string{"CN"},             // 证书所属的国家
			Organization:       []string{"Fast Os Docker"}, // 证书存放的公司名称
			OrganizationalUnit: []string{"Fast Os Docker"}, // 证书所属的部门名称
			Province:           []string{"BeiJing"},        // 证书签发机构所在省
			CommonName:         "hello.world.com",          // 证书域名
			Locality:           []string{"BeiJing"},        // 证书签发机构所在市
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		IsCA:                  true,
		BasicConstraintsValid: true,
	}

	return template
}
