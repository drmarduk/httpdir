package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	cookievalue string // holds the hash of the auth password in browser cookie
	noauth      bool   // wether to use a cookie based auth login
)

func main() {
	// main http conf
	host := flag.String("host", "", "hostname to listen on")
	port := flag.Int("port", 8080, "port number to listen on")
	root := flag.String("root", ".", "folder which we serve")
	// auth stuff
	pass := flag.String("pass", generateKey(), "the key to enter the directory")
	auth := flag.Bool("noauth", false, "if true == no auth")
	// TLS stuff
	cert := flag.String("cert", "server.crt", "the filename of the server certificate")
	key := flag.String("key", "server.key", "the filename of the server key")
	tls := flag.Bool("tls", true, "wether to use tls")
	flag.Parse()

	// TODO: meeh, we do not want global vars
	cookievalue = hash(*pass)
	noauth = *auth

	fileServer := http.FileServer(http.Dir(*root))
	http.Handle("/", authHandler(*pass, fileServer))
	http.HandleFunc("/login", loginHandler)

	url := fmt.Sprintf("%s:%d", *host, *port)

	if !*tls {
		log.Fatal(http.ListenAndServe(url, nil))
	} else {
		log.Fatal(http.ListenAndServeTLS(url, *cert, *key, nil))
	}
}

// TODO: rename filenam on request
