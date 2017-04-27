package sso

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"net/url"

	"github.com/pkg/errors"
	"github.com/y0ssar1an/q"
)

const (
	// These are response status codes described in the core SAML spec section
	// 3.2.2.1 See http://docs.oasis-open.org/security/saml/v2.0/saml-core-2.0-os.pdf
	Success int = iota
	Requestor
	Responder
	VersionMismatch
	AuthnFailed
	InvalidAttrNameOrValue
	InvalidNameIDPolicy
	NoAuthnContext
	NoAvailableIDP
	NoPassive
	NoSupportedIDP
	PartialLogout
	ProxyCountExceeded
	RequestDenied
	RequestUnsupported
	RequestVersionDeprecated
	RequestVersionTooHigh
	RequestVersionTooLow
	ResourceNotRecognized
	TooManyResponses
	UnknownAttrProfile
	UnknownPrincipal
	UnsupportedBinding
)

var statusMap = map[string]int{
	"urn:oasis:names:tc:SAML:2.0:status:Success":                  Success,
	"urn:oasis:names:tc:SAML:2.0:status:Requester":                Requestor,
	"urn:oasis:names:tc:SAML:2.0:status:Responder":                Responder,
	"urn:oasis:names:tc:SAML:2.0:status:VersionMismatch":          VersionMismatch,
	"urn:oasis:names:tc:SAML:2.0:status:AuthnFailed":              AuthnFailed,
	"urn:oasis:names:tc:SAML:2.0:status:InvalidAttrNameOrValue":   InvalidAttrNameOrValue,
	"urn:oasis:names:tc:SAML:2.0:status:InvalidNameIDPolicy":      InvalidNameIDPolicy,
	"urn:oasis:names:tc:SAML:2.0:status:NoAuthnContext":           NoAuthnContext,
	"urn:oasis:names:tc:SAML:2.0:status:NoAvailableIDP":           NoAvailableIDP,
	"urn:oasis:names:tc:SAML:2.0:status:NoPassive":                NoPassive,
	"urn:oasis:names:tc:SAML:2.0:status:NoSupportedIDP":           NoSupportedIDP,
	"urn:oasis:names:tc:SAML:2.0:status:PartialLogout":            PartialLogout,
	"urn:oasis:names:tc:SAML:2.0:status:ProxyCountExceeded":       ProxyCountExceeded,
	"urn:oasis:names:tc:SAML:2.0:status:RequestDenied":            RequestDenied,
	"urn:oasis:names:tc:SAML:2.0:status:RequestUnsupported":       RequestUnsupported,
	"urn:oasis:names:tc:SAML:2.0:status:RequestVersionDeprecated": RequestVersionDeprecated,
	"urn:oasis:names:tc:SAML:2.0:status:RequestVersionTooLow":     RequestVersionTooLow,
	"urn:oasis:names:tc:SAML:2.0:status:ResourceNotRecognized":    ResourceNotRecognized,
	"urn:oasis:names:tc:SAML:2.0:status:TooManyResponses":         TooManyResponses,
	"urn:oasis:names:tc:SAML:2.0:status:UnknownAttrProfile":       UnknownAttrProfile,
	"urn:oasis:names:tc:SAML:2.0:status:UnknownPrincipal":         UnknownPrincipal,
	"urn:oasis:names:tc:SAML:2.0:status:UnsupportedBinding":       UnsupportedBinding,
}

type AuthInfo interface {
	RelayState() string
	UserID() string
	Status() (int, error)
	StatusDescription() string
}

type resp struct {
	relayState string
	userID     string
	status     string
}

func (r resp) StatusDescription() string {
	return r.status
}

func (r resp) RelayState() string {
	return r.relayState
}

func (r resp) UserID() string {
	return r.userID
}

func (r resp) Status() (int, error) {
	if r.status == "" {
		return AuthnFailed, errors.New("no status present")
	}
	if s, ok := statusMap[r.status]; ok {
		return s, nil
	}
	return AuthnFailed, errors.Errorf("unhandled status %s", r.status)
}

// DecodeAuthResponse extracts SAML assertions from IDP response
func DecodeAuthResponse(samlResponse, relayState string) (AuthInfo, error) {
	q.Q(samlResponse)
	q.Q(relayState)
	var authInfo resp
	decoded, err := base64.StdEncoding.DecodeString(samlResponse)
	if err != nil {
		return nil, errors.Wrap(err, "decoding saml response")
	}
	var saml Response
	err = xml.NewDecoder(bytes.NewBuffer(decoded)).Decode(&saml)
	if err != nil {
		return nil, errors.Wrap(err, "decoding response xml")
	}
	authInfo.status = saml.Status.StatusCode.Value
	status, err := authInfo.Status()
	if err != nil {
		return nil, errors.Wrap(err, "decoding auth response")
	}
	if status == Success {
		authInfo.userID = saml.Assertion.Subject.NameID.Value
	}
	authInfo.relayState = relayState
	return &authInfo, nil
}

// See http://docs.oasis-open.org/security/saml/v2.0/saml-bindings-2.0-os.pdf
// Section 3.5.4 for details on decoding SAML Response
// Also we expect form to be application/x-www-form-urlencoded
// See https://www.w3.org/TR/html401/interact/forms.html section 17.13.4
func decodeAuthSAMLResponse(encodedResp []byte) (*Response, error) {
	unescaped, err := url.PathUnescape(string(encodedResp))
	if err != nil {
		return nil, errors.Wrap(err, "processing SAMLResponse")
	}
	decoded, err := base64.StdEncoding.DecodeString(unescaped)
	if err != nil {
		return nil, errors.Wrap(err, "base46 decoding auth response")
	}
	var resp Response
	reader := bytes.NewBuffer(decoded)
	err = xml.NewDecoder(reader).Decode(&resp)
	if err != nil {
		return nil, errors.Wrap(err, "decoding authorization response xml")
	}
	return &resp, nil
}
