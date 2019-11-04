package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

func (app *ServeCmd) HandleAuthToken(s *sessions.Session, w http.ResponseWriter, r *http.Request) error {
	b, ok := s.Values["token"].([]byte)
	if !ok {
		return &HTTPServerError{Code: http.StatusUnauthorized, Text: "Unauthorized"}
	}
	var token *oauth2.Token
	err := json.Unmarshal(b, &token)
	if err != nil {
		return err
	}
	ctx := r.Context()
	token, err = app.OAuth2Config.TokenSource(ctx, token).Token()
	if err != nil {
		return err
	}
	b, err = json.Marshal(token)
	if err != nil {
		return err
	}
	s.Values["token"] = b
	err = s.Save(r, w)
	if err != nil {
		return err
	}
	http.Error(w, string(b), http.StatusOK)
	return nil
}

func (app *ServeCmd) HandleAuthCallback(s *sessions.Session, w http.ResponseWriter, r *http.Request) error {
	nonce, ok := s.Values["nonce"].(string)
	if !ok {
		return &HTTPServerError{Code: http.StatusBadRequest, Text: "No state cookie present"}
	}
	delete(s.Values, "nonce")
	state := r.FormValue("state")
	if state != nonce {
		return &HTTPServerError{Code: http.StatusBadRequest, Text: "State verification failed"}
	}
	errorString := r.FormValue("error")
	if errorString != "" {
		errorDescription := r.FormValue("error_description")
		text := fmt.Sprintf("%s\n%s\n", errorString, errorDescription)
		return &HTTPServerError{Code: http.StatusBadRequest, Text: text}
	}
	code := r.FormValue("code")
	ctx := r.Context()
	token, err := app.OAuth2Config.Exchange(ctx, code)
	if err != nil {
		return err
	}
	b, err := json.Marshal(token)
	if err != nil {
		return err
	}
	s.Values["token"] = b
	err = s.Save(r, w)
	if err != nil {
		return err
	}
	url, ok := s.Values["redirect"].(string)
	if ok && url != "" {
		http.Redirect(w, r, url, http.StatusFound)
	} else {
		http.Error(w, "Signed in", http.StatusOK)
	}
	return nil
}

func (app *ServeCmd) HandleAuthSignIn(s *sessions.Session, w http.ResponseWriter, r *http.Request) error {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return err
	}
	nonce := fmt.Sprintf("%x", b)
	s.Values["nonce"] = nonce
	s.Values["redirect"] = r.FormValue("redirect")
	err = s.Save(r, w)
	if err != nil {
		return err
	}
	http.Redirect(w, r, app.OAuth2Config.AuthCodeURL(nonce), http.StatusFound)
	return nil
}

func (app *ServeCmd) HandleAuthSignOut(s *sessions.Session, w http.ResponseWriter, r *http.Request) error {
	delete(s.Values, "token")
	err := s.Save(r, w)
	if err != nil {
		return err
	}
	url := r.FormValue("redirect")
	if url != "" {
		http.Redirect(w, r, url, http.StatusFound)
	} else {
		http.Error(w, "Signed out", http.StatusOK)
	}
	return nil
}

func (app *ServeCmd) HandleAuth(w http.ResponseWriter, r *http.Request) {
	s, _ := app.SessionStore.Get(r, "auth-session")
	err := r.ParseForm()
	if err == nil {
		switch r.URL.Path {
		case "/auth/sign_in":
			err = app.HandleAuthSignIn(s, w, r)
		case "/auth/sign_out":
			err = app.HandleAuthSignOut(s, w, r)
		case "/auth/callback":
			err = app.HandleAuthCallback(s, w, r)
		case "/auth/token":
			err = app.HandleAuthToken(s, w, r)
		default:
			text := fmt.Sprintf("Unhandled path: %s", r.URL.Path)
			err = &HTTPServerError{Code: http.StatusNotFound, Text: text}
		}
	}
	if err != nil {
		herr, ok := err.(*HTTPServerError)
		if !ok {
			if os.IsNotExist(err) {
				herr = &HTTPServerError{Code: http.StatusNotFound, Text: err.Error()}
			} else if os.IsPermission(err) {
				herr = &HTTPServerError{Code: http.StatusForbidden, Text: err.Error()}
			} else {
				herr = &HTTPServerError{Code: http.StatusInternalServerError, Text: err.Error()}
			}
		}
		s.Save(r, w)
		http.Error(w, herr.Error(), herr.Code)
	}
}
