package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello world, this is User Service")
	})

	log.Println("Starting server on port 9090...")
	log.Fatal(http.ListenAndServe(":9090", nil))
}
