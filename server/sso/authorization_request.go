package sso

import (
	"bytes"
	"compress/zlib"
	"crypto/rand"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	samlVersion = "2.0"
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
func CreateAuthorizationRequest(settings *Settings, options ...func(o *opts)) (string, error) {
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
		Issuer: Issuer{
			XMLName: xml.Name{
				Local: "saml:Issuer",
			},
			Url: settings.Metadata.EntityID,
		},
	}

	queryVals := make(map[string]string)

	var writer bytes.Buffer
	err = xml.NewEncoder(&writer).Encode(request)
	if err != nil {
		return "", errors.Wrap(err, "encoding auth request xml")
	}
	authQueryVal, err := deflate(&writer)
	if err != nil {
		return "", errors.Wrap(err, "unable to compress auth info")
	}
	queryVals["SAMLRequest"] = authQueryVal
	if optionalParams.relayState != "" {
		queryVals["RelayState"] = optionalParams.relayState
	}
	if settings.Metadata.IDPSSODescriptor.WantAuthnRequestsSigned {
		signature, err := sign(settings, authQueryVal)
		if err != nil {
			return "", errors.Wrap(err, "signing auth request")
		}
		queryVals["Signature"] = signature
	}
	return buildRedirectURL(destinationURL, queryVals), nil
}

func buildRedirectURL(baseURL string, params map[string]string) string {
	queryString := ""
	for key, val := range params {
		if queryString == "" {
			queryString = "?"
		} else {
			queryString += "&"
		}
		queryString += key + "=" + val
	}
	return baseURL + queryString
}

func sign(settings *Settings, authString string) (string, error) {
	return "", nil
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
	writer := zlib.NewWriter(&deflated)
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
	// We have to remove the compression method, and flag bytes from the front
	// of the byte stream and the 32 bit checksum at the end. This is to
	// retain compatibility with PKZIP and GZIP
	// See https://tools.ietf.org/html/rfc1950
	// If you change this, the Auth won't work and you'll be sad.
	encbuff = encbuff[2 : len(encbuff)-4]
	encoded := base64.StdEncoding.EncodeToString(encbuff)
	// replace any whitespace, and URL encode base 64 output
	encoded = urlEncode(encoded)
	return encoded, nil
}

func urlEncode(val string) string {
	// replace any whitespace, and URL encode base 64 output
	return strings.NewReplacer("\n", "", "+", "%2B", "/", "%2F", "=", "%3D").Replace(val)
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
