package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"net/url"

	"github.com/jantb/jira/strip"
)

// Confluence struct for holding pages from confluence
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

// Page struct to hold info
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
		urll := conf.ConfluenceServer + "rest/api/content/search?cql=" + url.QueryEscape("space in ("+conf.ProjectConfluence+") and lastModified > \""+conf.LastUpdateConfluence.Format("2006/01/02")+"\"") + "&expand=body.storage.content&start=" + fmt.Sprintf("%d", i)
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
			conf.LastUpdateConfluence = time.Now()
			conf.store()
			break
		}
		for j, result := range confluence.Results {
			fmt.Printf("\r%d changed/new confluence pages", i+j)
			text := strip.StripTags(result.Body.Storage.Value)
			for key, value := range xml.HTMLEntity {
				text = strings.Replace(text, "&"+key+";", value, -1)
			}
			text = strings.ToLower(text)
			pages = append(pages, Page{Key: result.ID, Body: text, Title: result.Title, Link: conf.ConfluenceServer + result.Links.Webui[1:]})
		}
	}
	if len(pages) > 0 {
		fmt.Println()
	}

	return pages
}
