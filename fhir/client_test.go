package fhir

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestClient_DoWithBasicAuth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")
		if len(authHeader) == 0 || !strings.HasPrefix(authHeader, "Basic") {
			t.FailNow()
		}
	}))
	defer server.Close()

	auth := ClientAuth{BasicAuthUser: "foo", BasicAuthPassword: "bar"}
	baseURL, _ := url.ParseRequestURI(server.URL)
	client := NewClient(*baseURL, auth)

	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
	_, _ = client.Do(req)
}

func TestClient_DoWithBasicAuthWithoutPassword(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")
		if len(authHeader) == 0 || !strings.HasPrefix(authHeader, "Basic") {
			t.FailNow()
		}
	}))
	defer server.Close()

	auth := ClientAuth{BasicAuthUser: "foo", BasicAuthPassword: ""}
	baseURL, _ := url.ParseRequestURI(server.URL)
	client := NewClient(*baseURL, auth)

	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
	_, _ = client.Do(req)
}

func TestClient_DoWithoutBasicAuth(t *testing.T) {
	// we need a handler to check whether the basic auth was NOT set
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")
		if len(authHeader) != 0 {
			t.FailNow()
		}
	}))
	defer server.Close()

	auth := ClientAuth{BasicAuthUser: "", BasicAuthPassword: ""}
	baseURL, _ := url.ParseRequestURI(server.URL)
	client := NewClient(*baseURL, auth)

	req, _ := http.NewRequest("GET", "/", nil)
	_, _ = client.Do(req)
}

func TestNewClient(t *testing.T) {
	t.Run("BaseURL without path", func(t *testing.T) {
		parsedUrl, _ := url.ParseRequestURI("http://localhost:8080")
		client := NewClient(*parsedUrl, ClientAuth{})

		assert.Empty(t, client.baseURL.Path)
	})

	t.Run("BaseURL with path ending without slash", func(t *testing.T) {
		parsedUrl, _ := url.ParseRequestURI("http://localhost:8080/some-path")
		client := NewClient(*parsedUrl, ClientAuth{})

		assert.NotEmpty(t, client.baseURL.Path)
		assert.True(t, strings.HasSuffix(client.baseURL.Path, "some-path/"))
	})

	t.Run("BaseURL with path ending with slash", func(t *testing.T) {
		parsedUrl, _ := url.ParseRequestURI("http://localhost:8080/some-path/")
		client := NewClient(*parsedUrl, ClientAuth{})

		assert.NotEmpty(t, client.baseURL.Path)
		assert.True(t, strings.HasSuffix(client.baseURL.Path, "some-path/"))
	})
}
