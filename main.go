package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/wathuta/newsappMicroservice/model"

	"github.com/wathuta/newsappMicroservice/news"

	"github.com/joho/godotenv"
	"github.com/wathuta/newsappMicroservice/handlers"
)

func main() {
	//loadin the .env file which contains information that is sensitive
	err := godotenv.Load("./file.env")
	if err != nil {
		log.Fatal("unable to load the env file")
	}

	l := log.New(os.Stdout, "Newserver-Api", log.LstdFlags)

	mux := http.NewServeMux()

	port := os.Getenv("PORT")
	apiKey := os.Getenv("APIKEY")
	//setting up a client to use with the news API
	myClient := &http.Client{Timeout: 10 * time.Second}
	newsAPI := news.NewClient(myClient, apiKey, 20)
	//creating handlers
	nh := handlers.NewNews(l)
	fs := http.FileServer(http.Dir("./assets"))
	//handling the routing
	mux.HandleFunc("/", nh.IndexHandler)
	mux.HandleFunc("/search", nh.SearchHandler(newsAPI))

	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.ListenAndServe("localhost:"+port, mux)
}

//Search is a struct that represents the Query made by the user
type Search struct {
	Query      string
	NextPage   int
	TotalPages int
	Results    *model.Result
}
