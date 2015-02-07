package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"go/build"
	"net/url"
	"os"
	"path/filepath"
	"golang.org/x/tools/blog"
	"golang.org/x/tools/playground/socket"
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
	basePath := flag.String("base", "", "base path for slide template and static resources")
	originHost := flag.String("orighost", "", "host component of web origin URL (e.g., 'localhost')")
	flag.Parse()
	initBasePath(basePath)

	host, port, err := net.SplitHostPort(*httpAddr)
	if err != nil {
		log.Fatal(err)
	}

	origin := &url.URL{Scheme: "http"}
	if *originHost != "" {
		origin.Host = net.JoinHostPort(*originHost, port)
	} else {
		origin.Host = *httpAddr
	}
	

	srv, err := blog.NewServer(blog.Config{
		ContentPath:  filepath.Join(*basePath, "articles"),
		TemplatePath: filepath.Join(*basePath, "templates"),
		Hostname:     host,
		HomeArticles: 4,
		FeedArticles: 4,
		FeedTitle:    "Madhav's Blog",
		PlayEnabled:  true,
	})
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	
	http.Handle("/static/", http.FileServer(http.Dir(*basePath)))
	http.Handle("/socket", socket.NewHandler(origin))
	http.Handle("/", srv)

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
