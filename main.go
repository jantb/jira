package main

import (
	"net/http"
	"log"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"os"
	"flag"
	"net/url"
	"encoding/xml"
	"strings"
)

type Jira struct {
	Expand     string `json:"expand"`
	Issues     []struct {
		Expand string `json:"expand"`
		Fields struct {
			       Aggregateprogress             struct {
								     Percent  int `json:"percent"`
								     Progress int `json:"progress"`
								     Total    int `json:"total"`
							     } `json:"aggregateprogress"`
			       Aggregatetimeestimate         int `json:"aggregatetimeestimate"`
			       Aggregatetimeoriginalestimate int `json:"aggregatetimeoriginalestimate"`
			       Aggregatetimespent            int `json:"aggregatetimespent"`
			       Assignee                      struct {
								     Active       bool `json:"active"`
								     AvatarUrls   struct {
											  One6x16   string `json:"16x16"`
											  Two4x24   string `json:"24x24"`
											  Three2x32 string `json:"32x32"`
											  Four8x48  string `json:"48x48"`
										  } `json:"avatarUrls"`
								     DisplayName  string `json:"displayName"`
								     EmailAddress string `json:"emailAddress"`
								     Name         string `json:"name"`
								     Self         string `json:"self"`
							     } `json:"assignee"`
			       Components                    []interface{} `json:"components"`
			       Created                       string        `json:"created"`
			       Creator                       struct {
								     Active       bool `json:"active"`
								     AvatarUrls   struct {
											  One6x16   string `json:"16x16"`
											  Two4x24   string `json:"24x24"`
											  Three2x32 string `json:"32x32"`
											  Four8x48  string `json:"48x48"`
										  } `json:"avatarUrls"`
								     DisplayName  string `json:"displayName"`
								     EmailAddress string `json:"emailAddress"`
								     Name         string `json:"name"`
								     Self         string `json:"self"`
							     } `json:"creator"`
			       Customfield10000              interface{} `json:"customfield_10000"`
			       Customfield10010              interface{} `json:"customfield_10010"`
			       Customfield10041              struct {
								     ID    string `json:"id"`
								     Self  string `json:"self"`
								     Value string `json:"value"`
							     } `json:"customfield_10041"`
			       Customfield10042              interface{} `json:"customfield_10042"`
			       Customfield10043              interface{} `json:"customfield_10043"`
			       Customfield10045              interface{} `json:"customfield_10045"`
			       Customfield10046              interface{} `json:"customfield_10046"`
			       Customfield10047              interface{} `json:"customfield_10047"`
			       Customfield10049              interface{} `json:"customfield_10049"`
			       Customfield10050              interface{} `json:"customfield_10050"`
			       Customfield10051              interface{} `json:"customfield_10051"`
			       Customfield10055              interface{} `json:"customfield_10055"`
			       Customfield10057              interface{} `json:"customfield_10057"`
			       Customfield10060              interface{} `json:"customfield_10060"`
			       Customfield10061              interface{} `json:"customfield_10061"`
			       Customfield10062              interface{} `json:"customfield_10062"`
			       Customfield10063              interface{} `json:"customfield_10063"`
			       Customfield10072              struct {
								     ID    string `json:"id"`
								     Self  string `json:"self"`
								     Value string `json:"value"`
							     } `json:"customfield_10072"`
			       Customfield10073              struct {
								     ID    string `json:"id"`
								     Self  string `json:"self"`
								     Value string `json:"value"`
							     } `json:"customfield_10073"`
			       Customfield10081              interface{} `json:"customfield_10081"`
			       Customfield10082              interface{} `json:"customfield_10082"`
			       Customfield10084              interface{} `json:"customfield_10084"`
			       Customfield10100              interface{} `json:"customfield_10100"`
			       Customfield10130              interface{} `json:"customfield_10130"`
			       Customfield10360              interface{} `json:"customfield_10360"`
			       Customfield10361              string      `json:"customfield_10361"`
			       Customfield10363              interface{} `json:"customfield_10363"`
			       Customfield10366              interface{} `json:"customfield_10366"`
			       Customfield10660              interface{} `json:"customfield_10660"`
			       Customfield10860              interface{} `json:"customfield_10860"`
			       Customfield10960              string      `json:"customfield_10960"`
			       Customfield11060              interface{} `json:"customfield_11060"`
			       Description                   interface{} `json:"description"`
			       Duedate                       interface{} `json:"duedate"`
			       Environment                   interface{} `json:"environment"`
			       FixVersions                   []struct {
				       Archived bool   `json:"archived"`
				       ID       string `json:"id"`
				       Name     string `json:"name"`
				       Released bool   `json:"released"`
				       Self     string `json:"self"`
			       } `json:"fixVersions"`
			       Issuelinks                    []interface{} `json:"issuelinks"`
			       Issuetype                     struct {
								     Description string `json:"description"`
								     IconURL     string `json:"iconUrl"`
								     ID          string `json:"id"`
								     Name        string `json:"name"`
								     Self        string `json:"self"`
								     Subtask     bool   `json:"subtask"`
							     } `json:"issuetype"`
			       Labels                        []interface{} `json:"labels"`
			       LastViewed                    interface{}   `json:"lastViewed"`
			       Parent                        struct {
								     Fields struct {
										    Issuetype struct {
												      Description string `json:"description"`
												      IconURL     string `json:"iconUrl"`
												      ID          string `json:"id"`
												      Name        string `json:"name"`
												      Self        string `json:"self"`
												      Subtask     bool   `json:"subtask"`
											      } `json:"issuetype"`
										    Priority  struct {
												      IconURL string `json:"iconUrl"`
												      ID      string `json:"id"`
												      Name    string `json:"name"`
												      Self    string `json:"self"`
											      } `json:"priority"`
										    Status    struct {
												      Description    string `json:"description"`
												      IconURL        string `json:"iconUrl"`
												      ID             string `json:"id"`
												      Name           string `json:"name"`
												      Self           string `json:"self"`
												      StatusCategory struct {
															     ColorName string `json:"colorName"`
															     ID        int    `json:"id"`
															     Key       string `json:"key"`
															     Name      string `json:"name"`
															     Self      string `json:"self"`
														     } `json:"statusCategory"`
											      } `json:"status"`
										    Summary   string `json:"summary"`
									    } `json:"fields"`
								     ID     string `json:"id"`
								     Key    string `json:"key"`
								     Self   string `json:"self"`
							     } `json:"parent"`
			       Priority                      struct {
								     IconURL string `json:"iconUrl"`
								     ID      string `json:"id"`
								     Name    string `json:"name"`
								     Self    string `json:"self"`
							     } `json:"priority"`
			       Progress                      struct {
								     Percent  int `json:"percent"`
								     Progress int `json:"progress"`
								     Total    int `json:"total"`
							     } `json:"progress"`
			       Project                       struct {
								     AvatarUrls struct {
											One6x16   string `json:"16x16"`
											Two4x24   string `json:"24x24"`
											Three2x32 string `json:"32x32"`
											Four8x48  string `json:"48x48"`
										} `json:"avatarUrls"`
								     ID         string `json:"id"`
								     Key        string `json:"key"`
								     Name       string `json:"name"`
								     Self       string `json:"self"`
							     } `json:"project"`
			       Reporter                      struct {
								     Active       bool `json:"active"`
								     AvatarUrls   struct {
											  One6x16   string `json:"16x16"`
											  Two4x24   string `json:"24x24"`
											  Three2x32 string `json:"32x32"`
											  Four8x48  string `json:"48x48"`
										  } `json:"avatarUrls"`
								     DisplayName  string `json:"displayName"`
								     EmailAddress string `json:"emailAddress"`
								     Name         string `json:"name"`
								     Self         string `json:"self"`
							     } `json:"reporter"`
			       Resolution                    struct {
								     Description string `json:"description"`
								     ID          string `json:"id"`
								     Name        string `json:"name"`
								     Self        string `json:"self"`
							     } `json:"resolution"`
			       Resolutiondate                string `json:"resolutiondate"`
			       Status                        struct {
								     Description    string `json:"description"`
								     IconURL        string `json:"iconUrl"`
								     ID             string `json:"id"`
								     Name           string `json:"name"`
								     Self           string `json:"self"`
								     StatusCategory struct {
											    ColorName string `json:"colorName"`
											    ID        int    `json:"id"`
											    Key       string `json:"key"`
											    Name      string `json:"name"`
											    Self      string `json:"self"`
										    } `json:"statusCategory"`
							     } `json:"status"`
			       Subtasks                      []interface{} `json:"subtasks"`
			       Summary                       string        `json:"summary"`
			       Timeestimate                  int           `json:"timeestimate"`
			       Timeoriginalestimate          int           `json:"timeoriginalestimate"`
			       Timespent                     int           `json:"timespent"`
			       Updated                       string        `json:"updated"`
			       Versions                      []interface{} `json:"versions"`
			       Votes                         struct {
								     HasVoted bool   `json:"hasVoted"`
								     Self     string `json:"self"`
								     Votes    int    `json:"votes"`
							     } `json:"votes"`
			       Watches                       struct {
								     IsWatching bool   `json:"isWatching"`
								     Self       string `json:"self"`
								     WatchCount int    `json:"watchCount"`
							     } `json:"watches"`
			       Workratio                     int `json:"workratio"`
		       } `json:"fields"`
		ID     string `json:"id"`
		Key    string `json:"key"`
		Self   string `json:"self"`
	} `json:"issues"`
	MaxResults int `json:"maxResults"`
	StartAt    int `json:"startAt"`
	Total      int `json:"total"`
}
type JiraIssue struct {
	Expand string `json:"expand"`
	ID     string `json:"id"`
	Self   string `json:"self"`
	Key    string `json:"key"`
	Fields struct {
		       Issuetype                     struct {
							     Self        string `json:"self"`
							     ID          string `json:"id"`
							     Description string `json:"description"`
							     IconURL     string `json:"iconUrl"`
							     Name        string `json:"name"`
							     Subtask     bool `json:"subtask"`
						     } `json:"issuetype"`
		       Customfield10072              struct {
							     Self  string `json:"self"`
							     Value string `json:"value"`
							     ID    string `json:"id"`
						     } `json:"customfield_10072"`
		       Timespent                     interface{} `json:"timespent"`
		       Customfield10073              struct {
							     Self  string `json:"self"`
							     Value string `json:"value"`
							     ID    string `json:"id"`
						     } `json:"customfield_10073"`
		       Project                       struct {
							     Self       string `json:"self"`
							     ID         string `json:"id"`
							     Key        string `json:"key"`
							     Name       string `json:"name"`
							     AvatarUrls struct {
										Four8X48  string `json:"48x48"`
										Two4X24   string `json:"24x24"`
										One6X16   string `json:"16x16"`
										Three2X32 string `json:"32x32"`
									} `json:"avatarUrls"`
						     } `json:"project"`
		       FixVersions                   []struct {
			       Self        string `json:"self"`
			       ID          string `json:"id"`
			       Description string `json:"description"`
			       Name        string `json:"name"`
			       Archived    bool `json:"archived"`
			       Released    bool `json:"released"`
		       } `json:"fixVersions"`
		       Aggregatetimespent            interface{} `json:"aggregatetimespent"`
		       Resolution                    struct {
							     Self        string `json:"self"`
							     ID          string `json:"id"`
							     Description string `json:"description"`
							     Name        string `json:"name"`
						     } `json:"resolution"`
		       Resolutiondate                string `json:"resolutiondate"`
		       Workratio                     int `json:"workratio"`
		       LastViewed                    string `json:"lastViewed"`
		       Watches                       struct {
							     Self       string `json:"self"`
							     WatchCount int `json:"watchCount"`
							     IsWatching bool `json:"isWatching"`
						     } `json:"watches"`
		       Customfield10060              interface{} `json:"customfield_10060"`
		       Customfield10061              interface{} `json:"customfield_10061"`
		       Created                       string `json:"created"`
		       Customfield10062              interface{} `json:"customfield_10062"`
		       Customfield10063              interface{} `json:"customfield_10063"`
		       Customfield10660              interface{} `json:"customfield_10660"`
		       Priority                      struct {
							     Self    string `json:"self"`
							     IconURL string `json:"iconUrl"`
							     Name    string `json:"name"`
							     ID      string `json:"id"`
						     } `json:"priority"`
		       Customfield10100              interface{} `json:"customfield_10100"`
		       Customfield10860              interface{} `json:"customfield_10860"`
		       Customfield10101              []struct {
			       Self  string `json:"self"`
			       Value string `json:"value"`
			       ID    string `json:"id"`
		       } `json:"customfield_10101"`
		       Labels                        []string `json:"labels"`
		       Timeestimate                  interface{} `json:"timeestimate"`
		       Aggregatetimeoriginalestimate interface{} `json:"aggregatetimeoriginalestimate"`
		       Versions                      []struct {
			       Self        string `json:"self"`
			       ID          string `json:"id"`
			       Description string `json:"description,omitempty"`
			       Name        string `json:"name"`
			       Archived    bool `json:"archived"`
			       Released    bool `json:"released"`
		       } `json:"versions"`
		       Issuelinks                    []interface{} `json:"issuelinks"`
		       Assignee                      struct {
							     Self         string `json:"self"`
							     Name         string `json:"name"`
							     EmailAddress string `json:"emailAddress"`
							     AvatarUrls   struct {
										  Four8X48  string `json:"48x48"`
										  Two4X24   string `json:"24x24"`
										  One6X16   string `json:"16x16"`
										  Three2X32 string `json:"32x32"`
									  } `json:"avatarUrls"`
							     DisplayName  string `json:"displayName"`
							     Active       bool `json:"active"`
						     } `json:"assignee"`
		       Updated                       string `json:"updated"`
		       Status                        struct {
							     Self           string `json:"self"`
							     Description    string `json:"description"`
							     IconURL        string `json:"iconUrl"`
							     Name           string `json:"name"`
							     ID             string `json:"id"`
							     StatusCategory struct {
										    Self      string `json:"self"`
										    ID        int `json:"id"`
										    Key       string `json:"key"`
										    ColorName string `json:"colorName"`
										    Name      string `json:"name"`
									    } `json:"statusCategory"`
						     } `json:"status"`
		       Components                    []struct {
			       Self        string `json:"self"`
			       ID          string `json:"id"`
			       Name        string `json:"name"`
			       Description string `json:"description"`
		       } `json:"components"`
		       Customfield11060              interface{} `json:"customfield_11060"`
		       Customfield10050              interface{} `json:"customfield_10050"`
		       Customfield10051              interface{} `json:"customfield_10051"`
		       Timeoriginalestimate          interface{} `json:"timeoriginalestimate"`
		       Description                   string `json:"description"`
		       Customfield10130              interface{} `json:"customfield_10130"`
		       Customfield10010              interface{} `json:"customfield_10010"`
		       Customfield10055              interface{} `json:"customfield_10055"`
		       Customfield10057              interface{} `json:"customfield_10057"`
		       Timetracking                  struct {
						     } `json:"timetracking"`
		       Customfield10049              interface{} `json:"customfield_10049"`
		       Attachment                    []interface{} `json:"attachment"`
		       Aggregatetimeestimate         interface{} `json:"aggregatetimeestimate"`
		       Summary                       string `json:"summary"`
		       Creator                       struct {
							     Self         string `json:"self"`
							     Name         string `json:"name"`
							     EmailAddress string `json:"emailAddress"`
							     AvatarUrls   struct {
										  Four8X48  string `json:"48x48"`
										  Two4X24   string `json:"24x24"`
										  One6X16   string `json:"16x16"`
										  Three2X32 string `json:"32x32"`
									  } `json:"avatarUrls"`
							     DisplayName  string `json:"displayName"`
							     Active       bool `json:"active"`
						     } `json:"creator"`
		       Customfield10082              interface{} `json:"customfield_10082"`
		       Subtasks                      []interface{} `json:"subtasks"`
		       Customfield10084              interface{} `json:"customfield_10084"`
		       Customfield10360              interface{} `json:"customfield_10360"`
		       Customfield10041              struct {
							     Self  string `json:"self"`
							     Value string `json:"value"`
							     ID    string `json:"id"`
						     } `json:"customfield_10041"`
		       Customfield10361              string `json:"customfield_10361"`
		       Customfield10042              interface{} `json:"customfield_10042"`
		       Reporter                      struct {
							     Self         string `json:"self"`
							     Name         string `json:"name"`
							     EmailAddress string `json:"emailAddress"`
							     AvatarUrls   struct {
										  Four8X48  string `json:"48x48"`
										  Two4X24   string `json:"24x24"`
										  One6X16   string `json:"16x16"`
										  Three2X32 string `json:"32x32"`
									  } `json:"avatarUrls"`
							     DisplayName  string `json:"displayName"`
							     Active       bool `json:"active"`
						     } `json:"reporter"`
		       Customfield10043              interface{} `json:"customfield_10043"`
		       Customfield10363              interface{} `json:"customfield_10363"`
		       Customfield10000              interface{} `json:"customfield_10000"`
		       Aggregateprogress             struct {
							     Progress int `json:"progress"`
							     Total    int `json:"total"`
						     } `json:"aggregateprogress"`
		       Customfield10045              interface{} `json:"customfield_10045"`
		       Customfield10046              interface{} `json:"customfield_10046"`
		       Customfield10366              interface{} `json:"customfield_10366"`
		       Customfield10047              interface{} `json:"customfield_10047"`
		       Customfield10960              string `json:"customfield_10960"`
		       Customfield10125              struct {
							     Self  string `json:"self"`
							     Value string `json:"value"`
							     ID    string `json:"id"`
						     } `json:"customfield_10125"`
		       Environment                   interface{} `json:"environment"`
		       Duedate                       interface{} `json:"duedate"`
		       Progress                      struct {
							     Progress int `json:"progress"`
							     Total    int `json:"total"`
						     } `json:"progress"`
		       Comment                       struct {
							     StartAt    int `json:"startAt"`
							     MaxResults int `json:"maxResults"`
							     Total      int `json:"total"`
							     Comments   []struct {
								     Self         string `json:"self"`
								     ID           string `json:"id"`
								     Author       struct {
											  Self         string `json:"self"`
											  Name         string `json:"name"`
											  EmailAddress string `json:"emailAddress"`
											  AvatarUrls   struct {
													       Four8X48  string `json:"48x48"`
													       Two4X24   string `json:"24x24"`
													       One6X16   string `json:"16x16"`
													       Three2X32 string `json:"32x32"`
												       } `json:"avatarUrls"`
											  DisplayName  string `json:"displayName"`
											  Active       bool `json:"active"`
										  } `json:"author"`
								     Body         string `json:"body"`
								     UpdateAuthor struct {
											  Self         string `json:"self"`
											  Name         string `json:"name"`
											  EmailAddress string `json:"emailAddress"`
											  AvatarUrls   struct {
													       Four8X48  string `json:"48x48"`
													       Two4X24   string `json:"24x24"`
													       One6X16   string `json:"16x16"`
													       Three2X32 string `json:"32x32"`
												       } `json:"avatarUrls"`
											  DisplayName  string `json:"displayName"`
											  Active       bool `json:"active"`
										  } `json:"updateAuthor"`
								     Created      string `json:"created"`
								     Updated      string `json:"updated"`
							     } `json:"comments"`
						     } `json:"comment"`
		       Votes                         struct {
							     Self     string `json:"self"`
							     Votes    int `json:"votes"`
							     HasVoted bool `json:"hasVoted"`
						     } `json:"votes"`
		       Worklog                       struct {
							     StartAt    int `json:"startAt"`
							     MaxResults int `json:"maxResults"`
							     Total      int `json:"total"`
							     Worklogs   []interface{} `json:"worklogs"`
						     } `json:"worklog"`
	       } `json:"fields"`
}

func (jira *Jira )printIssues() {
	for _, issue := range jira.Issues {
		var fix = ""
		for _, fixversion := range issue.Fields.FixVersions {
			if (len(fix) == 0) {
				fix += fixversion.Name
			} else {
				fix += ", " + fixversion.Name
			}
		}
		var priority = ""
		if issue.Fields.Priority.Name == "Minor" {
			priority = fmt.Sprintf("\033[0;32m%-10s\033[m", issue.Fields.Priority.Name)
		} else if issue.Fields.Priority.Name == "Major" {
			priority = fmt.Sprintf("\033[0;31m%-10s\033[m", issue.Fields.Priority.Name)
		} else if issue.Fields.Priority.Name == "Blocker" {
			priority = fmt.Sprintf("\033[0;30;41m%-10s\033[m", issue.Fields.Priority.Name)
		} else {
			priority = fmt.Sprintf("%s", issue.Fields.Priority.Name)
		}
		assignee := issue.Fields.Assignee.Name
		creator := issue.Fields.Creator.Name
		if assignee == me {
			assignee = fmt.Sprintf("\033[1;31m%-10s\033[m", me)
		}
		if creator == me {
			assignee = fmt.Sprintf("\033[1;31m-10%s\033[m", me)
		}
		fmt.Printf("%-15s %-15s %-10s %-10s %-10s %-20s %-20s %s \033[34m%-10s\033[m %s\n",
			issue.Key, issue.Fields.Created[:len("2006-01-02T15:04:05")], priority, assignee, creator, fix, issue.Fields.Status.Name, issue.Fields.Summary,
			"http://" + jiraServer + "/browse/" + issue.Key, "jira " + issue.Key)
	}
	fmt.Printf("%d %s\n", len(jira.Issues), "issues found")
}

var ci = flag.Bool("ci", false, "Create a new issue")
var mine = flag.Bool("mine", false, "List my issues")
var me = ""
var creds = me + ":"
var jiraServer = ""

func er(err error) {
	if (err != nil) {
		log.Panic(err)
		os.Exit(-1)
	}
}
func main() {
	flag.Parse()
	version := getNextFixVersion()
	if version != "" {
		fmt.Printf("Fixversion: %s\n", version)
	}
	if *ci {
		return
	}
	if *mine {
		printMine()
		return
	}
	if len(flag.Args()) == 1 {
		printIssue(flag.Arg(0))
		return
	}
	printDefaultFilter()
	printTestFilter()
}

func printIssue(issueName string) {
	resp, err := http.DefaultClient.Get("http://" + creds + "@" + jiraServer + "/rest/api/2/issue/" + issueName)
	er(err)
	bytes, err := ioutil.ReadAll(resp.Body)
	er(err)
	var jiraIssue JiraIssue
	err = json.Unmarshal(bytes, &jiraIssue)
	er(err)
	var fix = ""
	for _, fixversion := range jiraIssue.Fields.FixVersions {
		if (len(fix) == 0) {
			fix += fixversion.Name
		} else {
			fix += ", " + fixversion.Name
		}
	}
	fmt.Printf("\033[32m%-10s\033[m ", jiraIssue.Fields.Created)
	fmt.Printf("\033[33m%-10s\033[m ", jiraIssue.Fields.Status.Name)
	fmt.Printf("\033[34m%-10s\033[m ", jiraIssue.Fields.Creator.Name)
	fmt.Printf("\033[35m%-10s\033[m ", jiraIssue.Fields.Assignee.Name)
	fmt.Printf("\033[36m%-10s\033[m ", fix)
	fmt.Printf("\n%s\n", jiraIssue.Fields.Summary)
	fmt.Print("\n", )
	fmt.Printf("%s\n", jiraIssue.Fields.Description)

	if len(jiraIssue.Fields.Comment.Comments) > 0 {
		fmt.Printf("%s\n", "Comments:")
		for _, comment := range jiraIssue.Fields.Comment.Comments {
			fmt.Printf("%s \033[0;36m%s\033[m %s\n", comment.Created, comment.Author.Name, comment.Body)
		}
	}

	fmt.Printf("\033[34m%-10s\033[m\n", "http://" + jiraServer + "/browse/" + issueName)
}
func printMine() {
	resp, err := http.DefaultClient.Get("http://" + creds + "@" + jiraServer + "/rest/api/2/search?jql=" + url.QueryEscape("assignee = " + me + " AND status not in (Lukket, Closed, Resolved) "))
	er(err)
	bytes, err := ioutil.ReadAll(resp.Body)
	er(err)
	var jira Jira
	err = json.Unmarshal(bytes, &jira)
	er(err)
	jira.printIssues()
}
func printDefaultFilter() {
	fmt.Print("Rettes:\n")
	resp, err := http.DefaultClient.Get("http://" + creds + "@" + jiraServer + "/rest/api/2/search?jql=filter=16228")
	er(err)
	bytes, err := ioutil.ReadAll(resp.Body)
	er(err)
	var jira Jira
	err = json.Unmarshal(bytes, &jira)
	er(err)
	jira.printIssues()
}
func printTestFilter() {
	fmt.Print("Testes:\n")
	resp, err := http.DefaultClient.Get("http://" + creds + "@" + jiraServer + "/rest/api/2/search?jql=filter=16319")
	er(err)
	bytes, err := ioutil.ReadAll(resp.Body)
	er(err)
	var jira Jira
	err = json.Unmarshal(bytes, &jira)
	er(err)
	jira.printIssues()
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
	er(err)
	return pom.Value[:strings.Index(pom.Value, "-")]
}

