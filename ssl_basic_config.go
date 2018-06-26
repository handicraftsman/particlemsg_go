package particlemsg

import (
	"crypto/tls"
	"os"
)

// GetBasicSSLConfig - generates a basic *tls.Config
func GetBasicSSLConfig(crt tls.Certificate) *tls.Config {
	if os.Getenv("PMSG_UNSAFE_SSL") == "true" {
		return &tls.Config{Certificates: []tls.Certificate{crt}, InsecureSkipVerify: true}
	}
	return &tls.Config{Certificates: []tls.Certificate{crt}}
}

// GetSSLCertFromFiles - loads SSL certificate and key from their files
func GetSSLCertFromFiles(c, k string) tls.Certificate {
	crt, err := tls.LoadX509KeyPair(c, k)
	if err != nil {
		panic(err)
	}
	return crt
}
