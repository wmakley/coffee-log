package util

import (
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type BasicCredentials interface {
	Username() string
	Password() string
}

func NewBasicCredentials(username string, password string) BasicCredentials {
	return &basicCredentials{
		username: username,
		password: password,
	}
}

type basicCredentials struct {
	username string
	password string
}

func (b *basicCredentials) Username() string {
	return b.username
}

func (b *basicCredentials) Password() string {
	return b.password
}

func NewTestRequest(method string, url string, body io.Reader, credentials BasicCredentials) *http.Request {
	req := httptest.NewRequest(method, url, body)
	req.SetBasicAuth(credentials.Username(), credentials.Password())
	return req
}

func AssertRedirectedTo(t *testing.T, expectedLocation string, status int, response *http.Response) {
	require.Equal(t, status, response.StatusCode)
	actualLoc, err := response.Location()
	require.NoError(t, err)
	require.Equal(t, expectedLocation, actualLoc.String())
}

func FollowRedirect(t *testing.T, redirect *http.Response, credentials BasicCredentials) *http.Request {
	loc, err := redirect.Location()
	require.NoError(t, err)
	req := httptest.NewRequest("GET", loc.String(), nil)
	req.SetBasicAuth(credentials.Username(), credentials.Password())
	return req
}

func ReadAndLogBody(t *testing.T, res *http.Response) string {
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	bodyString := string(body)
	t.Log("Response Body:", bodyString)
	return bodyString
}
