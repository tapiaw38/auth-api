package models

// SocialUserData is the data returned from the social login
type SocialUserData struct {
	Token         string `json:"token"`
	RefreshToken  string `json:"refresh_token"`
	Scopes        string `json:"scopes"`
	Email         string `json:"email"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Picture       string `json:"picture"`
	Birthday      string `json:"birthday"`
	VerifiedEmail bool   `json:"verified_email"`
}
