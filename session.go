package outlook

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// Session manages communication to microsoft's graph api as an authenticated user.
type Session struct {
	client       *Client
	basePath     string
	accessToken  string
	refreshToken string
}

// NewSession returns a new instance of a Session.
func NewSession(client *Client, refreshToken string) (*Session, error) {
	session := &Session{
		client:       client,
		basePath:     "/me",
		refreshToken: refreshToken,
	}

	if err := session.refreshAccessToken(); err != nil {
		return nil, err
	}

	return session, nil
}

func (session *Session) query(ctx context.Context, method, url string, params map[string]interface{}, data interface{}, result interface{}) (*http.Response, error) {
	var queryString string
	if params != nil {
		queryString = createQueryString(params)
	}

	path := fmt.Sprintf("%s%s%s", session.basePath, url, queryString)

	req, err := session.client.NewRequest(ctx, method, path, data)
	if err != nil {
		return nil, err
	}

	if session.accessToken == "" {
		return nil, ErrNoAccessToken
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", session.accessToken))

	// May want to detect failures due to invalid or expired tokens, then retry after attempting to refresh the token
	return session.client.Do(ctx, req, result)
}

// Get performs a get request to microsofts api with the underlying client and the sessions accessToken for authorization.
func (session *Session) Get(ctx context.Context, url string, params map[string]interface{}, result interface{}) (*http.Response, error) {
	return session.query(ctx, http.MethodGet, url, params, nil, result)
}

// Post performs a post request to microsofts api with the underlying client and the sessions accessToken for authorization.
func (session *Session) Post(ctx context.Context, url string, data interface{}, result interface{}) (*http.Response, error) {
	return session.query(ctx, http.MethodPost, url, nil, data, result)
}

// Patch performs a patch request to microsofts api with the underlying client and the sessions accessToken for authorization.
func (session *Session) Patch(ctx context.Context, url string, data interface{}, result interface{}) (*http.Response, error) {
	return session.query(ctx, http.MethodPatch, url, nil, data, result)
}

// Delete performs a delete request to microsofts api with the underlying client and the sessions accessToken for authorization.
func (session *Session) Delete(ctx context.Context, url string, params map[string]interface{}, result interface{}) (*http.Response, error) {
	return session.query(ctx, http.MethodDelete, url, params, nil, result)
}

// Calendars returns an instance of a CalendarService using this session.
func (session *Session) Calendars() *CalendarService {
	return NewCalendarService(session)
}

// Events returns an instance of a EventService using this session.
func (session *Session) Events() *EventService {
	return NewEventService(session)
}

// Folders returns an instance of a FolderService using this session.
func (session *Session) Folders() *FolderService {
	return NewFolderService(session)
}

// Messages returns an instance of a MessageService using this session.
func (session *Session) Messages() *MessageService {
	return NewMessageService(session)
}

func (session *Session) refreshAccessToken() error {
	body := url.Values{}
	body.Set("client_id", session.client.appID)
	body.Set("client_secret", session.client.appSecret)
	body.Set("refresh_token", session.refreshToken)
	body.Set("redirect_uri", session.client.redirectURI)
	body.Set("scope", session.client.scope)
	body.Set("grant_type", "refresh_token")

	// I suppose it's possible for this to change the media type of other requests being made at essentially the same time, will update sometime
	session.client.SetMediaType("application/x-www-form-urlencoded")

	req, err := session.client.NewRequest(context.Background(), http.MethodPost, DefaultOAuthTokenURL, body)
	if err != nil {
		return err
	}

	session.client.SetMediaType(mediaType)

	var tokenRes RefreshTokenResponse
	if _, err := session.client.Do(context.Background(), req, &tokenRes); err != nil {
		return err
	}

	session.accessToken = tokenRes.AccessToken

	return nil
}
