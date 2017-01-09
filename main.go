package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/urfave/cli.v2"
)

var index searchIndex
var conf config

func main() {
	index = Open()
	conf.load()
	app := &cli.App{
		EnableShellCompletion: true,
		Action: func(c *cli.Context) error {
			indexFunc()

			listCurrentFilter()
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "show",
				Usage: "show detailed information of an issue",
				Action: func(c *cli.Context) error {
					indexFunc()
					showDetails(c)
					return nil
				},
				ShellComplete: func(c *cli.Context) {
					indexFunc()
					if c.NArg() > 1 {
						return
					}
					if c.NArg() == 1 {
						res, _ := index.SearchAllMatchingSubString(c.Args().First())
						for _, r := range res {
							fmt.Println(r.key)
						}
					} else {
						res, _ := index.SearchAllMatchingSubString("")
						for _, r := range res {
							fmt.Println(r.key)
						}
					}

					fmt.Println("autocomplete")

				},
			},
			{
				Name:  "clearIndex",
				Usage: "clear the current index",
				Action: func(c *cli.Context) error {
					conf.LastUpdate = time.Time{}
					conf.store()
					index.Clear()
					return nil
				},
			},
		},
	}
	app.Run(os.Args)

}

func showDetails(c *cli.Context) {
	jiraClient, err := jira.NewClient(nil, conf.JiraServer)
	if err != nil {
		panic(err)
	}

	_, err = jiraClient.Authentication.AcquireSessionCookie(conf.Username, conf.Password)
	if err != nil {
		fmt.Printf("%s\n", err)
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}

		bytes, _ := json.MarshalIndent(conf, "", "    ")
		fmt.Printf("Invalid config in %s:\n%s\n", filepath.Join(usr.HomeDir, ".jira.conf"), string(bytes))
		os.Exit(0)
	}

	res, _ := index.getKey(c.Args().First())
	issue := jira.Issue{}
	err = json.Unmarshal([]byte(res.value), &issue)
	if err != nil {
		return
	}
	printIssueDet(issue)
	fmt.Println("\nSimilar issues:")
	resSearch, err := index.getKey(issue.Key)
	if err != nil {
		log.Panic(err)
	}
	index.calculateSimularities(resSearch.key, resSearch.value)
	sim, _ := index.getSimularities(issue.Key)
	if len(sim) == 0 {
		fmt.Println("No similar issues found, please run jira -index to generate them for this issue")
		return
	}
	for _, value := range sim {
		res, _ := index.getKey(value.Key)
		issue := jira.Issue{}
		json.Unmarshal([]byte(res.value), &issue)
		printIssue(issue)
	}

	fmt.Print("\n")
}

func indexFunc() {
	jiraClient, err := jira.NewClient(nil, conf.JiraServer)
	if err != nil {
		panic(err)
	}

	_, err = jiraClient.Authentication.AcquireSessionCookie(conf.Username, conf.Password)
	if err != nil {
		fmt.Printf("%s\n", err)
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}

		bytes, _ := json.MarshalIndent(conf, "", "    ")
		fmt.Printf("Invalid config in %s:\n%s\n", filepath.Join(usr.HomeDir, ".jira.conf"), string(bytes))
		os.Exit(0)
	}

	now := time.Now()
	for i := 0; ; i += 100 {
		searchString := "project in (" + conf.Project + ") AND updated > '" + conf.LastUpdate.Format("2006/01/02 15:04"+"'")
		list, response, err := jiraClient.Issue.Search(searchString, &jira.SearchOptions{StartAt: i, MaxResults: 100})
		if err != nil {
			i -= 100
			b, _ := ioutil.ReadAll(response.Body)
			fmt.Printf("Rolling back 100 commits to get around the error %s %s\n", err, b)
			continue
		}
		if i == 0 && len(list) == 0 {
			conf.LastUpdate = now
			conf.store()
			break
		}
		if len(list) == 0 && i > 0 {
			conf.LastUpdate = now
			conf.store()
			fmt.Println(" new/changed issues")
			index.calculateTfIdf()
			break
		}
		for j, l := range list {
			err = index.Index(l.Key, l)
			if err != nil {
				log.Panic(err)
			}
			fmt.Printf("\r%d", i+j+1)
		}
	}
	//fmt.Print("basic" + basicAuth(0))
	// Get pages from confluence

}

func listCurrentFilter() {
	jiraClient, err := jira.NewClient(nil, conf.JiraServer)
	if err != nil {
		panic(err)
	}

	_, err = jiraClient.Authentication.AcquireSessionCookie(conf.Username, conf.Password)
	if err != nil {
		fmt.Printf("%s\n", err)
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}

		bytes, _ := json.MarshalIndent(conf, "", "    ")
		fmt.Printf("Invalid config in %s:\n%s\n", filepath.Join(usr.HomeDir, ".jira.conf"), string(bytes))
		os.Exit(0)
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
	var priorityID = issue.Fields.Priority.ID
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
	if priorityID == "3" {
		priority = fmt.Sprintf("\033[0;32m%-10s\033[m", priorityValue)
	} else if priorityID == "2" {
		priority = fmt.Sprintf("\033[0;31m%-10s\033[m", priorityValue)
	} else if priorityID == "1" {
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
