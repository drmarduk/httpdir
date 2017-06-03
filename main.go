package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var cookievalue string

func main() {
	host := flag.String("host", "", "hostname to listen on")
	port := flag.Int("port", 8080, "port number to listen on")
	root := flag.String("root", ".", "folder which we serve")
	pass := flag.String("pass", generateKey(), "the key to enter the directory")
	cert := flag.String("cert", "server.crt", "the filename of the server certificate")
	key := flag.String("key", "server.key", "the filename of the server key")
	tls := flag.Bool("tls", true ,"wether to use tls")
	flag.Parse()

	// TODO: meeh, we do not want global vars
	cookievalue = hash(*pass)

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
