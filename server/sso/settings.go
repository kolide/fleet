package sso

import (
	"encoding/xml"
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
