package main

import (
	"log"
	"net/http"
	"test-backend-1-ArtyomRytikov/internal/handler"
)

func main() {
	r := handler.NewRouter()

	log.Println("server started on :8080")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}