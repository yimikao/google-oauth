package handlers

import (
	"github.com/yimikao/googleoauth/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func NewGoogleOauthConfig() *oauth2.Config {

	cfg, err := utils.LoadConfig(".")
	if err != nil {
		return nil
	}

	return &oauth2.Config{
		ClientID:     cfg.GoogleOauthClientID,
		ClientSecret: cfg.GoogleOauthClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  cfg.CallbackURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	}
}

var googleOauthConfig = NewGoogleOauthConfig()

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
