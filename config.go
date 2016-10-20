package main

import (
	"time"
	"io/ioutil"
	"path/filepath"
	"encoding/json"
	"log"
	"os/user"
)

type config struct {
	Username   string
	Password   string
	JiraServer string
	Project string
	Filter string
	LastUpdate time.Time
}

func (c *config) load() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	bytes, err := ioutil.ReadFile(filepath.Join(usr.HomeDir, ".jira.conf"))
	if err != nil {
		bytes, err := json.MarshalIndent(c, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		ioutil.WriteFile(filepath.Join(usr.HomeDir, ".jira.conf"), bytes, 0600)
	}

	json.Unmarshal(bytes, &c)
}
func (c *config) store() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	bytes, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile(filepath.Join(usr.HomeDir, ".jira.conf"), bytes, 0600)
}
