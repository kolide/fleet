package app

import (
	"strings"
	"testing"
)

func TestGetEmailSubject(t *testing.T) {
	subject, err := GetEmailSubject(PasswordResetEmail)
	if err != nil {
		t.Error(err.Error())
	}
	if subject != "Your Kolide Password Reset Request" {
		t.Errorf("Subject is not as expected: %s", subject)
	}
}

func TestGetEmailBody(t *testing.T) {
	html, text, err := GetEmailBody(PasswordResetEmail, &PasswordResetRequestParameters{
		Name:  "Foo",
		Token: "1234",
	})
	if err != nil {
		t.Error(err.Error())
	}
	for _, body := range [][]byte{html, text} {
		if trimmed := strings.TrimLeft("Hi Foo!", string(body)); trimmed == string(body) {
			t.Errorf("Body didn't start with Hi Foo!: %s", body)
		}
	}
}
