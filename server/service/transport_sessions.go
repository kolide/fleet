package service

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

func decodeSubmitAuthnResponseRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	redirectURL := r.URL.Query().Get("redirect_url")
	token := r.URL.Query().Get("token")

	return submitAuthnResponseRequest{RedirectURL: redirectURL, Token: token}, nil
}

func encodeSubmitAuthnResponseResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	tmpl, err := template.New("foo").Parse(`
<html>
<script type='text/javascript'>
var redirectURL = '{{.RedirectURL}}';
if (!redirectURL.startsWith('/')) {
  redirectURL = '/';
}
window.localStorage.setItem('KOLIDE::auth_token', '{{.Token}}');
window.location = redirectURL;
</script>
<body>
Redirecting to Kolide...
</body>
</html>
`)
	if err != nil {
		return errors.Wrap(err, "parsing template")
	}

	if r, ok := response.(submitAuthnResponseResponse); ok {
		err = tmpl.Execute(w, r)
		if err != nil {
			return errors.Wrap(err, "executing template")
		}
	} else {
		return errors.Errorf("Unknown response type: %+v", response)
	}

	return nil
}

func decodeGetInfoAboutSessionRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	return getInfoAboutSessionRequest{ID: id}, nil
}

func decodeGetInfoAboutSessionsForUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	return getInfoAboutSessionsForUserRequest{ID: id}, nil
}

func decodeDeleteSessionRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	return deleteSessionRequest{ID: id}, nil
}

func decodeDeleteSessionsForUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	return deleteSessionsForUserRequest{ID: id}, nil
}

func decodeLoginRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	req.Username = strings.ToLower(req.Username)
	return req, nil
}
