package oauth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"VDM2-BankBE/internal/config"
)

// GoogleUserInfo represents the user info we get from Google
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// GoogleOAuthClient handles OAuth2 flow with Google
type GoogleOAuthClient struct {
	config *oauth2.Config
}

// NewGoogleOAuthClient creates a new Google OAuth client
func NewGoogleOAuthClient(cfg *config.OAuthProviderConfig) *GoogleOAuthClient {
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	return &GoogleOAuthClient{
		config: oauthConfig,
	}
}

// GetAuthURL returns the URL to redirect the user to for authorization
func (g *GoogleOAuthClient) GetAuthURL(state string) string {
	return g.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// Exchange exchanges an authorization code for tokens
func (g *GoogleOAuthClient) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return nil, errors.Wrap(err, "failed to exchange Google OAuth code")
	}
	return token, nil
}

// GetUserInfo fetches user info from Google using the access token
func (g *GoogleOAuthClient) GetUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
	client := g.config.Client(ctx, token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch Google user info")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Google API returned status: %s", resp.Status)
	}

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, errors.Wrap(err, "failed to decode Google user info")
	}

	return &userInfo, nil
}
