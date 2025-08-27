package main

import (
	"fmt"
	"math/rand"
	"net/http"
)

func randomHandler(w http.ResponseWriter, r *http.Request) {
	num := rand.Intn(6) + 1 // от 1 до 6

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%d", num)
}

func main() {
	http.HandleFunc("/", randomHandler)

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
