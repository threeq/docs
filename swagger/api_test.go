package swagger

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestEndpoints_ServeHTTPNotFound(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://localhost", nil)
	w := httptest.NewRecorder()

	e := Endpoints{}
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestFilepathJoin(t *testing.T) {
	assert.Equal(t, "/api", filepath.Join("/", "/api"))
	assert.Equal(t, "/", filepath.Join("/", "/"))
}

func TestEndpoints_ServeHTTP(t *testing.T) {
	fn := func(v string) *Endpoint {
		return &Endpoint{
			Handler: func(w http.ResponseWriter, req *http.Request) {
				w.WriteHeader(http.StatusOK)
				io.WriteString(w, v)
			},
		}
	}

	e := Endpoints{
		Delete:  fn("Delete"),
		Head:    fn("Head"),
		Get:     fn("Get"),
		Options: fn("Options"),
		Post:    fn("Post"),
		Put:     fn("Put"),
		Patch:   fn("Patch"),
		Trace:   fn("Trace"),
		Connect: fn("Connect"),
	}

	methods := []string{
		http.MethodDelete,
		http.MethodHead,
		http.MethodGet,
		http.MethodOptions,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodTrace,
		http.MethodConnect,
	}
	for _, method := range methods {
		req, err := http.NewRequest(strings.ToUpper(method), "http://localhost", nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		e.ServeHTTP(w, req)
		assert.Equal(t, strings.ToUpper(w.Body.String()), strings.ToUpper(method))
	}
}

func TestSecuritySchemeDescription(t *testing.T) {
	scheme := &SecurityScheme{}
	description := "a security scheme"
	SecuritySchemeDescription(description)(scheme)
	assert.Equal(t, description, scheme.Description)
}

func TestBasicSecurity(t *testing.T) {
	scheme := &SecurityScheme{}
	BasicSecurity()(scheme)
	assert.Equal(t, scheme.Type, "basic")
}

func TestAPIKeySecurity(t *testing.T) {
	scheme := &SecurityScheme{}
	name := "Authorization"
	in := "header"
	APIKeySecurity(name, in)(scheme)
	assert.Equal(t, scheme.Type, "apiKey")
	assert.Equal(t, scheme.Name, name)
	assert.Equal(t, scheme.In, in)

	assert.Panics(t, func() {
		APIKeySecurity(name, "invalid")
	}, "expected APIKeySecurity to panic with invalid \"in\" parameter")
}

func TestOAuth2Security(t *testing.T) {
	scheme := &SecurityScheme{}

	flow := "accessCode"
	authURL := "http://example.com/oauth/authorize"
	tokenURL := "http://example.com/oauth/token"
	OAuth2Security(flow, authURL, tokenURL)(scheme)

	assert.Equal(t, scheme.Type, "oauth2")
	assert.Equal(t, scheme.Flow, "accessCode")
	assert.Equal(t, scheme.AuthorizationURL, authURL)
	assert.Equal(t, scheme.TokenURL, tokenURL)
}

func TestOAuth2Scope(t *testing.T) {
	scheme := &SecurityScheme{}

	OAuth2Scope("read", "read data")(scheme)
	OAuth2Scope("write", "write data")(scheme)

	assert.Len(t, scheme.Scopes, 2)
	assert.Contains(t, scheme.Scopes, "read")
	assert.Contains(t, scheme.Scopes, "write")

	assert.Equal(t, "read data", scheme.Scopes["read"])
	assert.Equal(t, "write data", scheme.Scopes["write"])
}
