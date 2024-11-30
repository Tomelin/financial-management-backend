package authProvider

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const (
	// AuthProviderGoogle is the google provider
	key    = "randomString"
	MaxAge = 86400 * 30
	IsProd = false
)

type IAuthProvider interface {
	Login(w http.ResponseWriter, r *http.Request) (*SessionStore, error)
	Callback(w http.ResponseWriter, r *http.Request) (*goth.User, error)
	Logout(c *gin.Context, w http.ResponseWriter, r *http.Request) error
	IsLoggedIn(w http.ResponseWriter, r *http.Request) (*goth.User, error)
	Store() (*sessions.CookieStore, error)
}

type SessionRequest struct {
	Request  *http.Request
	Response *http.ResponseWriter
}

type SessionStore struct {
	CookieStore *sessions.CookieStore
	Request     *SessionRequest
	Session     *sessions.Session
	User        *goth.User
}

type AuthConfig struct {
	Providers ProvidersConfig `json:"providers"`
}

type ProviderConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURL  string `json:"redirect_url"`
}

type ProvidersConfig struct {
	Google ProviderConfig `json:"google"`
}

func NewAuthProvider(fields any) (IAuthProvider, error) {

	session := SessionStore{}
	rest, err := session.parseConfig(fields)
	if err != nil {
		return nil, err
	}

	store := sessions.NewCookieStore([]byte(key))
	sessions := &sessions.Options{
		Path:     "/",
		MaxAge:   MaxAge,
		HttpOnly: true,
		Secure:   IsProd,
	}
	store.Options(*sessions)

	gothic.Store = store
	goth.UseProviders(
		google.New(rest.Providers.Google.ClientID, rest.Providers.Google.ClientSecret, rest.Providers.Google.RedirectURL),
	)

	session.CookieStore = &store
	return &session, nil
}

func (a *SessionStore) Login(w http.ResponseWriter, r *http.Request) (*SessionStore, error) {

	return a, nil
}

func (a *SessionStore) Callback(w http.ResponseWriter, r *http.Request) (*goth.User, error) {

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		return nil, err
	}

	session, err := gothic.Store.Get(r, "X-Authorization")
	if err != nil {
		return nil, err
	}

	session.Values["user"] = user
	err = session.Save(r, w)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (a *SessionStore) Logout(c *gin.Context, w http.ResponseWriter, r *http.Request) error {

	session := sessions.Default(c)
	session.Clear()
	session.Options(sessions.Options{
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   IsProd,
	})
	err := session.Save()
	if err != nil {
		return fmt.Errorf("error saving session: %s", err.Error())
	}

	err = gothic.Logout(w, r)
	if err != nil {
		return fmt.Errorf("error during logout: %s", err.Error())
	}

	return nil
}

func (a *SessionStore) IsLoggedIn(w http.ResponseWriter, r *http.Request) (*goth.User, error) {

	session, err := gothic.Store.Get(r, "X-Authorization")
	if err != nil {
		return nil, err
	}

	user, ok := session.Values["user"].(goth.User)
	if !ok {
		return nil, errors.New("user or session not found")
	}
	return &user, nil
}

func (a *SessionStore) Store() (*sessions.CookieStore, error) {
	return a.CookieStore, nil
}

func (a *SessionStore) parseConfig(fields any) (*AuthConfig, error) {
	b, err := json.Marshal(fields)
	if err != nil {
		return nil, err
	}

	var rest AuthConfig
	err = json.Unmarshal(b, &rest)
	if err != nil {
		return nil, err
	}

	return &rest, nil
}
