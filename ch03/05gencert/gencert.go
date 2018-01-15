package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)

const shiftWidth = 128

func main() {
	max := new(big.Int).Lsh(big.NewInt(1), shiftWidth)
	serialNumber, _ := rand.Int(rand.Reader, max)
	template := createTemplate(serialNumber)

	pk, _ := rsa.GenerateKey(rand.Reader, 2048)
	if err := createCertDotPem("cert.pem", template, pk); err != nil {
		os.Exit(1)
	}
	if err := createKeyDotPem("key.pem", pk); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

func createTemplate(serialNumber *big.Int) x509.Certificate {
	const (
		year    = 365
		day     = 24
		address = "127.0.0.1"
	)
	return x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization:       []string{"manning Publications Co."},
			OrganizationalUnit: []string{"Books"},
			CommonName:         "Go Web Programming",
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(year * day * time.Hour),
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses: []net.IP{net.ParseIP(address)},
	}
}

func createCertDotPem(filepath string, template x509.Certificate, pk *rsa.PrivateKey) error {
	certOut, err1 := os.Create(filepath)
	if err1 != nil {
		fmt.Printf("cannot open: %s\n", filepath)
		return err1
	}
	defer certOut.Close()

	derBytes, err2 := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk)
	if err2 != nil {
		return err2
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	return nil
}

func createKeyDotPem(filepath string, pk *rsa.PrivateKey) error {
	keyOut, err := os.Create(filepath)
	if err != nil {
		fmt.Printf("cannot open: %s\n", filepath)
		return err
	}
	defer keyOut.Close()

	pem.Encode(keyOut, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	return nil
}
