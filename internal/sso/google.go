package sso

import (
	"context"

	"github.com/tapiaw38/auth-api/internal/models"
	"golang.org/x/oauth2"
	googleauth "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type GoogleClient struct {
	ClientID     string
	ClientSecret string
	FrontendURL  string
}

func NewGoogleClient(config *GoogleClient) *GoogleClient {
	return &GoogleClient{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		FrontendURL:  config.FrontendURL,
	}
}

func (g *GoogleClient) GoogleClientInit() *oauth2.Config {
	oauth := &oauth2.Config{
		ClientID:     g.ClientID,
		ClientSecret: g.ClientSecret,
		RedirectURL:  g.FrontendURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
	}
	return oauth
}

// ExchangeCode exchanges the code for a token
func (g *GoogleClient) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := g.GoogleClientInit().Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// GetUserInfo gets the user info from the token
func (g *GoogleClient) GetUserInfo(ctx context.Context, token *oauth2.Token) (*models.SocialUserData, error) {

	svc, err := googleauth.NewService(ctx, option.WithTokenSource(oauth2.StaticTokenSource(token)))
	if err != nil {
		return nil, err
	}

	userInfo, err := svc.Userinfo.Get().Do()
	if err != nil {
		return nil, err
	}

	return &models.SocialUserData{
		Email:         userInfo.Email,
		FirstName:     userInfo.GivenName,
		LastName:      userInfo.FamilyName,
		Picture:       userInfo.Picture,
		VerifiedEmail: true,
	}, nil
}
