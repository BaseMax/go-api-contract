package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Post struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

var posts = []Post{
	{UserID: 1, ID: 1, Title: "First Post", Body: "This is the first post"},
	{UserID: 1, ID: 2, Title: "Second Post", Body: "This is the second post"},
}

func main() {
	http.HandleFunc("/posts/", handlePost)
	http.HandleFunc("/posts", handlePosts)

	fmt.Println("Test server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from path
	path := strings.TrimPrefix(r.URL.Path, "/posts/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Find post
	for _, post := range posts {
		if post.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(post)
			return
		}
	}

	http.Error(w, "Post not found", http.StatusNotFound)
}

func handlePosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(posts)
	case "POST":
		var newPost Post
		if err := json.NewDecoder(r.Body).Decode(&newPost); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		newPost.ID = len(posts) + 1
		posts = append(posts, newPost)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newPost)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
