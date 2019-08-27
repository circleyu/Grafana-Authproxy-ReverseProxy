package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dghubble/sling"
)

type createDashboardResp struct {
	ID      int    `json:"id"`
	Slug    string `json:"slug"`
	Status  string `json:"status"`
	UID     string `json:"uid"`
	URL     string `json:"url"`
	Version int    `json:"version"`
}

func createDashboard(slingClient *sling.Sling, dashboard string) (createDashboardResp, error) {
	var data createDashboardResp

	f, err := os.Open(fmt.Sprintf("./dashboard/%s.json", dashboard))
	if err != nil {
		return data, err
	}
	defer f.Close()

	path := "/api/dashboards/db"

	var err2 error
	_, err = slingClient.New().Post(path).Body(f).Receive(&data, err2)
	if err == nil {
		err = err2
	}
	if err != nil {
		log.Printf("Create %v Error: %q", dashboard, err)
		return data, err
	}
	log.Printf("CreateDashboard: %v %v", dashboard, data.Status)

	return data, nil
}

type getResp struct {
	Dashboard myDashboard `json:"dashboard"`
	Meta      myMeta      `json:"meta"`
}
type myDashboard struct {
	ID            int      `json:"id"`
	UID           string   `json:"uid"`
	Title         string   `json:"title"`
	Tags          []string `json:"tags"`
	Timezone      string   `json:"timezone"`
	SchemaVersion int      `json:"schemaVersion"`
	Version       int      `json:"version"`
}
type myMeta struct {
	IsStarred bool   `json:"isStarred"`
	URL       string `json:"url"`
	Slug      string `json:"slug"`
}

func hasDashboard(slingClient *sling.Sling, uid string) (bool, error) {
	path := "/api/dashboards/uid/" + uid
	var data getResp
	var err2 error
	r, err := slingClient.New().Get(path).Receive(&data, err2)
	if err == nil {
		err = err2
	}
	if err != nil {
		log.Print(err)
		return false, err
	}
	if r.StatusCode != 200 {
		return false, err
	}
	return true, nil
}
