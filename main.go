package main

import (
	"github.com/andygrunwald/go-jira"
	"fmt"
	"github.com/blevesearch/bleve"
	"log"
	"github.com/blevesearch/bleve/mapping"
	"github.com/blevesearch/bleve/analysis/lang/en"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/document"
)

var username = ""
var password = ""
var me = username

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
	m, err := buildIndexMapping()
	if err != nil {
		panic(err)
	}
	index, err := bleve.Open("example.bleve")
	if err != nil {
		index, err = bleve.New("example.bleve", m)
	}


	for ; ;  {
		list, _, _ := jiraClient.Issue.Search("", &jira.SearchOptions{StartAt:0, MaxResults:100})

		for _, l := range list {
			err = index.Index(l.ID, l)
			if err != nil {
				log.Panic(err)
			}
		}
	}


	query := bleve.NewMatchAllQuery()
	search := bleve.NewSearchRequest(query)
	search.Size = 100
	search.SortBy([]string{"-_id"})
	searchResults, err := index.Search(search)
	if err != nil {
		log.Panic(err)
	}
	for _, h := range searchResults.Hits {
		b, _ := index.Document(h.ID)
		printIssue(b)
	}
}

func buildIndexMapping() (mapping.IndexMapping, error) {

	// a generic reusable mapping for english text
	textFieldMapping := bleve.NewTextFieldMapping()
	textFieldMapping.Analyzer = en.AnalyzerName

	// a generic reusable mapping for keyword text
	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = keyword.Name

	// a generic reusable mapping for keyword text
	dateFieldMapping := bleve.NewDateTimeFieldMapping()


	m := bleve.NewDocumentMapping()

	m.AddFieldMappingsAt("fields.summary", textFieldMapping)
	m.AddFieldMappingsAt("fields.updated", dateFieldMapping)

	m.AddFieldMappingsAt("fields.description", textFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("issue", m)

	indexMapping.TypeField = "type"
	indexMapping.DefaultAnalyzer = "en"

	return indexMapping, nil
}

func printIssue(issue *document.Document) {
	var priorityValue = ""
	var creator = ""
	var assignee = ""
	var key = ""
	var updated = ""
	var summary = ""
	var status = "status"
	var fix = ""
	for _, value := range issue.Fields {
		if value.Name() == "fields.priority.name" {
			priorityValue = string(value.Value())
		}
		if value.Name() == "fields.Creator.name" {
			creator = string(value.Value())
		}
		if value.Name() == "fields.assignee.name" {
			assignee = string(value.Value())
		}
		if value.Name() == "key" {
			key = string(value.Value())
		}
		if value.Name() == "fields.updated" {
			updated = string(value.Value())
		}
		if value.Name() == "fields.summary" {
			value.Analyze()
			summary = string(value.Value())
		}
		if value.Name() == "fields.status.name" {
			status = string(value.Value())
		}
		if value.Name() == "fields.fixVersions.name" {
			fix = string(value.Value())
		}
	}
	var priority = ""
	if priorityValue == "Minor" {
		priority = fmt.Sprintf("\033[0;32m%-10s\033[m", priorityValue)
	} else if priorityValue == "Major" {
		priority = fmt.Sprintf("\033[0;31m%-10s\033[m", priorityValue)
	} else if priorityValue == "Blocker" {
		priority = fmt.Sprintf("\033[0;30;41m%-10s\033[m", priorityValue)
	} else {
		priority = fmt.Sprintf("%s", priorityValue)
	}
	if assignee == me {
		assignee = fmt.Sprintf("\033[1;31m%-10s\033[m", me)
	}
	if creator == me {
		assignee = fmt.Sprintf("\033[1;31m-10%s\033[m", me)
	}
	fmt.Printf("%-15s %-15s %-10s %-10s %-10s %-20s %-20s %s\n",
		key, updated[:len("2006-01-02T15:04:05")], priority, assignee, creator, fix, status, summary)
}
