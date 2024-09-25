package main

import (
	"fmt"
	"net/http"

	"github.com/thenets/notes/kvstore"
)

func pageApiNoteGet(w http.ResponseWriter, r *http.Request) {
    // Extract the key from the URL path
    keys := r.URL.Path[len("/api/"):]
    if keys == "" {
        http.Error(w, "Key is required", http.StatusBadRequest)
        return
    }

    // Print the key to the console
    fmt.Println("Received key:", keys)

    // Send a response back to the client
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Key received: %s", keys)
}


func pageHomepage(w http.ResponseWriter, r *http.Request) {
	// Create a new instance of InMemoryKVStore
	kv := &kvstore.InMemoryKVStore{}

	// Set a value
	err := kv.Set("name", "Alice")
	if err != nil {
		fmt.Println("Error setting key:", err)
		return
	}

	// Get the value
	value, err := kv.Get("name")
	if err != nil {
		fmt.Println("Error getting key:", err)
		return
	}
	fmt.Println("Value for 'name':", value)

	http.ServeFile(w, r, "./static/index.html")
}

func main() {
	http.HandleFunc("/", pageHomepage)
	http.HandleFunc("/api/", pageApiNoteGet)
	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
