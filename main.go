package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprintf(w, "Hello, GitHub Actions!"); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			fmt.Printf("Error writing response: %v\n", err)
		}
	})
	fmt.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
