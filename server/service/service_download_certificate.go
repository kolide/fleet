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

	u, err := url.Parse(config.KolideServerURL)
	if err != nil {
		return nil, errors.Wrap(err, "parsing serverURL")
	}

	conn, err := connectTLS(u, insecure)
	if err != nil {
		return nil, err
	}

	return chain(conn.ConnectionState(), u.Hostname())
}

func connectTLS(serverURL *url.URL, insecure bool) (*tls.Conn, error) {
	var hostport string
	if serverURL.Port() == "" {
		hostport = net.JoinHostPort(serverURL.Host, "443")
	} else {
		hostport = serverURL.Host
	}

	conn, err := tls.Dial("tcp", hostport, &tls.Config{InsecureSkipVerify: insecure})
	if err != nil {
		return nil, errors.Wrap(err, "dial tls")
	}
	defer conn.Close()
	return conn, nil
}

// chain builds a PEM encoded certificate chain using the PeerCertificates
// in tls.ConnectionState. chain uses the hostname to omit the Leaf certificate
// from the chain.
func chain(cs tls.ConnectionState, hostname string) ([]byte, error) {
	buf := bytes.NewBuffer([]byte(""))

	verifyEncode := func(chain []*x509.Certificate) error {
		for _, cert := range chain {
			if len(chain) > 1 {
				// drop the leaf certificate from the chain. osqueryd does not
				// need it to establish a secure connection
				if err := cert.VerifyHostname(hostname); err == nil {
					continue
				}
			}
			if err := encodePEMCertificate(buf, cert); err != nil {
				return err
			}
		}
		return nil
	}

	// use verified chains if available(which adds the root CA), otherwise
	// use the certificate chain offered by the server (if terminated with
	// self-signed certs)
	if len(cs.VerifiedChains) != 0 {
		for _, chain := range cs.VerifiedChains {
			if err := verifyEncode(chain); err != nil {
				return nil, errors.Wrap(err, "encode verified chains pem")
			}
		}
	} else {
		if err := verifyEncode(cs.PeerCertificates); err != nil {
			return nil, errors.Wrap(err, "encode peer certificates pem")
		}
	}
	return buf.Bytes(), nil
}

func encodePEMCertificate(buf io.Writer, cert *x509.Certificate) error {
	block := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	}
	return pem.Encode(buf, block)
}
