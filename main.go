package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	cookievalue string
	noauth      bool
)

func redirect(w http.ResponseWriter, req *http.Request) {
    // remove/add not default ports from req.Host
    target := "https://" + req.Host + req.URL.Path 
    if len(req.URL.RawQuery) > 0 {
        target += "?" + req.URL.RawQuery
    }
    http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}

func main() {
	host := flag.String("host", "", "hostname to listen on")
	port := flag.Int("port", 8080, "port number to listen on")
	root := flag.String("root", ".", "folder which we serve")
	pass := flag.String("pass", generateKey(), "the key to enter the directory")
	cert := flag.String("cert", "server.crt", "the filename of the server certificate")
	key := flag.String("key", "server.key", "the filename of the server key")
	auth := flag.Bool("noauth", false, "if true == no auth")
	flag.Parse()

	// TODO: meeh, we do not want global vars
	cookievalue = hash(*pass)
	noauth = *auth

	fileServer := http.FileServer(http.Dir(*root))
	http.Handle("/", authHandler(*pass, fileServer))
	http.HandleFunc("/login", loginHandler)

	url := fmt.Sprintf("%s:%d", *host, *port)
	go log.Fatal(http.ListenAndServe(":80", http.HandlerFunc(redirect)))
	log.Fatal(http.ListenAndServeTLS(url, *cert, *key, nil))
}

// TODO: rename filenam on request
