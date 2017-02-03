package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/jantb/go-jira"
	"gopkg.in/urfave/cli.v2"
	"strconv"
	"unicode/utf8"
	"regexp"
)

var index SearchIndex
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
				Name:  "search",
				Usage: "search issues and Confluence pages for string",
				Action: func(c *cli.Context) error {
					indexFunc()

					var issues = []jira.Issue{}
					var confluence = []string{}
					sim, _ := index.IndexSearch(strings.ToLower(strings.Join(c.Args().Slice(), " ")))
					for _, value := range sim {
						res, _ := index.getKey(value.Key)

						issue := jira.Issue{}
						err := json.Unmarshal([]byte(res.value), &issue)
						if err != nil {
							res, _ = index.getConfluenceKey(value.Key)
							page := Page{}
							json.Unmarshal([]byte(res.value), &page)
							confluence = append(confluence, fmt.Sprintf("%s \033[34m%-10s\033[m", page.Title, page.Link))
						} else {
							issues = append(issues, issue)
						}

					}
					if len(issues) > 20 {
						issues = issues[:20]
					}
					if len(confluence) > 20 {
						confluence = confluence[:20]
					}

					for _, value := range formatIssues(issues) {
						fmt.Println(value)
					}

					fmt.Println("\nRelated confluence pages:")
					for _, c := range confluence {
						fmt.Println(c)
					}
					return nil
				},
			},
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

				},
			},
			{
				Name:  "assign",
				Usage: "assign to user",
				Action: func(c *cli.Context) error {
					indexFunc()
					assignToUser(c)
					return nil
				},
				ShellComplete: func(c *cli.Context) {
					indexFunc()
					if c.NArg() > 2 {
						return
					}

					if c.NArg() == 1 {
						res, _ := index.SearchAllMatchingSubString(c.Args().First())
						for _, r := range res {
							fmt.Println(r.key)
						}
					} else if c.NArg() == 0 {
						res, _ := index.SearchAllMatchingSubString("")
						for _, r := range res {
							fmt.Println(r.key)
						}
					}
				},
			},
			{
				Name:  "clearIndex",
				Usage: "clear the current index",
				Action: func(c *cli.Context) error {
					conf.LastUpdate = time.Time{}
					conf.LastUpdateConfluence = time.Time{}
					conf.store()
					index.Clear()
					return nil
				},
			},
		},
	}
	app.Run(os.Args)

}

func formatIssues(issues []jira.Issue) []string {
	print := [][]string{}
	for _, issue := range issues {
		print = append(print, printIssue(issue))
	}
	length := make([]int, len(print[0]))
	for _, i := range print {
		for index, i2 := range i {
			length[index] = Max(utf8.RuneCount([]byte(removeColor(i2))), length[index])
		}
	}
	for _, i := range print {
		for key, leng := range length {
			i[key] = i[key] +fmt.Sprintf("%-" + strconv.Itoa(leng+1 - utf8.RuneCount([]byte(removeColor(i[key])))) + "s", " ")
		}
	}

	ret := []string{}
	for _, value := range print {
		ret = append(ret, strings.Join(value, ""))
	}
	return ret
}


func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func assignToUser(c *cli.Context) {
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
	r, err := jiraClient.Issue.Assign(c.Args().First(), conf.Username)
	if err != nil {
		fmt.Printf("%s\n", err)
		b, _ := ioutil.ReadAll(r.Body)
		fmt.Println(string(b))
	}
}

func showDetails(c *cli.Context) {
	res, _ := index.getKey(c.Args().First())
	issue := jira.Issue{}
	err := json.Unmarshal([]byte(res.value), &issue)
	if err != nil {
		return
	}
	printIssueDet(issue)
	fmt.Println("\nSimilar issues:")
	resSearch, err := index.getKey(issue.Key)
	if err != nil {
		log.Panic(err)
	}
	index.calculateSimularities(resSearch.key)
	sim, _ := index.getSimularities(issue.Key)
	if len(sim) == 0 {
		fmt.Println("No similar issues found, please run jira -index to generate them for this issue")
		return
	}
	var issues = []jira.Issue{}
	var confluence = []string{}
	for _, value := range sim {
		res, _ := index.getKey(value.Key)

		issue := jira.Issue{}
		err := json.Unmarshal([]byte(res.value), &issue)
		if err != nil {
			res, _ = index.getConfluenceKey(value.Key)
			page := Page{}
			json.Unmarshal([]byte(res.value), &page)
			confluence = append(confluence, fmt.Sprintf("%s \033[34m%-10s\033[m", page.Title, page.Link))
		} else {
			issues = append(issues, issue)
		}

	}
	if len(issues) > 20 {
		issues = issues[:20]
	}
	if len(confluence) > 20 {
		confluence = confluence[:20]
	}

	for _, issue := range formatIssues(issues) {
		fmt.Println(issue)
	}
	fmt.Println("\nRelated confluence pages:")
	for _, c := range confluence {
		fmt.Println(c)
	}

	fmt.Print("\n")
}

func indexFunc() {
	if conf.ConfluenceServer != "" {
		for _, page := range getConfluencePages() {
			index.IndexConfluence(page.Key, page)
		}
	}
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

	for i := 0; ; i += 25 {
		searchString := "project in (" + conf.Project + ") AND updated > '" + conf.LastUpdate.Format("2006/01/02 15:04" + "'")
		list, response, err := jiraClient.Issue.Search(searchString, &jira.SearchOptions{StartAt: i, MaxResults: 25})
		if err != nil {
			i -= 25
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
			index.calculateTfIdf()
			break
		}
		for j, l := range list {
			err = index.Index(l.Key, l)
			if err != nil {
				log.Panic(err)
			}
			fmt.Printf("                                               \r%d new/changed issues", i + j + 1)
		}
	}
	fmt.Println()

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
	for _, value := range formatIssues(list) {
		fmt.Println(value)
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
		fmt.Printf("\033[34m%-10s\033[m ", issue.Fields.Creator.DisplayName)
	}
	if issue.Fields.Assignee != nil {
		fmt.Printf("\033[35m%-10s\033[m ", issue.Fields.Assignee.DisplayName)
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

	fmt.Printf("\033[34m%-10s\033[m\n", conf.JiraServer + "browse/" + issue.Key)
}

func printIssue(issue jira.Issue) (ret []string) {

	var priorityValue = issue.Fields.Priority.Name
	var priorityID = issue.Fields.Priority.ID
	var creator = ""
	if issue.Fields.Creator != nil {
		creator = issue.Fields.Creator.DisplayName
	}
	var assignee = ""
	if issue.Fields.Assignee != nil {
		assignee = issue.Fields.Assignee.DisplayName
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
		priority = fmt.Sprintf("\033[0;32m%s\033[m", priorityValue)
	} else if priorityID == "2" {
		priority = fmt.Sprintf("\033[0;31m%s\033[m", priorityValue)
	} else if priorityID == "1" {
		priority = fmt.Sprintf("\033[0;30;41m%s\033[m", priorityValue)
	} else {
		priority = fmt.Sprintf("%s", priorityValue)
	}
	if issue.Fields.Assignee !=nil && issue.Fields.Assignee.Name == conf.Username {
		assignee = fmt.Sprintf("\033[1;31m%-10s\033[m", assignee)
	}
	ret = append(ret, key)
	ret = append(ret, updated[:len("2006-01-02T15:04:05")])
	ret = append(ret, priority)
	ret = append(ret, assignee)
	ret = append(ret, creator)
	ret = append(ret, fix)
	ret = append(ret, status)
	ret = append(ret, summary)
	return ret
}
func removeColor(x string) string{
	re := regexp.MustCompile("\\033\\[[0-9;]*m")
	x = re.ReplaceAllString(x, "")
	return x
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
