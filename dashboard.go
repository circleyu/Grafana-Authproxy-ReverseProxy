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
		log.Print(err)
		return data, err
	}
	log.Print("CreateDashboard: " + data.Status)

	return data, nil
}
