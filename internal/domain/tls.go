package domain

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

type TLSConfig struct {
	CertFile      string
	KeyFile       string
	CAFile        string
	ServerAddress string
	Server        bool
}

func (c TLSConfig) SetupTLS() (tlsConfig *tls.Config, err error) {
	tlsConfig = &tls.Config{}

	if c.CertFile != "" && c.KeyFile != "" {
		tlsConfig.Certificates = make([]tls.Certificate, 1)
		tlsConfig.Certificates[0], err = tls.LoadX509KeyPair(c.CertFile, c.KeyFile)
		if err != nil {
			return nil, err
		}
	}

	if c.CAFile != "" {
		b, err := os.ReadFile(c.CAFile)
		if err != nil {
			return nil, err
		}

		ca := x509.NewCertPool()
		if !ca.AppendCertsFromPEM(b) {
			return nil, fmt.Errorf("failed to parse root certificate %q", c.CAFile)
		}

		if c.Server {
			tlsConfig.ClientCAs = ca
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		} else {
			tlsConfig.RootCAs = ca
		}

		tlsConfig.ServerName = c.ServerAddress
	}

	return
}
