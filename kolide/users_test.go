package kolide

import "testing"

func TestValidatePassword(t *testing.T) {

	var passwordTests = []struct {
		Username, Password, Email string
		Admin, PasswordReset      bool
	}{
		{"marpaia", "foobar", "mike@kolide.co", true, false},
		{"jason", "bar0baz!?", "jason@kolide.co", true, false},
	}

	for _, tt := range passwordTests {
		user := &User{
			Username:           tt.Username,
			Email:              tt.Email,
			Admin:              tt.Admin,
			NeedsPasswordReset: tt.PasswordReset,
		}
		user.setPassword(tt.Password, 60, 10)

		{
			err := user.ValidatePassword(tt.Password)
			if err != nil {
				t.Errorf("Password validation failed for user %s", user.Username)
			}
		}

		{
			err := user.ValidatePassword("different")
			if err == nil {
				t.Errorf("Incorrect password worked for user %s", user.Username)
			}
		}
	}
}
