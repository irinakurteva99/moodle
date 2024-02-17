package moodle

import (
	"context"
	"path"

	"github.com/irinakurteva99/moodle/pkg/urlutil"
)

type AuthAPI interface {
	Login(ctx context.Context, username, password string) (*LoginResponse, error)
}

type authAPI struct {
	*apiClient
}

func newAuthAPI(apiClient *apiClient) *authAPI {
	return &authAPI{apiClient}
}

type LoginResponse struct {
	Token        string `json:"token"`
	PrivateToken string `json:"privatetoken"`
}

func (a *authAPI) Login(ctx context.Context, username, password string) (*LoginResponse, error) {
	u := urlutil.CopyWithQueries(a.serviceURL, map[string]string{
		"username": username,
		"password": password,
		"service":  "moodle_mobile_app",
	})
	u.Path = path.Join(u.Path, "/login/token.php")
	res := LoginResponse{}
	if err := a.getAndUnmarshal(ctx, u, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
