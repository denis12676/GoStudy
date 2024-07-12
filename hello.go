package main

import (
	"log"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!\n"))
}
func main() {
	http.HandleFunc("/", Handler)

	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
