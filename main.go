package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	httpFlag = flag.String("http", "localhost:8080", "HTTP listen address")
)

func handleHelloRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World\n")
}

func main() {
	flag.Parse()
	log.Printf("Listening on %s\n", *httpFlag)

	http.HandleFunc("/hello", handleHelloRoute)
	log.Fatal(http.ListenAndServe(*httpFlag, nil))
}
