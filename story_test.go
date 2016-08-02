package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type integrationRequests struct {
	r  *gin.Engine
	db *gorm.DB
}

func (req *integrationRequests) New() {
	req.r = createTestServer()
	req.r.Use(testSessionMiddleware)
	req.r.Use(JWTRenewalMiddleware)

	req.db, _ = openTestDB()
	injectedTestDB = req.db

	// Until we have a better solution for first-user onboarding, manually
	// create an admin
	_, err := NewUser(req.db, "admin", "foobar", "admin@kolide.co", true, false)
	if err != nil {
		panic(err.Error())
	}

	req.r.POST("/login", Login)
	req.r.GET("/logout", Logout)

	req.r.POST("/user", GetUser)
	req.r.PUT("/user", CreateUser)
	req.r.PATCH("/user", ModifyUser)
	req.r.DELETE("/user", DeleteUser)

	req.r.PATCH("/user/password", ChangeUserPassword)
	req.r.PATCH("/user/admin", SetUserAdminState)
	req.r.PATCH("/user/enabled", SetUserEnabledState)
}

func (req *integrationRequests) Login(username, password string, sessionOut *string) error {
	response := httptest.NewRecorder()
	body, err := json.Marshal(LoginRequestBody{
		Username: username,
		Password: password,
	})
	if err != nil {
		return err
	}

	buff := new(bytes.Buffer)
	buff.Write(body)
	request, _ := http.NewRequest("POST", "/login", buff)
	request.Header.Set("Content-Type", "application/json")
	req.r.ServeHTTP(response, request)

	if response.Code != 200 {
		return errors.New(fmt.Sprintf("Response code: %d", response.Code))
	}
	*sessionOut = response.Header().Get("Set-Cookie")

	return nil
}

func (req *integrationRequests) CreateUser(username, password, email string, admin, reset bool, session *string) (*GetUserResponseBody, error) {
	response := httptest.NewRecorder()
	body, err := json.Marshal(CreateUserRequestBody{
		Username:           username,
		Password:           password,
		Email:              email,
		Admin:              admin,
		NeedsPasswordReset: reset,
	})
	if err != nil {
		return nil, err
	}

	buff := new(bytes.Buffer)
	buff.Write(body)
	request, _ := http.NewRequest("PUT", "/user", buff)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cookie", *session)
	req.r.ServeHTTP(response, request)

	if response.Code != 200 {
		return nil, errors.New(fmt.Sprintf("Response code: %d", response.Code))
	}
	*session = response.Header().Get("Set-Cookie")

	var responseBody GetUserResponseBody
	err = json.Unmarshal(response.Body.Bytes(), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

func (req *integrationRequests) GetUser(username string, session *string) (*GetUserResponseBody, error) {
	response := httptest.NewRecorder()
	body, err := json.Marshal(GetUserRequestBody{
		Username: username,
	})
	if err != nil {
		return nil, err
	}

	buff := new(bytes.Buffer)
	buff.Write(body)
	request, _ := http.NewRequest("POST", "/user", buff)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cookie", *session)
	req.r.ServeHTTP(response, request)

	if response.Code != 200 {
		return nil, errors.New(fmt.Sprintf("Response code: %d", response.Code))
	}
	*session = response.Header().Get("Set-Cookie")

	var responseBody GetUserResponseBody
	err = json.Unmarshal(response.Body.Bytes(), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

func (req *integrationRequests) ModifyUser(username, name, email string, session *string) (*GetUserResponseBody, error) {
	response := httptest.NewRecorder()
	body, err := json.Marshal(ModifyUserRequestBody{
		Username: username,
		Name:     name,
		Email:    email,
	})
	if err != nil {
		return nil, err
	}

	buff := new(bytes.Buffer)
	buff.Write(body)
	request, _ := http.NewRequest("PATCH", "/user", buff)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cookie", *session)
	req.r.ServeHTTP(response, request)

	if response.Code != 200 {
		return nil, errors.New(fmt.Sprintf("Response code: %d", response.Code))
	}
	*session = response.Header().Get("Set-Cookie")

	var responseBody GetUserResponseBody
	err = json.Unmarshal(response.Body.Bytes(), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

func (req *integrationRequests) DeleteUser(username string, session *string) error {
	response := httptest.NewRecorder()
	body, err := json.Marshal(DeleteUserRequestBody{
		Username: username,
	})
	if err != nil {
		return err
	}

	buff := new(bytes.Buffer)
	buff.Write(body)
	request, _ := http.NewRequest("DELETE", "/user", buff)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cookie", *session)
	req.r.ServeHTTP(response, request)

	if response.Code != 200 {
		return errors.New(fmt.Sprintf("Response code: %d", response.Code))
	}
	*session = response.Header().Get("Set-Cookie")

	return nil
}

func (req *integrationRequests) ChangePassword(username, currentPassword, newPassword string, session *string) (*GetUserResponseBody, error) {
	response := httptest.NewRecorder()
	body, err := json.Marshal(ChangePasswordRequestBody{
		Username:          username,
		CurrentPassword:   currentPassword,
		NewPassword:       newPassword,
		NewPasswordConfim: newPassword,
	})
	if err != nil {
		return nil, err
	}

	buff := new(bytes.Buffer)
	buff.Write(body)
	request, _ := http.NewRequest("PATCH", "/user/password", buff)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cookie", *session)
	req.r.ServeHTTP(response, request)

	if response.Code != 200 {
		return nil, errors.New(fmt.Sprintf("Response code: %d", response.Code))
	}
	*session = response.Header().Get("Set-Cookie")

	var responseBody GetUserResponseBody
	err = json.Unmarshal(response.Body.Bytes(), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

func (req *integrationRequests) SetAdminState(username string, admin bool, session *string) (*GetUserResponseBody, error) {
	response := httptest.NewRecorder()
	body, err := json.Marshal(SetUserAdminStateRequestBody{
		Username: username,
		Admin:    admin,
	})
	if err != nil {
		return nil, err
	}

	buff := new(bytes.Buffer)
	buff.Write(body)
	request, _ := http.NewRequest("PATCH", "/user/admin", buff)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cookie", *session)
	req.r.ServeHTTP(response, request)

	if response.Code != 200 {
		return nil, errors.New(fmt.Sprintf("Response code: %d", response.Code))
	}
	*session = response.Header().Get("Set-Cookie")

	var responseBody GetUserResponseBody
	err = json.Unmarshal(response.Body.Bytes(), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

func (req *integrationRequests) SetEnabledState(username string, enabled bool, session *string) (*GetUserResponseBody, error) {
	response := httptest.NewRecorder()
	body, err := json.Marshal(SetUserEnabledStateRequestBody{
		Username: username,
		Enabled:  enabled,
	})
	if err != nil {
		return nil, err
	}

	buff := new(bytes.Buffer)
	buff.Write(body)
	request, _ := http.NewRequest("PATCH", "/user/enabled", buff)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cookie", *session)
	req.r.ServeHTTP(response, request)

	if response.Code != 200 {
		return nil, errors.New(fmt.Sprintf("Response code: %d", response.Code))
	}
	*session = response.Header().Get("Set-Cookie")

	var responseBody GetUserResponseBody
	err = json.Unmarshal(response.Body.Bytes(), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

func (req *integrationRequests) CheckUser(username, email, name string, admin, reset, enabled bool) error {
	var user User
	err := req.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return err
	}
	if user.Email != email {
		return errors.New(fmt.Sprintf("user's email was not set in the DB: %s", user.Email))
	}
	if (user.Admin && !admin) || (!user.Admin && admin) {
		return errors.New("user shouldn't be admin")
	}
	if (user.NeedsPasswordReset && !reset) || (!user.NeedsPasswordReset && reset) {
		return errors.New("user reset settings don't match")
	}
	if (user.Enabled && !enabled) || (!user.Enabled && enabled) {
		return errors.New("user enabled settings don't match")
	}
	if user.Name != name {
		return errors.New(fmt.Sprintf("user names don't match: %s and %s", user.Name, name))
	}
	return nil
}

func (req *integrationRequests) GetAndCheckUser(username string, session *string) error {
	resp, err := req.GetUser(username, session)
	if err != nil {
		return err
	}

	err = req.CheckUser(username, resp.Email, resp.Name, resp.Admin, resp.NeedsPasswordReset, resp.Enabled)
	if err != nil {
		return err
	}

	return nil
}

func (req *integrationRequests) CreateAndCheckUser(username, password, email, name string, admin, reset bool, session *string) error {
	resp, err := req.CreateUser(username, password, email, admin, reset, session)
	if err != nil {
		return err
	}

	err = req.CheckUser(username, email, name, admin, reset, resp.Enabled)
	if err != nil {
		return err
	}

	return nil
}

func (req *integrationRequests) ModifyAndCheckUser(username, email, name string, admin, reset bool, session *string) error {
	resp, err := req.ModifyUser(username, name, email, session)
	if err != nil {
		return err
	}

	err = req.CheckUser(username, email, name, admin, reset, resp.Enabled)
	if err != nil {
		return err
	}

	return nil
}

func (req *integrationRequests) DeleteAndCheckUser(username string, session *string) error {
	err := req.DeleteUser(username, session)
	if err != nil {
		return err
	}

	var user User
	err = req.db.Where("username = ?", username).First(&user).Error
	if err == nil {
		return errors.New("User should have been deleted.")
	}

	return nil
}

func (req *integrationRequests) SetEnabledStateAndCheckUser(username string, enabled bool, session *string) error {
	resp, err := req.SetEnabledState(username, enabled, session)
	if err != nil {
		return err
	}

	err = req.CheckUser(username, resp.Email, resp.Name, resp.Admin, resp.NeedsPasswordReset, enabled)
	if err != nil {
		return err
	}

	return nil
}

func (req *integrationRequests) SetAdminStateAndCheckUser(username string, admin bool, session *string) error {
	resp, err := req.SetAdminState(username, admin, session)
	if err != nil {
		return err
	}

	err = req.CheckUser(username, resp.Email, resp.Name, admin, resp.NeedsPasswordReset, resp.Enabled)
	if err != nil {
		return err
	}

	return nil
}

func TestUserAndAccountManagement(t *testing.T) {

	// Create and configure the webserver which will be used to handle the tests
	var req integrationRequests
	req.New()

	// Instantiate the variables that will store the most recent session cookie
	// for each user context that will be created
	var adminSession string
	var admin2Session string
	var user1Session string
	var user2Session string

	var err error

	// Test logging in with the first admin
	err = req.Login("admin", "foobar", &adminSession)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Once admin is logged in, create a user using a valid admin session
	err = req.CreateAndCheckUser("user1", "foobar", "user1@kolide.co", "", false, false, &adminSession)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Once admin is logged in, create another admin account using a valid
	// admin session
	err = req.CreateAndCheckUser("admin2", "foobar", "admin2@kolide.co", "", true, false, &adminSession)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Once admin has created admin2, log in with admin2 to get a session
	// context for admin2
	err = req.Login("admin2", "foobar", &admin2Session)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Use an admin created via the API to create a user via the API
	err = req.CreateAndCheckUser("user2", "foobar", "user2@kolide.co", "", false, false, &admin2Session)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Once admin has created user1, log in with user1 to get a session context
	// for user1
	err = req.Login("user1", "foobar", &user1Session)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Once admin2 has created user2, log in with user1 to get a session context
	// for user2
	err = req.Login("user2", "foobar", &user2Session)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Get info on user2 as admin2
	err = req.GetAndCheckUser("user2", &admin2Session)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Get info on admin2 as user2
	err = req.GetAndCheckUser("admin2", &user2Session)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Modify user1 as admin
	err = req.ModifyAndCheckUser("user1", "user1@kolide.co", "User One", false, false, &adminSession)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Modify user2 as user2
	err = req.ModifyAndCheckUser("user2", "user2@kolide.co", "User Two", false, false, &user2Session)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Modify user2 as user1
	err = req.ModifyAndCheckUser("user2", "user1@kolide.co", "Less Cool User", false, false, &user1Session)
	if err == nil {
		t.Fatal("Action should not be authorized")
	}

	// admin resets user1 password
	_, err = req.ChangePassword("user1", "", "bazz1", &adminSession)
	if err != nil {
		t.Fatal(err.Error())
	}

	// user1 logs in with new password
	err = req.Login("user1", "bazz1", &user1Session)
	if err != nil {
		t.Fatal(err.Error())
	}

	// user2 resets user2 password
	_, err = req.ChangePassword("user2", "foobar", "bazz2", &user2Session)
	if err != nil {
		t.Fatal(err.Error())
	}

	// user2 logs in with new password
	err = req.Login("user2", "bazz2", &user2Session)
	if err != nil {
		t.Fatal(err.Error())
	}

	// user2 tries to change user1 password
	_, err = req.ChangePassword("user1", "", "fake", &user2Session)
	if err == nil {
		t.Fatal("Action should not be authorized")
	}

	// admin2 promotes user2 to admin
	err = req.SetAdminStateAndCheckUser("user2", true, &admin2Session)
	if err != nil {
		t.Fatal(err.Error())
	}

	// user2 is admin
	resp, err := req.GetUser("user2", &user2Session)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !resp.Admin {
		t.Fatal("user2 should be an admin")
	}

	// admin demotes user2 from admin
	err = req.SetAdminStateAndCheckUser("user2", false, &adminSession)
	if err != nil {
		t.Fatal(err.Error())
	}

	// user2 is no longer an admin
	resp, err = req.GetUser("user2", &user2Session)
	if err != nil {
		t.Fatal(err.Error())
	}
	if resp.Admin {
		t.Fatal("user2 shouldn't be an admin")
	}

	// admin sets user1 as no longer enabled
	err = req.SetEnabledStateAndCheckUser("user1", false, &adminSession)
	if err != nil {
		t.Fatal(err.Error())
	}

	// user1 is no longer enabled
	resp, err = req.GetUser("user1", &user2Session)
	if err != nil {
		t.Fatal(err.Error())
	}
	if resp.Enabled {
		t.Fatal("user1 shouldn't be enabled")
	}

	// user1 can't view user2
	_, err = req.GetUser("user2", &user1Session)
	if err == nil {
		t.Fatal("Action shouldn't be authorized")
	}

	// admin2 re-enables user1
	err = req.SetEnabledStateAndCheckUser("user1", true, &admin2Session)
	if err != nil {
		t.Fatal(err.Error())
	}

	// user1 can view user2
	_, err = req.GetUser("user2", &user2Session)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Delete admin2 as admin1
	err = req.DeleteAndCheckUser("admin2", &adminSession)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Delete admin as user1
	err = req.DeleteAndCheckUser("admin", &user1Session)
	if err == nil {
		t.Fatal("Action should not be authorized")
	}

	// Delete user2 as user1
	err = req.DeleteAndCheckUser("user2", &user1Session)
	if err == nil {
		t.Fatal("Action should not be authorized")
	}

	// Delete user2 as admin
	err = req.DeleteAndCheckUser("user2", &adminSession)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Test get info from deleted account
	_, err = req.GetUser("admin", &admin2Session)
	if err == nil {
		t.Fatal("User session still works after account delete")
	}
}
