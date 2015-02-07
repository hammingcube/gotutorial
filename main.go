package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"go/build"
	"os"
)

const basePkg = "github.com/maddyonline/gotutorial"
var basePath = ""

func handleHelloRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World\n")
}

func initBasePath(basePath *string) {
	if *basePath == "" {
		p, err := build.Default.Import(basePkg, "", build.FindOnly)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't find blog files: %v\n", err)
			fmt.Fprintf(os.Stderr, basePathMessage, basePkg)
			os.Exit(1)
		}
		*basePath = p.Dir
	}
	log.Printf("Using %s as directory for content and static files.", *basePath)
}

func main() {
	httpAddr := flag.String("http", "127.0.0.1:3999", "HTTP service address (e.g., '127.0.0.1:3999')")
	flag.StringVar(&basePath, "base", "", "base path for slide template and static resources")
	flag.Parse()
	initBasePath(&basePath)
	
	http.HandleFunc("/hello", handleHelloRoute)
	log.Printf("Listening on %s\n", *httpAddr)
	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}

const basePathMessage = `
By default, goblog locates the content files and associated
static content by looking for a %q package
in your Go workspaces (GOPATH).
You may use the -base flag to specify an alternate location.
`
