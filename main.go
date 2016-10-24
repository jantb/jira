package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

var index searchIndex
var conf config
var indexIssues = flag.Bool("index", false, "Index jira, and generate similarities, uses a timestamp to only update new issues")
var clearIndex = flag.Bool("clearIndex", false, "Clear the index and reset the timestamp")

func main() {
	flag.Parse()
	conf.load()
	os.Args = flag.Args()
	index = Open()

	if *clearIndex {
		conf.LastUpdate = time.Time{}
		conf.store()
		index.Clear()
		os.Exit(0)
	}

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

	if *indexIssues {
		now := time.Now()
		for i := 0; ; i += 100 {
			searchString := "project in (" + conf.Project + ") AND updated > '" + conf.LastUpdate.Format("2006/01/02 15:04"+"'")
			list, _, err := jiraClient.Issue.Search(searchString, &jira.SearchOptions{StartAt: i, MaxResults: 100})
			if err != nil {
				i -= 100
				continue
			}
			if i == 0 && len(list) == 0 {
				conf.LastUpdate = now
				conf.store()
				os.Exit(0)
			}
			if len(list) == 0 && i > 0 {
				resSearch, err := index.SearchAllMatching(1000000)
				if err != nil {
					log.Panic(err)
				}
				for _, value := range resSearch {
					index.calculateSimularities(value.key, value.value)
				}

				conf.LastUpdate = now
				conf.store()
				fmt.Println()
				fmt.Println("Done indexing")

				os.Exit(0)
			}
			for _, l := range list {
				comments := ""
				if l.Fields.Comments != nil {
					for _, comment := range l.Fields.Comments.Comments {
						if len(comments) != 0 {
							comments += " "
						}
						comments += comment.Body
					}
				}
				err = index.Index(l.Key, fmt.Sprintf("%s %s %s", l.Fields.Summary, l.Fields.Description, comments))
				if err != nil {
					log.Panic(err)
				}
			}
			fmt.Printf("\r%d", i)
		}
	}

	if len(os.Args) == 1 {
		list, _, _ := jiraClient.Issue.Search("key = "+os.Args[0], &jira.SearchOptions{StartAt: 0, MaxResults: 100})
		for _, value := range list {
			printIssueDet(value)
			fmt.Println("\nSimilar issues:")
			sim, _ := index.getSimularities(value.Key)
			keys := ""
			for _, value := range sim[:25] {
				if len(keys) != 0 {
					keys += ","
				}
				keys += value.Key
			}
			list, _, _ := jiraClient.Issue.Search("key in ("+keys+")", &jira.SearchOptions{StartAt: 0, MaxResults: 100})
			for _, value := range list {
				printIssue(value)
			}
			fmt.Print("\n")
		}

		if len(list) == 0 {
			list, _, _ := jiraClient.Issue.Search("text ~ \""+os.Args[0]+"\" ", &jira.SearchOptions{StartAt: 0, MaxResults: 100})
			for _, value := range list {
				printIssue(value)
			}
		}
		return
	}

	searchString := "filter=" + conf.Filter
	list, _, _ := jiraClient.Issue.Search(searchString, &jira.SearchOptions{StartAt: 0, MaxResults: 100})
	fmt.Println("Next fix version " + getNextFixVersion())
	for _, issue := range list {
		printIssue(issue)
	}

}

func printIssueDet(issue jira.Issue) {
	var fix = ""
	for _, fixversion := range issue.Fields.FixVersions {
		if len(fix) == 0 {
			fix += fixversion.Name
		} else {
			fix += ", " + fixversion.Name
		}
	}
	fmt.Printf("\033[32m%-10s\033[m ", issue.Fields.Created)
	fmt.Printf("\033[31m%-10s\033[m ", issue.Fields.Type.Name)
	fmt.Printf("\033[33m%-10s\033[m ", issue.Fields.Status.Name)
	if issue.Fields.Creator != nil {
		fmt.Printf("\033[34m%-10s\033[m ", issue.Fields.Creator.Name)
	}
	if issue.Fields.Assignee != nil {
		fmt.Printf("\033[35m%-10s\033[m ", issue.Fields.Assignee.Name)
	}
	fmt.Printf("\033[36m%-10s\033[m ", fix)
	fmt.Printf("\n%s\n", issue.Fields.Summary)
	fmt.Print("\n")
	fmt.Printf("%s\n", issue.Fields.Description)

	if issue.Fields.Comments != nil {
		fmt.Printf("%s\n", "Comments:")
		for _, comment := range issue.Fields.Comments.Comments {
			fmt.Printf("%s \033[0;36m%s\033[m %s\n", comment.Created, comment.Author.Name, comment.Body)
		}
	}

	fmt.Printf("\033[34m%-10s\033[m\n", conf.JiraServer+"browse/"+issue.Key)
}

func printIssue(issue jira.Issue) {
	var priorityValue = issue.Fields.Priority.Name
	var creator = ""
	if issue.Fields.Creator != nil {
		creator = issue.Fields.Creator.Name
	}
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

func getNextFixVersion() string {
	type Result struct {
		Value string `xml:"version"`
	}
	bytes, err := ioutil.ReadFile("pom.xml")
	if err != nil {
		return ""
	}
	var pom Result
	err = xml.Unmarshal(bytes, &pom)
	if err != nil {
		log.Panic(err)
	}
	return pom.Value[:strings.Index(pom.Value, "-")]
}
