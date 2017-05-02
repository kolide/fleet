package sso

import (
	"bytes"
	"compress/flate"
	"crypto/rand"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

const (
	samlVersion   = "2.0"
	cacheLifetime = 300 // five minutes
)

// RelayState sets optional relay state
func RelayState(v string) func(*opts) {
	return func(o *opts) {
		o.relayState = v
	}
}

type opts struct {
	relayState string
}

// CreateAuthorizationRequest creates a url suitable for use to satisfy the SAML
// redirect binding.
// See http://docs.oasis-open.org/security/saml/v2.0/saml-bindings-2.0-os.pdf Section 3.4
func CreateAuthorizationRequest(settings *Settings, issuer string, options ...func(o *opts)) (string, error) {
	var optionalParams opts
	for _, opt := range options {
		opt(&optionalParams)
	}
	if settings.Metadata == nil {
		return "", errors.New("missing settings metadata")
	}
	requestID, err := getAuthnRequestID()
	if err != nil {
		return "", errors.Wrap(err, "creating auth request id")
	}
	destinationURL, err := getDestinationURL(settings)
	if err != nil {
		return "", errors.Wrap(err, "creating auth request")
	}
	request := AuthnRequest{
		XMLName: xml.Name{
			Local: "samlp:AuthnRequest",
		},
		ID:    requestID,
		SAMLP: "urn:oasis:names:tc:SAML:2.0:protocol",
		SAML:  "urn:oasis:names:tc:SAML:2.0:assertion",
		AssertionConsumerServiceURL: settings.AssertionConsumerServiceURL,
		Destination:                 destinationURL,
		IssueInstant:                time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		ProtocolBinding:             RedirectBinding,
		Version:                     samlVersion,
		ProviderName:                "Kolide",
		Issuer: Issuer{
			XMLName: xml.Name{
				Local: "saml:Issuer",
			},
			Url: issuer,
		},
	}
	var reader bytes.Buffer
	err = xml.NewEncoder(&reader).Encode(settings.Metadata)
	if err != nil {
		return "", errors.Wrap(err, "encoding metadata creating auth request")
	}
	// cache metadata so we can check the signatures on the response we get from the IDP
	err = settings.SessionStore.create(requestID,
		settings.OriginalURL,
		reader.String(),
		cacheLifetime,
	)
	if err != nil {
		return "", errors.Wrap(err, "caching cert while creating auth request")
	}
	u, err := url.Parse(destinationURL)
	if err != nil {
		return "", errors.Wrap(err, "parsing destination url")
	}
	qry := u.Query()

	var writer bytes.Buffer
	err = xml.NewEncoder(&writer).Encode(request)
	if err != nil {
		return "", errors.Wrap(err, "encoding auth request xml")
	}
	authQueryVal, err := deflate(&writer)
	if err != nil {
		return "", errors.Wrap(err, "unable to compress auth info")
	}
	qry.Set("SAMLRequest", authQueryVal)
	if optionalParams.relayState != "" {
		qry.Set("RelayState", optionalParams.relayState)
	}
	u.RawQuery = qry.Encode()
	return u.String(), nil
}

func getDestinationURL(settings *Settings) (string, error) {
	for _, sso := range settings.Metadata.IDPSSODescriptor.SingleSignOnService {
		if sso.Binding == RedirectBinding {
			return sso.Location, nil
		}
	}
	return "", errors.New("IDP does not support redirect binding")
}

// See SAML Bindings http://docs.oasis-open.org/security/saml/v2.0/saml-bindings-2.0-os.pdf
// Section 3.4.4.1
func deflate(xmlBuffer *bytes.Buffer) (string, error) {
	var deflated bytes.Buffer
	writer, err := flate.NewWriter(&deflated, flate.DefaultCompression)
	if err != nil {
		return "", err
	}
	defer writer.Close()
	n, err := writer.Write(xmlBuffer.Bytes())
	if n != xmlBuffer.Len() {
		return "", errors.New("incomplete write during compression")
	}
	if err != nil {
		return "", errors.Wrap(err, "compressing auth request")
	}
	writer.Flush()
	encbuff := deflated.Bytes()
	encoded := base64.StdEncoding.EncodeToString(encbuff)
	return encoded, nil
}

// UUID per RFC 4122
func getAuthnRequestID() (string, error) {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		return "", err
	}
	// Variant field section 4.1.2
	// Set bit 7 (msb) to 1, bit 6 to 0, leave the rest of the bits alone
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// Version 4 psuedo random uuid msb octet = 0100
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
