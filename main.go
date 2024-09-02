package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gocolly/colly/v2"
	"github.com/gorilla/mux"
)

type Article struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type Standing struct {
	Name     string `json:"name"`
	Position int    `json:"pos"`
}

var articles []Article
var standings []Standing

func getStandings(w http.ResponseWriter, r *http.Request) {
	standings = []Standing{}
	scrapeURL := "https://www.formula1.com/en/results/2024/drivers"
	w.Header().Set("Content-Type", "application/json")
	c := colly.NewCollector(colly.AllowedDomains("www.formula1.com", "formula1.com"))
	index := 0
	position := 0
	c.OnHTML("a.underline.underline-offset-normal.decoration-1.decoration-greyLight.hover\\:decoration-brand-primary", func(h *colly.HTMLElement) {
		if index%2 == 0 {
			standings = append(standings, Standing{
				Name:     h.Text[:len(h.Text)-3],
				Position: position + 1,
			})
			position++
		}
		index++
	})
	c.Visit(scrapeURL)
	json.NewEncoder(w).Encode(standings)
}

func getArticles(w http.ResponseWriter, r *http.Request) {
	scrapeURL := "https://www.formula1.com/"
	c := colly.NewCollector(colly.AllowedDomains("www.formula1.com", "formula1.com"))
	index := 0
	c.OnHTML("p.f1--s.no-margin", func(h *colly.HTMLElement) {
		if h.Text != "Sorry" {
			articles = append(articles, Article{
				ID:    index + 1,
				Title: h.Text,
			})
			index++
		}
	})
	c.Visit(scrapeURL)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(articles)
}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/articles", getArticles).Methods("GET")
	router.HandleFunc("/driver-standings", getStandings).Methods("GET")
	fmt.Printf(("Starting server at 8000 \n"))
	log.Fatal((http.ListenAndServe(":8000", router)))
}
