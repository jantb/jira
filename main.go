package main

import (
	"github.com/andygrunwald/go-jira"
	"fmt"
	"log"
	"encoding/json"
)

var username = ""
var password = ""
var me = username

var jiraServer = ""
var index datastore
func main() {

	jiraClient, err := jira.NewClient(nil, jiraServer)
	if err != nil {
		panic(err)
	}

	res, err := jiraClient.Authentication.AcquireSessionCookie(username, password)
	if err != nil || res == false {
		fmt.Printf("Result: %v\n", res)
		panic(err)
	}

	index = Open()

	for i := 0; ;i+=100  {
		list, _, _ := jiraClient.Issue.Search("", &jira.SearchOptions{StartAt:i, MaxResults:i+100})
		if len(list) == 0 {
			resSearch, err := index.SearchAllMatching(1000000)
			if err != nil {
				log.Panic(err)
			}
			for i, value := range resSearch {
				fmt.Println(i)
				var issue jira.Issue
				json.Unmarshal(value, &issue)
				index.calculateSimularities(issue.Key,string(value))
			}
			break
		}
		for _, l := range list {
			err = index.Index(l.Key, l)
			if err != nil {
				log.Panic(err)
			}
		}
		fmt.Println(i)
	}

	resSearch, err := index.SearchAllMatching(100)
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
	b, _:=json.Marshal(issue)
	index.calculateSimularities(issue.Key,string(b))
	sim,_ := index.getSimularities(issue.Key)
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
		if len(fix)!=0 {
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
	if assignee == me {
		assignee = fmt.Sprintf("\033[1;31m%-10s\033[m", me)
	}
	if creator == me {
		assignee = fmt.Sprintf("\033[1;31m-10%s\033[m", me)
	}
	fmt.Printf("%-15s %-15s %-10s %-10s %-10s %-20s %-20s %s\n",
		key, updated[:len("2006-01-02T15:04:05")], priority, assignee, creator, fix, status, summary)
}
