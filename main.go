package main

import (
	"crypto/sha256"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net/http"
)

func generateKey() string {
	return "this.is.sparta"
}

func hash(s string) string {
	hasher := sha256.New()
	hasher.Write([]byte(s))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}

func main() {
	host := flag.String("host", "", "hostname to listen on")
	port := flag.Int("port", 8080, "port number to listen on")
	tls := flag.Bool("tls", true, "should we use tls")
	root := flag.String("root", ".", "folder which we serve")
	pass := flag.String("pass", generateKey(), "the key to enter the directory")
	cert := flag.String("cert", "server.crt", "the filename of the server certificate")
	key := flag.String("key", "server.key", "the filename of the server key")
	flag.Parse()

	cookievalue = hash(*pass)
	log.Printf("Pass: %s - %s\n", *pass, cookievalue)
	fileServer := http.FileServer(http.Dir(*root))
	http.Handle("/", authHandler(*pass, fileServer))
	http.HandleFunc("/login", loginHandler)

	url := fmt.Sprintf("%s:%d", *host, *port)

	if *tls {
		log.Fatal(http.ListenAndServeTLS(url, *cert, *key, nil))
	} else {
		log.Fatal(http.ListenAndServe(url, nil))
	}
}

var cookiename string = "keykkk"
var cookievalue string = ""

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		key := r.Form.Get("key")
		_key := hash(key)
		if _key == cookievalue {
			c := &http.Cookie{}
			c.HttpOnly = true
			c.Name = cookiename
			c.Secure = true
			c.Value = cookievalue
			log.Printf("SetCookie: %s\n", r.RemoteAddr)
			http.SetCookie(w, c)
		}
		http.Redirect(w, r, "/", 301)
		return
	}
	html := `<html><head><title>Login</title><body><form method="post" action="/login">Key<input name="key" type="text"/><input type="submit" /></form></body></html>`
	fmt.Fprint(w, html)
}

func check(r *http.Request) bool {
	log.Printf("check %d cookies\n", len(r.Cookies()))
	for _, c := range r.Cookies() {
		if c.Name == cookiename {
			if c.Value == cookievalue {
				return true
			}
		}
	}
	return false
}

func authHandler(key string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !check(r) {
			http.Redirect(w, r, "/login", 301)
			return
		}
		logger(r)
		r.Header.Get("")
		h.ServeHTTP(w, r)
	})
}

func logger(r *http.Request) {
	fmt.Printf("%s: %s - %s\n", r.Method, r.URL.Path, r.RemoteAddr)
}
