package service

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io"
	"io/ioutil"
	"net"
	"net/url"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

// Certificate returns the PEM encoded certificate chain for osqueryd TLS termination.
func (svc service) CertificateChain(ctx context.Context, insecure bool) ([]byte, error) {
	if svc.config.Server.TLS {
		cert, err := ioutil.ReadFile(svc.config.Server.Cert)
		if err != nil {
			return nil, errors.Wrap(err, "reading certificate file")
		}
		return cert, nil
	}

	// if kolide is not using a TLS listener itself, it must be terminated upstream.
	// we can still retrieve the certificate chain if we can establish a
	// connection to the KolideServerURL and get the cert from the connection info.
	config, err := svc.AppConfig(ctx)
	if err != nil {
		return nil, err
	}

	conn, err := connectTLS(config.KolideServerURL, insecure)
	if err != nil {
		return nil, err
	}

	return chain(conn.ConnectionState())
}

func connectTLS(serverURL string, insecure bool) (*tls.Conn, error) {
	u, err := url.Parse(serverURL)
	if err != nil {
		return nil, errors.Wrap(err, "parsing serverURL")
	}
	var hostport string
	if u.Port() == "" {
		hostport = net.JoinHostPort(u.Host, "443")
	} else {
		hostport = u.Host
	}

	conn, err := tls.Dial("tcp", hostport, &tls.Config{InsecureSkipVerify: insecure})
	if err != nil {
		return nil, errors.Wrap(err, "dial tls")
	}
	defer conn.Close()
	return conn, nil
}

func chain(cs tls.ConnectionState) ([]byte, error) {
	buf := bytes.NewBuffer([]byte(""))
	if len(cs.VerifiedChains) != 0 {
		for _, chain := range cs.VerifiedChains {
			for _, cert := range chain {
				if err := encodePEM(buf, cert); err != nil {
					return nil, errors.Wrap(err, "encode verified chains pem")
				}
			}
		}
	} else {
		for _, cert := range cs.PeerCertificates {
			if err := encodePEM(buf, cert); err != nil {
				return nil, errors.Wrap(err, "encode peer certificates pem")
			}
		}
	}
	return buf.Bytes(), nil
}

func encodePEM(buf io.Writer, cert *x509.Certificate) error {
	block := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	}
	return pem.Encode(buf, block)
}
