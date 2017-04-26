package sso

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type IDPKeyDescriptor struct {
	Use         string `xml:"use,attr"`
	Certificate string `xml:"KeyInfo>X509Data>X509Certificate"`
}

type IDPSSODescriptor struct {
	WantAuthnRequestsSigned    bool     `xml:"WantAuthnRequestsSigned,attr"`
	ProtocolSupportEnumeration string   `xml:"protocolSupportEnumeration,attr"`
	NameIDFormat               []string `xml:"NameIDFormat"`
	SingleSignOnService        []SingleSignOnService
	KeyDescriptor              IDPKeyDescriptor
}

type SingleSignOnService struct {
	Binding  string `xml:"Binding,attr"`
	Location string `xml:"Location,attr"`
}

type IDPMetadata struct {
	XMLName          xml.Name `xml:"EntityDescriptor"`
	EntityID         string   `xml:"entityID,attr"`
	IDPSSODescriptor IDPSSODescriptor
}

const (
	PasswordProtectedTransport = "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"
	PostBinding                = "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST"
	RedirectBinding            = "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect"
)

type Settings struct {
	Metadata *IDPMetadata
	// AssertionConsumerServiceURL is the call back on the service provider which responds
	// to the IDP
	AssertionConsumerServiceURL string
}

// ParseMetadata writes metadata xml to a struct
func ParseMetadata(metadata string) (*IDPMetadata, error) {
	var md IDPMetadata
	err := xml.Unmarshal([]byte(metadata), &md)
	if err != nil {
		return nil, err
	}
	return &md, nil
}

// GetMetadata retrieves information describing how to interact with a particular
// IDP via a remote URL. metadataURL is the location where the metadata is located
// and timeout defines how long to wait to get a response form the metadata
// server.
func GetMetadata(metadataURL string, timeout time.Duration) (*IDPMetadata, error) {
	request, err := http.NewRequest(http.MethodGet, metadataURL, nil)
	if err != nil {
		return nil, err
	}
	client := http.Client{Timeout: timeout}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("SAML metadata server at %s returned %s", metadataURL, resp.Status)
	}
	xmlData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var md IDPMetadata
	err = xml.Unmarshal(xmlData, &md)
	if err != nil {
		return nil, err
	}
	return &md, nil
}
