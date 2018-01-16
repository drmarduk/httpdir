package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	cookievalue string
	noauth      bool
)

func main() {
	// Server info
	host := flag.String("host", "", "hostname to listen on")
	port := flag.Int("port", 8080, "port number to listen on")
	root := flag.String("root", ".", "folder which we serve")
	// TLS
	cert := flag.String("cert", "server.crt", "the filename of the server certificate")
	key := flag.String("key", "server.key", "the filename of the server key")
	tls := flag.Bool("tls", false, "wether to use tls")
	// Auth
	useauth := flag.Bool("useauth", false, "weather we want to use pass based auth. default: false")
	pass := flag.String("pass", "tomatensaft", "the key to enter the directory")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.LstdFlags)

	// Create AuthManager
	authmanager := NewAuthManager(*useauth, logger)
	authmanager.AddPassphrase(*pass)

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(*root))
	mux.Handle("/files", authmanager.Check(fileServer))
	// mux.Handle("/stream", authmanager.TokenCheck(fileserver))
	mux.HandleFunc("/login", authmanager.Login)

	url := ""
	if !*tls {
		url = fmt.Sprintf("http://%s:%d", *host, *port)
		logger.Printf("Listening on %s\n", url)
		logger.Fatal(http.ListenAndServe(url, mux))
	} else {
		url = fmt.Sprintf("https://%s:%d", *host, *port)
		logger.Printf("Listening on %s\n", url)
		logger.Fatal(http.ListenAndServeTLS(url, *cert, *key, mux))
	}
}
