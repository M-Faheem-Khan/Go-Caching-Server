package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
)

var REDIS_PASSWORD = os.Getenv("REDIS_PASSWORD")
var client *redis.Client
var redis_ctx = context.Background()

func getUserDataFromCache(userId int) (string, error) {
	userData, err := client.Get(redis_ctx, strconv.Itoa(userId)).Result()
	if err != nil {
		return "", err
	}

	return userData, nil
}

// getUserID - return user id from the url path
func getUserID(path string) (int, error) {
	uri := strings.Split(path, "/")
	if len(uri) <= 2 || uri[2] == "" {
		return 0, errors.New("no userID given")
	}

	// Verify the userID is valid
	// userId must be an integer
	// userId must be greater than 0
	userID, err := strconv.Atoi(uri[2])
	if err != nil || userID <= 0 {
		return 0, errors.New("invalid userID")
	}

	return userID, nil
}

func user(w http.ResponseWriter, r *http.Request) {
	// Getting userID
	userID, err := getUserID(r.URL.Path)
	if err != nil {
		fmt.Fprintf(w, "UserID: %s\n", err.Error())
	} else {
		// Getting user data from cache
		userData, err := getUserDataFromCache(userID)
		if err != nil {
			// Fetch user data from mongodb
			// add to redis
			// send user data to requester
			fmt.Fprintf(w, "Error: %s\n", err.Error())
		} else {
			fmt.Fprintf(w, "UserID: %s\n", userData)
		}
	}

}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Twitch from docker!\n")
}

func main() {
	// Creating Redis Client
	client = redis.NewClient(&redis.Options{
		Addr:     "redis-server:6379",
		Password: REDIS_PASSWORD,
		DB:       0, // use default DB
	})

	// Ping Redis to make sure its working
	ping, err := client.Ping(redis_ctx).Result()
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1) // exit
	}

	if ping != "PONG" {
		log.Fatalf("Expected PONG instead got '%s'", ping)
		os.Exit(1) // exit
	}

	// Setting up HTTP server routes
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/user", user)  // /user/:id
	http.HandleFunc("/user/", user) // /user/:id

	fmt.Println("Starting Server")
	http.ListenAndServe(":8001", nil)
}
