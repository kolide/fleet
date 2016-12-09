package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/y0ssar1an/q"
)

func (s *EndpointTestSuite) TestGetAppConfig() {
	q.Q(s.server.URL + "/api/v1/kolide/config")
	req, err := http.NewRequest("GET", s.server.URL+"/api/v1/kolide/config", nil)
	s.Require().Nil(err)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.userToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	s.Require().Nil(err)

	s.Require().Equal(http.StatusOK, resp.StatusCode)
	var configInfo getAppConfigResponse
	err = json.NewDecoder(resp.Body).Decode(&configInfo)
	s.Require().Nil(err)
	s.Require().NotNil(configInfo.AppConfig)
	config := configInfo.AppConfig
	s.Equal(uint(465), config.Port)
	s.Equal("Kolide", config.OrgName)
	s.Equal("http://foo.bar/image.png", config.OrgLogoURL)

}
