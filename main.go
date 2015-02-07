package main

import "fmt"
import "net/http"


func handleHelloRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World\n")
}

func main() {
	http.HandleFunc("/hello", handleHelloRoute)
	http.ListenAndServe(":8080", nil)
}
