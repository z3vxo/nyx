package teamserver

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"
)

func GenSelSigned(host string) (tls.Certificate, error) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: host},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		DNSNames:     []string{host},
	}

	cert, _ := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)

	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key)})

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert})

	return tls.X509KeyPair(certPEM, keyPEM)
}

// func NewLetsEncrypt(domain string) *autocert.Manager {

// 	return &autocert.Manager{
// 		Cache:      autocert.DirCache("/var/certs/cache"),
// 		Prompt:     autocert.AcceptTOS,
// 		HostPolicy: autocert.HostWhitelist(domain),
// 		Email:      fmt.Sprintf("admin@%s", domain),
// 	}
// }
