package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var redisClient *redis.Client
var redis_ctx = context.Background()
var usersCollection *mongo.Collection

func getUserDataFromMongo(userId int) {
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{
					{"id", bson.D{{"$eq", userId}}},
				},
			},
		},
	}
	var result bson.M
	err := usersCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("displaying the first result from the search filter")
	fmt.Println(result)
}

func getUserDataFromCache(userId int) (string, error) {
	userData, err := redisClient.Get(redis_ctx, strconv.Itoa(userId)).Result()
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
			getUserDataFromMongo(userID)
			// add user data to redis
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
	var mongoDbName = "Twitch"
	var mongoCollectionName = "Users"
	var mongoUsername = os.Getenv("MONGO_USERNAME")
	var mongoPassword = os.Getenv("MONGO_PASSWORD")
	var redisPassword = os.Getenv("REDIS_PASSWORD")

	// Creating Redis Client
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis-server:6379",
		Password: redisPassword,
		DB:       0, // use default DB
	})

	// Ping Redis to make sure its working
	ping, err := redisClient.Ping(redis_ctx).Result()
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1) // exit
	}

	if ping != "PONG" {
		log.Fatalf("Expected PONG instead got '%s'", ping)
		os.Exit(1) // exit
	}

	// Creating MongoDB Client
	// Creating mongodb connection string
	URI := fmt.Sprintf("mongodb://%s:%s@%s:%d", mongoUsername, mongoPassword, "mongo-server", 27017)
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1) // exit
	}

	err = mongoClient.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1) // exit
	}

	usersCollection = mongoClient.Database(mongoDbName).Collection(mongoCollectionName)

	// Setting up HTTP server routes
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/user", user)  // /user/:id
	http.HandleFunc("/user/", user) // /user/:id

	fmt.Println("Starting Server")
	http.ListenAndServe(":8001", nil)
}
