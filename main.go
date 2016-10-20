package main

import (
	"github.com/andygrunwald/go-jira"
	"fmt"
	"log"
	"encoding/json"
	"time"
	"flag"
	"os"
)

var index searchIndex
var conf config
var indexIssues = flag.Bool("index", false, "Index jira, and generate similarities, uses a timestamp to only update new issues")

func main() {
	flag.Parse()
	conf.load()
	jiraClient, err := jira.NewClient(nil, conf.JiraServer)
	if err != nil {
		panic(err)
	}

	_, err = jiraClient.Authentication.AcquireSessionCookie(conf.Username, conf.Password)
	if err != nil {
		bytes, _ := json.MarshalIndent(conf, "", "    ")
		fmt.Printf("Invalid config:\n%s\n", string(bytes))
		panic(err)
	}
	index = Open()
	if *indexIssues {
		now := time.Now()
		for i := 0; ; i += 100 {
			searchString := "project=" + conf.Project + " AND updated > '" + conf.LastUpdate.Format("2006/01/02 15:04" + "'")
			list, _, err := jiraClient.Issue.Search(searchString, &jira.SearchOptions{StartAt:i, MaxResults:100})
			if err != nil {
				i -= 100
				continue
			}
			if len(list) == 0 {
				resSearch, err := index.SearchAllMatching(1000000)
				if err != nil {
					log.Panic(err)
				}
				for _, value := range resSearch {
					var issue jira.Issue
					err := json.Unmarshal(value, &issue)
					if err != nil {
						continue
					}
					index.calculateSimularities(issue.Key, string(value))
				}

				conf.LastUpdate = now
				conf.store()
				os.Exit(0)
			}
			for _, l := range list {
				err = index.Index(l.Key, l)
				if err != nil {
					log.Panic(err)
				}
			}
			fmt.Printf("\r%d", i)
		}
	}

	resSearch, err := index.SearchAllMatching(1000)
	if err != nil {
		log.Panic(err)
	}
	for _, value := range resSearch {
		var issue jira.Issue
		json.Unmarshal(value, &issue)
		printIssue(issue)
		printSimularities(issue)
	}

}

func printSimularities(issue jira.Issue) {
	b, _ := json.Marshal(issue)
	index.calculateSimularities(issue.Key, string(b))
	sim, _ := index.getSimularities(issue.Key)
	for _, value := range sim[:10] {
		fmt.Print(value.Key + " ")
	}
	fmt.Print("\n")
}
func printIssue(issue jira.Issue) {
	var priorityValue = issue.Fields.Priority.Name
	var creator = issue.Fields.Creator.Name
	var assignee = ""
	if issue.Fields.Assignee != nil {
		assignee = issue.Fields.Assignee.Name
	}
	var key = issue.Key
	var updated = issue.Fields.Updated
	var summary = issue.Fields.Summary
	var status = issue.Fields.Status.Name
	var fix = ""

	for _, value := range issue.Fields.FixVersions {
		if len(fix) != 0 {
			fix += " "
		}
		fix += value.Name
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
	if assignee == conf.Username {
		assignee = fmt.Sprintf("\033[1;31m%-10s\033[m", conf.Username)
	}
	if creator == conf.Username {
		assignee = fmt.Sprintf("\033[1;31m-10%s\033[m", conf.Username)
	}
	fmt.Printf("%-15s %-15s %-10s %-10s %-10s %-20s %-20s %s\n",
		key, updated[:len("2006-01-02T15:04:05")], priority, assignee, creator, fix, status, summary)
}
