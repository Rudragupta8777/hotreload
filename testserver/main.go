package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	// This simulates a server starting up
	log.Println("Test server is starting up...")
	time.Sleep(1 * time.Second) 
	log.Println("Test server is running on http://localhost:8080")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// We will change this text to test the hot reload!
		fmt.Fprintf(w, "Hello from Version 2!") 
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}