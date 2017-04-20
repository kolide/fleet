package sso

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"io"
	"net/url"

	"github.com/pkg/errors"
)

type AuthInfo interface {
	RelayState() (string, bool)
	UserID() (string, bool)
}

type resp struct {
	relayState string
	userID     string
}

func (r *resp) RelayState() (string, bool) {
	if r.relayState == "" {
		return "", false
	}
	return r.relayState, true
}

func (r *resp) UserID() (string, bool) {
	if r.userID == "" {
		return "", false
	}
	return r.userID, true
}

func DecodeAuthResponse(body io.Reader) (AuthInfo, error) {
	var dest bytes.Buffer
	_, err := io.Copy(&dest, body)
	if err != nil {
		return nil, errors.Wrap(err, "malformed auth response")
	}
	// parse form name/value pairs
	args := bytes.Split(dest.Bytes(), []byte("&"))
	params := make(map[string][]byte)
	for _, arg := range args {
		// seperate name and values and assign to map to contextualize the
		// form data
		vals := bytes.Split(arg, []byte("="))
		if len(vals) != 2 {
			return nil, errors.New("auth response form argument malformed")
		}
		params[string(vals[0])] = vals[1]
	}
	// We MUST have SAMLResponse, RelayState is also required per the spec as we supply
	// it in the auth request
	var response resp
	if _, ok := params["RelayState"]; !ok {
		return nil, errors.New("missing required RelayState")
	}
	response.relayState, err = url.PathUnescape(string(params["RelayState"]))
	// SAMLResponse is required
	if _, ok := params["SAMLResponse"]; !ok {
		return nil, errors.New("missing required SAMLResponse parameter")
	}
	authResp, err := decodeAuthSAMLResponse(params["SAMLResponse"])
	if err != nil {
		return nil, errors.Wrap(err, "decoding authorization response")
	}
	response.userID = authResp.Assertion.Subject.NameID.Value
	return &response, nil
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
