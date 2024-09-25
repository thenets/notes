package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/thenets/notes/kvstore"
)

var kv = &kvstore.InMemoryKVStore{}

type Note struct {
	Content   string `json:"content"`
	UpdatedAt string `json:"updated_at"`
	Error     string `json:"error"`
}

func getCurrentDatetime() string {
	now := time.Now().UTC()
	iso8601Format := now.Format(time.RFC3339)
	return iso8601Format
}

func pageApiNoteGet(w http.ResponseWriter, r *http.Request) {
	// Extract the key from the URL path
	keys := r.URL.Path[len("/api/"):]
	if keys == "" {
		http.Error(w, "Key is required", http.StatusBadRequest)
		return
	}
	if strings.Count(keys, "/") > 0 {
		http.Error(w, "Invalid key format. Key must be in the form /api/{note_id}", http.StatusBadRequest)
		return
	}

	// Update value
	fmt.Println("RETRIEVE:", keys)
	before_value, err := kv.Get(keys)
	if err != nil {
		fmt.Println("Error getting key:", err)
	}
	if err != nil {
		fmt.Println("Setting empty value to it.")
		before_value = ""
	}
	fmt.Println("  - before:", before_value)
	after_value := before_value + "+"

	err = kv.Set(keys, after_value)
	if err != nil {
		fmt.Println("Error setting key:", err)
		return
	}
	fmt.Println("  - after:", after_value)

	// Send a response back to the client
	note := Note{
		Content:   after_value,
		UpdatedAt: getCurrentDatetime(),
		Error:     "",
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(note)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode JSON: %v", err), http.StatusInternalServerError)
		return
	}
}

func staticFileHandler(w http.ResponseWriter, r *http.Request) {
	keys := r.URL.Path[len("/"):]

	var keys_slice []string

	// Static pattern identified
	if len(keys) == 0 || strings.HasPrefix(keys, "assets/") {
		// Default for homepage
		if len(keys) == 0 {
			keys_slice = append(keys_slice, "static")
			keys_slice = append(keys_slice, "index.html")
		}

		// Assets
		if strings.HasPrefix(keys, "assets/") {
			keys_slice = append(keys_slice, "static")
			keys_slice = append(keys_slice, strings.Split(keys, "/")...)
		}

		filePath := filepath.Join(keys_slice...)
		fmt.Println(filePath)

		// Check if the file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			http.Error(w, "Not Found1", http.StatusNotFound)
			return
		}

		// Serve the static file
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, filePath)
	}

	// Loading notes UI
	if len(strings.Split(keys, "/")) > 1 {
		http.Error(w, "Not Found2", http.StatusNotFound)
		return
	}
	keys_slice = append(keys_slice, "static")
	keys_slice = append(keys_slice, "note.html")
	filePath := filepath.Join(keys_slice...)
	fmt.Println(filePath)
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, filePath)

}

// func pageHomepage(w http.ResponseWriter, r *http.Request) {
// 	http.ServeFile(w, r, "./static/index.html")
// }

func main() {
	http.HandleFunc("/api/", pageApiNoteGet)
	http.HandleFunc("/", staticFileHandler)
	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
