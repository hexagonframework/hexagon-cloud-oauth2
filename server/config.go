package server

import (
	"time"

	"github.com/hexagon-cloud/oauth2"
)

// Config configuration parameters
type Config struct {
	TokenType             string                // token type
	AllowGetAccessRequest bool                  // to allow GET requests for the token
	AllowedResponseTypes  []oauth2.ResponseType // allow the authorization type
	AllowedGrantTypes     []oauth2.GrantType    // allow the grant type
}

// NewConfig create to configuration instance
func NewConfig() *Config {
	return &Config{
		TokenType:            "Bearer",
		AllowedResponseTypes: []oauth2.ResponseType{oauth2.CodeRsp, oauth2.TokenRsp},
		AllowedGrantTypes: []oauth2.GrantType{
			oauth2.AuthorizationCode,
			oauth2.PasswordCredentials,
			oauth2.ClientCredentials,
			oauth2.RefreshToken,
		},
	}
}

// AuthorizeRequest authorization request
type AuthorizeRequest struct {
	ResponseType   oauth2.ResponseType
	ClientID       string
	Scope          string
	RedirectURI    string
	State          string
	UserID         string
	AccessTokenExp time.Duration
}
