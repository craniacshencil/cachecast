package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/craniacshencil/cachecast/internal"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("ERR: While opening .env file")
		log.Println(err)
		return
	}
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("./web/index.html")
		if err != nil {
			log.Println("ERR: While creating template for HTML file")
			log.Println(err)
			return
		}

		err = t.Execute(w, nil)
		if err != nil {
			log.Println("ERR: While executing template")
			log.Println(err)
			return
		}
	})
	router.HandleFunc("POST /getWeather", internal.GetWeather)

	server := &http.Server{
		Addr:    os.Getenv("SERVER_ADDR"),
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Println("ERR: While serving")
		log.Println(err)
		return
	}
}
