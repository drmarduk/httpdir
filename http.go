package main

import (
	"fmt"
	"net/http"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		key := r.Form.Get("key")
		_key := hash(key)
		if _key == cookievalue {
			c := &http.Cookie{
				HttpOnly: true,
				Name:     "key",
				Secure:   true,
				Value:    cookievalue,
			}
			http.SetCookie(w, c)
		}
		http.Redirect(w, r, "/", 301)
	} else {
		html := `<html><head><title>Login</title><body><form method="post" action="/login">Key<input name="key" type="text"/><input type="submit" /></form></body></html>`
		fmt.Fprint(w, html)
	}
}

func check(r *http.Request) bool {
	if noauth {
		return true // OOOOOOHHHH, watch out
	}
	k, err := r.Cookie("key")
	if err != nil {
		return false
	}
	if k.Value == cookievalue {
		return true
	}
	return false
}

func authHandler(key string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s: %s - %s\n", r.Method, r.URL.Path, r.RemoteAddr)
		if !check(r) {
			http.Redirect(w, r, "/login", 301)
			return
		}
		h.ServeHTTP(w, r)
	})
}
