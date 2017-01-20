package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/jantb/jira/strip"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Confluence struct {
	Results []struct {
		ID     string `json:"id"`
		Type   string `json:"type"`
		Status string `json:"status"`
		Title  string `json:"title"`
		Body   struct {
			Storage struct {
				Value          string `json:"value"`
				Representation string `json:"representation"`
			} `json:"storage"`
		} `json:"body"`
		Links struct {
			Webui  string `json:"webui"`
			Tinyui string `json:"tinyui"`
			Self   string `json:"self"`
		} `json:"_links"`
	} `json:"results"`
	Start int `json:"start"`
	Limit int `json:"limit"`
	Size  int `json:"size"`
}
type Page struct {
	Key   string
	Body  string
	Link  string
	Title string
}

func getConfluencePages() []Page {
	client := &http.Client{}
	pages := []Page{}
	for i := 0; ; {
		urll := conf.ConfluenceServer + "rest/api/content?spaceKey=" + conf.ProjectConfluence + "&expand=body.storage.content&start=" + fmt.Sprintf("%d", i)
		req, err := http.NewRequest("GET", urll, nil)
		req.SetBasicAuth(conf.UsernameConfluence, conf.PasswordConfluence)
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		bodyText, err := ioutil.ReadAll(resp.Body)
		confluence := Confluence{}
		err = json.Unmarshal(bodyText, &confluence)
		if err != nil {
			fmt.Println(string(bodyText))
			log.Fatal(err)
		}
		i += confluence.Size
		if confluence.Size == 0 {
			break
		}
		for _, result := range confluence.Results {
			text := strip.StripTags(result.Body.Storage.Value)
			for key, value := range xml.HTMLEntity {
				text = strings.Replace(text, "&"+key+";", value, -1)
			}
			pages = append(pages, Page{Key: result.ID, Body: text, Title: result.Title, Link: conf.ConfluenceServer + result.Links.Webui[1:]})
		}
	}
	return pages
}
