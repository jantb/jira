package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Confluence struct {
	Results []struct {
		ID   string `json:"id"`
		Type string `json:"type"`
		Body struct {
			Expandable struct {
				Storage string `json:"storage"`
			} `json:"_expandable"`
		} `json:"body"`
	} `json:"results"`
	Start int `json:"start"`
	Limit int `json:"limit"`
	Size  int `json:"size"`
}

func basicAuth(start int) string {
	client := &http.Client{}

	urll := conf.ConfluenceServer + "rest/api/content?spaceKey=" + conf.ProjectConfluence + "&expand=body.storage&start=" + fmt.Sprintf("%d", start)
	fmt.Println(urll)
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

	s := string(bodyText)
	return s
}
