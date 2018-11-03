package main

import (
	"fmt"
	"log"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello world</h1>")
}

func main() {
	http.HandleFunc("/", hello)
	log.Fatal(http.ListenAndServe(":9911", nil))
}
