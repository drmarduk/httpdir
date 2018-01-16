package main

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
)

// AuthManager handles all requests and look if a user can do requests
type AuthManager struct {
	UseAuth bool

	Passphrase string            // single pass login
	Users      map[string]string // user:pass login

	logger *log.Logger
}

// NewAuthManager returns a new http handler for checking cookievalues
func NewAuthManager(useauth bool, logger *log.Logger) *AuthManager {
	am := &AuthManager{
		UseAuth: useauth,
		logger:  logger,
	}

	am.logger.Printf("NewAuthManager: useAuth: %v\n", am.UseAuth)
	return am
}

// AddUser enables a user-based login
func (am *AuthManager) AddUser(user, pass string) {
	if !am.UseAuth {
		return
	}
	am.Users[user] = am.hashstring(pass)
	am.logger.Printf("Added User %s:%s\n", user, am.Users[user])
}

// AddPassphrase enables a single pass login
func (am *AuthManager) AddPassphrase(pass string) {
	if !am.UseAuth {
		return
	}
	am.Passphrase = am.hashstring(pass)
	am.logger.Printf("Added Passphrase %s\n", am.Passphrase)
}

// Check is the http handler which checks the request if the key cookie is set
// can be expanded to support user based cookies
func (am *AuthManager) Check(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if am.UseAuth {
			c, err := r.Cookie("passphrase")
			if err != nil {
				http.Error(w, "not authorized", http.StatusUnauthorized)
				return
			}

			if c.Value != am.Passphrase {
				http.Error(w, "wrong passphrase", http.StatusUnauthorized)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

// Login is our main handler for log in
func (am *AuthManager) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "%s", htmlLogin)
		return
	}

	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Error while parsing /login form: %v\n", err)
		return
	}

	pass := r.Form.Get("passphrase")
	hashed := am.hashstring(pass)

	if am.checkpassphrase(hashed) {
		c := &http.Cookie{
			Name:  "passphrase",
			Value: am.Passphrase,
		}
		http.SetCookie(w, c)
		am.logger.Printf("/Login: added Cookie %v\n", c)
	}
}

func (am *AuthManager) hashstring(s string) string {
	hasher := sha512.New()
	hasher.Write([]byte(s))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}

func (am *AuthManager) checkpassphrase(pass string) bool {
	if am.Passphrase == pass {
		return true
	}
	return false
}
