package news

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/wathuta/newsappMicroservice/model"
)

//Client is used to make an http calls to the news api
type Client struct {
	http     *http.Client //a http client that will be used to make requests to the news api
	Key      string       //the api key used to access the news api
	PageSize int          //holds the number of results to return per page
}

//NewClient is the entry point.It returns a new client
func NewClient(http *http.Client, Key string, PageSize int) *Client {
	if PageSize > 100 {
		PageSize = 100
	}
	return &Client{
		http:     http,
		Key:      Key,
		PageSize: PageSize,
	}
}

//FetchEverything makes a request to the everithing endpoint of the newsapi using the client parameters
func (c *Client) FetchEverything(query, page string) (*model.Result, error) {
	endpoint := fmt.Sprintf("https://newsapi.org/v2/everything?q=%s&pageSize=%d&page=%s&apiKey=%s&sortBy=publishedAt&language=en", url.QueryEscape(query), c.PageSize, page, c.Key)
	resp, err := http.Get(endpoint)
	if err != nil {
		fmt.Println("unable to send request to api")
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(string(body))
	}
	res := &model.Result{}
	return res, json.Unmarshal(body, &res)

}
