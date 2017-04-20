package sso

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestCompression(t *testing.T) {
	input := "<samlp:AuthnRequest AssertionConsumerServiceURL='https://sp.example.com/acs' Destination='https://idp.example.com/sso' ID='_18185425-fd62-477c-b9d4-4b5d53a89845' IssueInstant='2017-04-16T15:32:42Z' ProtocolBinding='urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST' Version='2.0' xmlns:saml='urn:oasis:names:tc:SAML:2.0:assertion' xmlns:samlp='urn:oasis:names:tc:SAML:2.0:protocol'><saml:Issuer>https://sp.example.com/saml2</saml:Issuer><samlp:NameIDPolicy AllowCreate='true' Format='urn:oasis:names:tc:SAML:2.0:nameid-format:transient'/></samlp:AuthnRequest>"
	expected := "fJJf79IwFIa%2FSu961f0pG4yGLZkQ4xLUBaYX3piyHaTJ2s6eTvHbmw2McPHjtnne9u1zzgal7gdRjv5iDvBzBPSkRATnlTVba3DU4I7gfqkWvhz2Ob14P6AIQxwCuEo99BC0VoeyRUp2gF4ZOUX%2Fg6p7JhEtJdUup9%2FjLM7ShKfs3C05S1arlp3WXcKSU9qlC5mtsySlpEIcoTLopfE55VG8YlHC4mUTp2LBRcK%2FUVI7621r%2B3fKdMr8yOnojLASFQojNaDwrTiWH%2FeCB5E43SAUH5qmZvXnY0PJV3A4t%2BZBRMlV9wbFZOb1TfKfqMfI8Doz3KvSYlYv5u%2B54g2tE8I34SN5n9gnqaHa1bZX7R9S9r39vXUgPeTUuxEoeW%2Bdlv51l%2BlEdew8o8I7aVCB8TQsbk8%2B70XxFw%3D%3D"
	buff := bytes.NewBufferString(input)
	compressed, err := deflate(buff)
	require.Nil(t, err)
	assert.Equal(t, expected, compressed)

	// reader, err := inflate([]byte(compressed))
	// require.Nil(t, err)
	// var authReq AuthnRequest
	// err = xml.NewDecoder(reader).Decode(&authReq)
	// require.Nil(t, err)
	// assert.Equal(t, "https://sp.example.com/acs", authReq.AssertionConsumerServiceURL)
	// assert.Equal(t, "https://idp.example.com/sso", authReq.Destination)
	// assert.Equal(t, "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST", authReq.ProtocolBinding)

}
