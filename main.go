package main

import (
	"github.com/andygrunwald/go-jira"
	"fmt"
	"github.com/blevesearch/bleve"
	"log"
	"github.com/blevesearch/bleve/mapping"
	"github.com/blevesearch/bleve/analysis/lang/en"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
)

var username = ""
var password = ""

// Jira server without http://
var jiraServer = ""

func main() {

	jiraClient, err := jira.NewClient(nil, "http://" + jiraServer)
	if err != nil {
		panic(err)
	}

	res, err := jiraClient.Authentication.AcquireSessionCookie(username, password)
	if err != nil || res == false {
		fmt.Printf("Result: %v\n", res)
		panic(err)
	}

	list, _, _ := jiraClient.Issue.Search("", &jira.SearchOptions{StartAt:0, MaxResults:1})

	m, err := buildIndexMapping()
	if err != nil {
		panic(err)
	}
	index, err := bleve.Open("example.bleve")
	if err != nil {
		index, err = bleve.New("example.bleve", m)
	}

	for _, l := range list {
		fmt.Println(l.Fields.Description)
		fmt.Println(l.Fields.Summary)
		err = index.Index(l.ID, l)
		if err != nil {
			log.Panic(err)
		}
	}

	query := bleve.NewMatchQuery("begge")
	search := bleve.NewSearchRequest(query)
	searchResults, err := index.Search(search)
	if err != nil {
		log.Panic(err)
	}
	b, err := index.Document(searchResults.Hits[0].ID)
	fmt.Println(string(b.Fields[2].Value()))
	//index.
}

func buildIndexMapping() (mapping.IndexMapping, error) {

	// a generic reusable mapping for english text
	textFieldMapping := bleve.NewTextFieldMapping()
	textFieldMapping.Analyzer = en.AnalyzerName

	// a generic reusable mapping for keyword text
	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = keyword.Name

	mapping := bleve.NewDocumentMapping()

	mapping.AddFieldMappingsAt("Summary", textFieldMapping)

	mapping.AddFieldMappingsAt("Description", textFieldMapping)


	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("issue", mapping)

	indexMapping.TypeField = "type"
	indexMapping.DefaultAnalyzer = "en"

	return indexMapping, nil
}
