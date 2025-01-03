package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/craniacshencil/cachecast/internal"
	"github.com/craniacshencil/cachecast/utils"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

// Returns False if file does not exist
func checkFileExists(name string) bool {
	_, err := os.Stat(name)
	return !errors.Is(err, os.ErrNotExist)
}

func init() {
	// This works only if bash is installed in docker image
	// APP_ENV := os.Getenv("APP_ENV")
	log.Println(checkFileExists(".env"))
	if checkFileExists(".env") {
		if err := godotenv.Load(".env"); err != nil {
			log.Println("ERR: While opening .env file:", err)
			return
		}
	}
	if err := godotenv.Load(".env.docker"); err != nil {
		log.Println("ERR: While opening .env.docker file:", err)
		return
	}
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       1,
		Protocol: 2,
	})

	cacheClient := internal.NewCacheClient(client)

	router := http.NewServeMux()
	router.Handle("/", rateLimiter(servePage))
	router.Handle("POST /", rateLimiter(cacheClient.GetWeather))

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Println("ERR: While serving ", err)
		return
	}
}

func rateLimiter(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	limiter := rate.NewLimiter(2, 4)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			message := "Too many requests"
			utils.WriteJSON(w, http.StatusTooManyRequests, message)
			return
		} else {
			next(w, r)
		}
	})
}

func servePage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./web/index.html")
	if err != nil {
		log.Println("ERR: While creating template for HTML file:", err)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		log.Println("ERR: While executing template:", err)
		return
	}
}
