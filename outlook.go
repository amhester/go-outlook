package outlook

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// ClientVersion the current version of this sdk
	ClientVersion = "0.1.0"
	// DefaultBaseURL the root host url for the microsoft outlook api
	DefaultBaseURL = "https://graph.microsoft.com/v1.0"
	// DefaultOAuthTokenURL the url used to exchange a user's refreshToken for a usable accessToken
	DefaultOAuthTokenURL = "https://login.microsoftonline.com/common/oauth2/v2.0/token"
	// DefaultAuthScopes the set of permissions the client will request from the user
	DefaultAuthScopes = "mail.read calendars.read user.read offline_access"
	// DefaultQueryDateTimeFormat time format for the datetime query parameters used in outlook
	DefaultQueryDateTimeFormat = "2006-01-02T15:04:05Z"

	mediaType = "application/json"
)

var (
	// ErrNoDeltaLink error when our email paging fails to return a delta token at the end.
	ErrNoDeltaLink = errors.New("no delta link on response")

	// DefaultClient the http client that the sdk will use to make calls.
	DefaultClient = &http.Client{Timeout: time.Second * 60}

	// DefaultUserAgent the user agent to get passed in request headers on each call
	DefaultUserAgent = fmt.Sprintf("go-outlook/%s", ClientVersion)
)

// Client manages communication with microsoft's graph api, specifically for Mail and Calendar.
type Client struct {
	client      *http.Client
	baseURL     *url.URL
	userAgent   string
	mediaType   string
	appID       string
	appSecret   string
	redirectURI string
	scope       string
}

// ClientOpt functions to configure options on a Client.
type ClientOpt func(*Client)

// SetClientAppID returns a ClientOpt function which set the clients AppID.
func SetClientAppID(appID string) ClientOpt {
	return func(c *Client) {
		c.appID = appID
	}
}

// SetClientAppSecret returns a ClientOpt function which set the clients App Secret.
func SetClientAppSecret(secret string) ClientOpt {
	return func(c *Client) {
		c.appSecret = secret
	}
}

// SetClientRedirectURI returns a ClientOpt function which set the clients redirectURI.
func SetClientRedirectURI(uri string) ClientOpt {
	return func(c *Client) {
		c.redirectURI = uri
	}
}

// SetClientScope returns a ClientOpt function which set the clients auth scope.
func SetClientScope(scope string) ClientOpt {
	return func(c *Client) {
		c.scope = scope
	}
}

// SetClientMediaType returns a ClientOpt function which sets the clients mediaType.
func SetClientMediaType(mType string) ClientOpt {
	return func(c *Client) {
		c.mediaType = mType
	}
}

// NewClient returns a new instance of a Client with the given options set.
func NewClient(opts ...ClientOpt) (*Client, error) {
	baseURL, err := url.Parse(DefaultBaseURL)
	if err != nil {
		return nil, err
	}
	client := &Client{
		client:    DefaultClient,
		baseURL:   baseURL,
		userAgent: DefaultUserAgent,
		scope:     DefaultAuthScopes,
		mediaType: mediaType,
	}
	for _, opt := range opts {
		opt(client)
	}
	return client, nil
}

// SetAppID fluent configuration of the client's microsoft AppID.
func (client *Client) SetAppID(appID string) *Client {
	client.appID = appID
	return client
}

// SetAppSecret fluent configuration of the client's microsoft App Secret.
func (client *Client) SetAppSecret(secret string) *Client {
	client.appSecret = secret
	return client
}

// SetRedirectURI fluent configuration of the client's microsoft RedirectURI.
func (client *Client) SetRedirectURI(uri string) *Client {
	client.redirectURI = uri
	return client
}

// SetScope fluent configuration of the client's microsoft auth Scope.
func (client *Client) SetScope(scope string) *Client {
	client.scope = scope
	return client
}

// SetMediaType fluent configuration of the client's mediaType.
func (client *Client) SetMediaType(mType string) *Client {
	client.mediaType = mType
	return client
}

// NewRequest creates a new request with some reasonable defaults based on the client.
func (client *Client) NewRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	var fullURL string
	pathURL, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	if pathURL.Hostname() != "" {
		fullURL = path
	} else {
		fullURL = fmt.Sprintf("%s%s", client.baseURL.String(), path)
	}

	encodedBody := new(bytes.Buffer)
	if body != nil {
		switch client.mediaType {
		case "application/json":
			if err := json.NewEncoder(encodedBody).Encode(body); err != nil {
				return nil, err
			}
		case "application/x-www-form-urlencoded":
			if v, ok := body.(url.Values); ok {
				bodyReader := strings.NewReader(v.Encode())
				if _, err := io.Copy(encodedBody, bodyReader); err != nil {
					return nil, err
				}
			} else {
				return nil, fmt.Errorf("Body must be of type url.Values when Content-Type is set to application/x-www-form-urlencoded")
			}
		}
	}

	req, err := http.NewRequest(method, fullURL, encodedBody)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", client.mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", client.userAgent)

	return req, nil
}

// Do executes the given http request and will bind the response body with v. Returns the http response as well as any error.
func (client *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)
	response, err := client.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := response.Body.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	err = checkResponse(response)
	if err != nil {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, response.Body)
			if err != nil {
				return response, err
			}
		} else {
			err = json.NewDecoder(response.Body).Decode(v)
			if err != nil {
				return response, err
			}
		}
	}

	return response, err
}

// NewSession returns a new instance of a Session using this client and the given refreshToken.
func (client *Client) NewSession(refreshToken string) (*Session, error) {
	session, err := NewSession(client, refreshToken)
	if err != nil {
		return nil, err
	}
	return session, nil
}
