package client

// Client details interface for OAuth 2
type ClientDetails interface {
	GetId() string
	GetSecret() string
	GetDomain() string
	GetUserID() string
}

// Client client model
type Client struct {
	ID     string
	Secret string
	Domain string
	UserID string
}

// GetID client id
func (c *Client) GetID() string {
	return c.ID
}

// GetSecret client domain
func (c *Client) GetSecret() string {
	return c.Secret
}

// GetDomain client domain
func (c *Client) GetDomain() string {
	return c.Domain
}

// GetUserID user id
func (c *Client) GetUserID() string {
	return c.UserID
}