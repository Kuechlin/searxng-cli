package main

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type SearchResponse struct {
	Query   string         `json:"query"`
	Results []SearchResult `json:"results"`
}
type SearchResult struct {
	URL       string   `json:"url"`
	ParsedURL []string `json:"parsed_url"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Engine    string   `json:"engine"`
	Engines   []string `json:"engines"`
	Score     float32  `json:"score"`
	Category  string   `json:"category"`
}
type SearchAnswers struct {
	URL       string   `json:"url"`
	ParsedURL []string `json:"parsed_url"`
	Engine    string   `json:"engine"`
	Answer    string   `json:"answer"`
}
type SearchInfobox struct {
	ID      string `json:"id"`
	Engine  string `json:"engine"`
	Infobox string `json:"infobox"`
	Content string `json:"content"`
}

func Search(promt string) (SearchResponse, error) {
	var data SearchResponse

	resp, err := http.PostForm("https://search.kuechlin.dev/search", url.Values{"q": {promt}, "format": {"json"}})
	if err != nil {
		return data, err
	}

	err = json.NewDecoder(resp.Body).Decode(&data)

	return data, err
}
