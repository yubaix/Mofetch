package omdb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// JSON Structuring for future requests
type Film struct {
	Title      string `json:"Title"`
	Year       string `json:"Year"`
	Director   string `json:"Director"`
	Rating     string `json:"Rated"`
	Actors     string `json:"Actors"`
	Plot       string `json:"Plot"`
	Boxoffice  string `json:"BoxOffice"`
	Language   string `json:"Language"`
	Runtime    string `json:"Runtime"`
	ImdbRating string `json:"imdbRating"`
	Poster     string `json:"Poster"`
}

type Api struct {
	apiKey string
}

// Again, thanks to Ashish Kumar for the idea of this function
func NewOMDBClient(apiKey string) *Api {
	return &Api{apiKey: apiKey}
}

func (a *Api) Search(query string) (*Film, error) {
	// Fetching the film data from the OMDB API
	// Thanks to Ashish Kumar's Mufetch for the idea for params
	params := url.Values{}
	params.Add("t", query)
	params.Add("plot", "short")

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://www.omdbapi.com/?apikey="+a.apiKey+"&"+params.Encode(), nil) // Incredibly messy, PRs welcome
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
	}

	defer resp.Body.Close()

	var film Film
	if err := json.NewDecoder(resp.Body).Decode(&film); err != nil {
		return nil, err
	}

	return &film, nil
}
