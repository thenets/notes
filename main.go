package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"math/rand"

	"github.com/thenets/notes/kvstore"

	"github.com/sirupsen/logrus"
)

var kv = &kvstore.InMemoryKVStore{}

type Note struct {
	Content   string `json:"content"`
	UpdatedAt string `json:"updated_at"`
	Error     string `json:"error"`
}

type ResponseNotePost struct {
	Note string `json:"note"`
}

func getCurrentDatetime() string {
	now := time.Now().UTC()
	iso8601Format := now.Format(time.RFC3339)
	return iso8601Format
}

func generateRandomString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func pageApiNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	keys := r.URL.Path[len("/api/"):]

	if r.Method == http.MethodGet {
		logrus.Debug("GET", " ", keys)

		// Extract the key from the URL path
		if keys == "" {
			logrus.Warn(keys, " ", "Key is required")
			http.Error(w, "Key is required", http.StatusBadRequest)
			return
		}
		if strings.Count(keys, "/") > 0 {
			logrus.Warn(keys, " ", "Invalid key format. Key must be in the form /api/{note_id}")
			http.Error(w, "Invalid key format. Key must be in the form /api/{note_id}", http.StatusBadRequest)
			return
		}

		// Retrieving value
		before_value, err := kv.Get(keys)
		if err != nil {
			before_value = ""
		}
		// fmt.Println("  - before:", before_value)

		// Send a response back to the client
		note := Note{
			Content:   before_value,
			UpdatedAt: getCurrentDatetime(),
			Error:     "",
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(note)
		if err != nil {
			logrus.WithError(err).Error(keys, " ", "Failed to encode JSON")
			http.Error(w, fmt.Sprintf("Failed to encode JSON: %v", err), http.StatusInternalServerError)
			return
		}
	}

	if r.Method == http.MethodPost {
		logrus.Info("POST", " ", keys)

		// Receive the response
		responseBodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			logrus.WithError(err).Error(keys, " ", "Error reading response body")
			http.Error(w, "Error reading response body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		// Declare a variable to hold the parsed JSON data
		var requestData map[string]interface{}

		// Parse as struct
		err = json.Unmarshal(responseBodyBytes, &requestData)
		if err != nil {
			logrus.WithError(err).Error(keys, " ", "Error parsing JSON")
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(requestData)

		dataString := requestData["note"].(string)

		// Update value
		err = kv.Set(keys, dataString)
		if err != nil {
			logrus.WithError(err).Error(keys, " ", "Error setting key")
			http.Error(w, "Internal error", http.StatusBadRequest)
			return
		}

	}
}

func staticFileHandler(w http.ResponseWriter, r *http.Request) {
	keys := r.URL.Path[len("/"):]

	// Create new note if reaches the root path
	if len(keys) == 0 {
		// Generate a random value
		randomValue := generateRandomString(32)

		// Construct the target URL with the random value
		targetURL := fmt.Sprintf("//%s/%s", r.Host, randomValue)

		// Redirect the client to the constructed URL
		http.Redirect(w, r, targetURL, http.StatusFound)
	}

	var keys_slice []string

	// Static pattern identified
	if strings.HasPrefix(keys, "assets/") {
		// Assets
		if strings.HasPrefix(keys, "assets/") {
			keys_slice = append(keys_slice, "static")
			keys_slice = append(keys_slice, strings.Split(keys, "/")...)
		}

		filePath := filepath.Join(keys_slice...)
		fmt.Println(filePath)

		// Check if the file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		// Serve the static file
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, filePath)
		return
	}

	// Loading notes UI
	if len(strings.Split(keys, "/")) > 1 {
		logrus.Error(keys, " ", "Not found")
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	keys_slice = append(keys_slice, "static")
	keys_slice = append(keys_slice, "note.html")
	filePath := filepath.Join(keys_slice...)
	logrus.Debug("Static file ", filePath)
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, filePath)

}

func main() {
	// Setup log system
	// logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)

	/// HTTP server
	http.HandleFunc("/api/", pageApiNote)
	http.HandleFunc("/", staticFileHandler)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	fmt.Println("Starting server at port " + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
