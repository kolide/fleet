package service

import (
	"testing"
	"time"

	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/kolide/kolide-ose/server/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func TestLicenseService(t *testing.T) {
	tokenString := "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJsaWNlbnNlX3V1aWQiOiI3M" +
		"DA0ODhiNy03OTliLTRlNDYtOTM0Yy1lOWYxNWMwNjFlMGQiLCJvcmdhbml6YXRpb25fbmFtZSI6Il" +
		"BoYW50YXNtLCBJbmMuIiwib3JnYW5pemF0aW9uX3V1aWQiOiJhZTVhNDU0OC02NTJhLTQ4YjktYTE1Y" +
		"y01NDRlNWYyZmFlY2IiLCJob3N0X2xpbWl0IjowLCJldmFsdWF0aW9uIjp0cnVlLCJleHBpcmVzX2F0" +
		"IjoiMjAxNy0wMi0yNiAxNTo0ODowNyBVVEMifQ.wveyNQNQ6YXA8eQuznWifiZgzYBR9hofqe6lA" +
		"Kh-5sBVSo-4RgliEwQaBBc5DfDcQkkim9TtfYHouWUQ8AG3wN2fTwkLK2thF41cCDivN5nhC93KEb" +
		"VkE2M2yvm2zPzlCV79KX9y7d9bPpG6hZnfJWHK_wpJv-iu6BhlVJhR8-1-5jQOLmsLWSNsZju9RrlL" +
		"1njDw6ktPqz-kfL7jELokj_6rWXZz4q7rM0gKQnGSYZUJqMnfP2F807DbdCxHV-c6bEFw7cTYqGl" +
		"sAOUH2JLC-3OMObEvSK8Mi2WrqqTWnTnGHp1WNyrS6Wl0-kcp_tc31ijOfcmyFzhBea5WCYf8A"
	ds := new(mock.Store)
	ds.LicenseFunc = func() (*kolide.License, error) {
		result := &kolide.License{
			UpdateTimestamp: kolide.UpdateTimestamp{
				UpdatedAt: time.Now().Add(-5 * time.Minute),
			},
			Token: &tokenString,
			PublicKey: "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AM" +
				"IIBCgKCAQEA0ZhY7r6HmifXPtServt4\nD3MSi8Awe9u132vLf8yzlknvnq+8CSnOPSSbC" +
				"D+HajvZ6dnNJXjdcAhuZ32ShrH8\nrEQACEUS8Mh4z8Mo5Nlq1ou0s2JzWCx049kA34jP" +
				"3u6AiPgpWUf8JRGstTlisxMn\nH6B7miDs1038gVbN5rk+j+3ALYzllaTnCX3Y0C7f6IW7B" +
				"jNO/tvFB84/95xfOLEz\no2MeFMqkD29hvcrUW+8+fQGJaVLvcEqBDnIEVbCCk8Wnoi48d" +
				"UE06WHUl6voJecD\ndW1E6jHcq8PQFK+4bI1gKZVbV4dFGSSMUyD7ov77aWHjxdQe6YEGc" +
				"SXKzfyMaUtQ\nvQIDAQAB\n-----END PUBLIC KEY-----\n",
			Revoked: false,
			ID:      1,
		}
		return result, nil
	}

	svc, err := newTestService(ds, nil)
	require.Nil(t, err)
	ctx := context.Background()
	lic, err := svc.License(ctx)
	require.Nil(t, err)
	claims, err := lic.Claims()
	require.Nil(t, err)
	require.NotNil(t, claims)

	assert.False(t, claims.Revoked)
	assert.Equal(t, "700488b7-799b-4e46-934c-e9f15c061e0d", claims.LicenseUUID)
	assert.Equal(t, "Phantasm, Inc.", claims.OrganizationName)
	assert.Equal(t, "ae5a4548-652a-48b9-a15c-544e5f2faecb", claims.OrganizationUUID)
	assert.Equal(t, 0, claims.HostLimit)
	assert.Equal(t, "2017-02-26 15:48:07 UTC", claims.ExpiresAt.Format(kolide.LicenseTimeLayout))

}
