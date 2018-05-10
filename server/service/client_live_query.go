package service

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

// LiveQuery creates a new live query and begins streaming results.
func (c *Client) LiveQuery(query string, labels []uint, hosts []uint) error {
	req := createDistributedQueryCampaignRequest{
		Query:    query,
		Selected: distributedQueryCampaignTargets{Labels: labels, Hosts: hosts},
	}
	response, err := c.AuthenticatedDo("POST", "/api/v1/kolide/queries/run", req)
	if err != nil {
		return errors.Wrap(err, "POST /api/v1/kolide/queries/run")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.Errorf(
			"create live query received status %d %s",
			response.StatusCode,
			extractServerErrorText(response.Body),
		)
	}

	var responseBody createDistributedQueryCampaignResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return errors.Wrap(err, "decode create live query response")
	}
	if responseBody.Err != nil {
		return errors.Errorf("create live query: %s", responseBody.Err)
	}

	// Copy default dialer but skip cert verification if set.
	dialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
		TLSClientConfig:  &tls.Config{InsecureSkipVerify: c.insecureSkipVerify},
	}

	wssURL := *c.baseURL
	wssURL.Scheme = "wss"
	wssURL.Path = "/api/v1/kolide/results/websocket"
	conn, _, err := dialer.Dial(wssURL.String(), nil)
	if err != nil {
		return errors.Wrap(err, "upgrade live query result websocket")
	}
	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"auth","data":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzZXNzaW9uX2tleSI6Ik5JMzFyZitjQVk0RUFPWTlvVFU5L1NSK2g1cGlWcFZ4bVpMbTNUeEFET2hoME00d0liaWR3OHRyM1JWbHovVU5SeWMveEZKZStqVHo3TzNQYUFxMWt3PT0ifQ.7tb_fIjq94EybEmtwvO_n_54ii_YLvZIRriYmVQGPc0"}}`))
	if err != nil {
		return errors.Wrap(err, "auth for results")
	}

	err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(`{"type":"select_campaign","data":{"campaign_id":%d}}`, responseBody.Campaign.ID)))
	if err != nil {
		return errors.Wrap(err, "auth for results")
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return errors.Wrap(err, "receive ws message")
		}
		fmt.Println(string(message))
	}

	return nil
}
