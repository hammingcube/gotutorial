package main

import (
	"flag"
	"fmt"
	_ "go/build"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/maddyonline/tools/blog"
	"github.com/maddyonline/tools/playground/socket"
	//"golang.org/x/tools/blog"
	//"golang.org/x/tools/playground/socket"
)

const (
	packagePath = "github.com/maddyonline/gotutorial"
)

var (
	httpFlag   = flag.String("http", "localhost:8080", "HTTP listen address")
	originFlag = flag.String("origin", "", "web socket origin for Go Playground (e.g. localhost)")
	baseFlag   = flag.String("base", "", "base path for articles and resources")
)

func main() {
	flag.Parse()

	if *baseFlag == "" {
		// By default, the base is the blog package location.
		*baseFlag = filepath.Join("/go/src", packagePath)
		fmt.Fprintf(os.Stderr, "Trying %s", *baseFlag)
		//os.Exit(1)
	}

	ln, err := net.Listen("tcp", *httpFlag)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	host, port, err := net.SplitHostPort(*httpFlag)
	if err != nil {
		log.Fatal(err)
	}

	srv, err := blog.NewServer(blog.Config{
		ContentPath:  filepath.Join(*baseFlag, "articles"),
		TemplatePath: filepath.Join(*baseFlag, "templates"),
		Hostname:     host,
		HomeArticles: 4,
		FeedArticles: 4,
		FeedTitle:    "Madhav's Blog",
		PlayEnabled:  true,
	})
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	origin := &url.URL{Scheme: "http"}
	if *originFlag != "" {
		origin.Host = net.JoinHostPort(*originFlag, port)
	} else {
		origin.Host = *httpFlag
	}

	http.Handle("/static/", http.FileServer(http.Dir(*baseFlag)))
	http.Handle("/socket", socket.NewHandler(origin))
	http.Handle("/", srv)
	log.Fatal(http.Serve(ln, nil))
}
