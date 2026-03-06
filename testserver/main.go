package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	log.Println("Test server is starting up...")

	// This simulates a background worker (like a DB connection or job queue)
	// If our kill logic fails, this will keep printing and become a "zombie"
	go func() {
		counter := 1
		for {
			log.Printf("[Worker] Background task running... %d\n", counter)
			time.Sleep(2 * time.Second)
			counter++
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello Trademarkia! ")
	})

	log.Println("Test server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}