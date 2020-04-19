package main

import (
	"flag"
	"log"
	"net/http"
)

type commandLineConfig struct {
	buildPath string
	hostport  string
}

func main() {
	config := commandLineConfig{}
	flag.StringVar(&config.buildPath, "build", "web/build", "Path to built front end html and javascript")
	flag.StringVar(&config.hostport, "hostport", "localhost:8080", "Port to start server on")
	flag.Parse()

	log.Printf("Using build path %s\n", config.buildPath)
	fs := http.FileServer(http.Dir(config.buildPath))
	mux := http.NewServeMux()
	mux.Handle("/", fs)

	serve := http.Server{
		Addr:    config.hostport,
		Handler: mux,
	}

	log.Printf("Listening on %s\n", config.hostport)
	log.Fatal(serve.ListenAndServe())
}
