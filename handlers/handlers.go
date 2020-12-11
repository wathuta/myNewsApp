package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"

	"github.com/wathuta/newsappMicroservice/model"

	"github.com/wathuta/newsappMicroservice/news"
)

var tmpl = template.Must(template.ParseFiles("./index.html"))

//News an struct that is used for dependancy injjection
type News struct {
	l *log.Logger
}

//NewNews the entry point to the News handler
func NewNews(l *log.Logger) *News {
	return &News{l}
}

//IndexHandler implements the handler interface and serves the template
func (n *News) IndexHandler(w http.ResponseWriter, r *http.Request) {
	//the template is first written to buf to catch any error that might occure during execution
	buf := &bytes.Buffer{}
	err := tmpl.Execute(buf, nil)
	if err != nil {
		http.Error(w, "[error] unable to write templates to bufffer", http.StatusInternalServerError)
	}
	_, err = buf.WriteTo(w)
	if err != nil {
		n.l.Fatal("error unable to render template from buffer")
		http.Error(w, "[error] unable to write templates from bufffer to w", http.StatusInternalServerError)
	}
}

//SearchHandler implements the handler interface and is responsible for sending queries to the newsApi endpoint
func (n *News) SearchHandler(newsAPI *news.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.URL.String())
		if err != nil {
			http.Error(w, "Error invalid url", http.StatusBadRequest)
			n.l.Fatal("unable to parse url")
			return
		}
		params := u.Query()
		searchQuery := params.Get("q")
		page := params.Get("page")
		if page == "" {
			page = "1"
		}
		results, err := newsAPI.FetchEverything(searchQuery, page)
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			n.l.Fatal(err)
			return
		}
		nextPage, err := strconv.Atoi(page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		search := &model.Search{
			Query:      searchQuery,
			NextPage:   nextPage,
			TotalPages: int(math.Ceil(float64(results.TotalResults / newsAPI.PageSize))),
			Results:    results,
		}
		buf := &bytes.Buffer{}
		err = tmpl.Execute(buf, search)
		if err != nil {
			n.l.Println("unable to parse template  into the buffer")
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}
		_, err = buf.WriteTo(w)
		if err != nil {
			fmt.Println("unable to parse buffer to the response writer ")
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}

	}

}
