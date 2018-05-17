package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	opts   options
	logger *log.Logger
)

type options struct {
	host    string
	port    int
	root    string
	cert    string
	key     string
	tls     bool
	useauth bool
	pass    string
}

func main() {
	logger = log.New(os.Stdout, "", log.LstdFlags)

	flag.StringVar(&opts.host, "host", "localhost", "hostname to listen on")
	flag.IntVar(&opts.port, "port", 8080, "port number to listen on")
	flag.StringVar(&opts.root, "root", "./", "folder which we serve")
	flag.StringVar(&opts.cert, "cert", "server.crt", "the filename of the server certificate")
	flag.StringVar(&opts.key, "key", "server.key", "the filename of the server key")
	flag.BoolVar(&opts.tls, "tls", false, "wether to use tls")
	flag.BoolVar(&opts.useauth, "useauth", true, "use user/pass or not. default: false")
	flag.StringVar(&opts.pass, "pass", "tomatensaft", "the key to enter the directory")
	flag.Parse()

	// Create AuthManager
	authmanager := NewAuthManager(opts.useauth, logger)
	authmanager.AddPassphrase(opts.pass)

	mux := http.NewServeMux()
	fileServer := NewStreamFileSystem(opts.root, "/files")
	mux.Handle("/files/", authmanager.Check(fileServer))
	mux.HandleFunc("/stream/", fileServer.Stream)
	mux.HandleFunc("/login", authmanager.Login)

	url := fmt.Sprintf("%s:%d", opts.host, opts.port)
	logger.Printf("Listening on %s\n", url)
	if !opts.tls {
		logger.Fatal(http.ListenAndServe(url, mux))
	} else {
		logger.Fatal(http.ListenAndServeTLS(url, opts.cert, opts.key, mux))
	}
}
