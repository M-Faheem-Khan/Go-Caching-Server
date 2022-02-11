package main

import (
	"fmt"
	"net/http"
	"strings"
)

// getUserID - Extracts user id from the url path
func getUserID(path string) string {
	uri := strings.Split(path, "/")
	if len(uri) <= 2 || uri[2] == "" {
		return "no userID given!"
	}
	return uri[2]
}

func user(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "UserID: %s\n", getUserID(r.URL.Path))
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Twitch from docker!")
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/user/", user)

	fmt.Println("Starting Server")
	http.ListenAndServe(":8001", nil)
}
