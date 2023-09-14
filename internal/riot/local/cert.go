package local

import (
	"crypto/x509"
	_ "embed"
)

//go:embed riot.pem
var riotCertificate []byte

var caCertPool *x509.CertPool

func init() {
	caCertPool = x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(riotCertificate)
}
