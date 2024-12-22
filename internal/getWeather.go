package internal

import (
	"log"
	"net/http"
)

func GetWeather(w http.ResponseWriter, r *http.Request) {
	log.Println("Hitting this endpoint")
}
