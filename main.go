package main

import (
	"github.com/andygrunwald/go-jira"
	"fmt"
	"github.com/blevesearch/bleve"
	"log"
)

func main() {
	jiraClient, _ := jira.NewClient(nil, "https://jira.atlassian.com")
	list, r, _ := jiraClient.Issue.Search("", &jira.SearchOptions{StartAt:0, MaxResults:100})

	mapping := bleve.NewIndexMapping()
	index, err := bleve.Open("example.bleve")
	if err != nil {
		index, err = bleve.New("example.bleve", mapping)
	}

	for _, l := range list {
		fmt.Println(l.Fields.Description)
		err = index.Index(l.ID, l)
		if err != nil {
			log.Panic(err)
		}
	}
	fmt.Println(list)
	fmt.Println(*r)

	query := bleve.NewMatchQuery("Atlassian in Korea")
	search := bleve.NewSearchRequest(query)
	searchResults, err := index.Search(search)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(searchResults)
	index.
}
