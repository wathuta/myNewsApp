package model

import "time"

//Article is a model struct for a news article from the news api
type Article struct {
	Source      Source    `json:"source"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	URLToImage  string    `json:"urlToImage"`
	PublishedAt time.Time `json:"publishedAt"`
	Content     string    `json:"content"`
}

//Result is a model struct for the results obtained after the call to the News api
type Result struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
}

//Source is a struct for source of the article
type Source struct {
	ID   interface{} `json:"id"`
	Name string      `json:"name"`
}

//Search is a struct that represents the Query made by the user
type Search struct {
	Query      string
	NextPage   int
	TotalPages int
	Results    *Result
}
