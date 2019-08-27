package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/dghubble/sling"
)

var mapDashboardUID map[string]string

type addData struct {
	UserName  string   `json:"userName"`
	Dashboard []string `json:"dashboard"`
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	slingClient := sling.New().Base(grafanaURL).Client(nil)
	slingClient.Set("Accept", "application/json")
	slingClient.Set("Content-Type", "application/json")
	slingClient.SetBasicAuth("admin", adminPassWord)

	// 這邊傳進來的時候會帶著userName、dashboard
	buf := new(bytes.Buffer)
	io.Copy(buf, r.Body)
	r.Body.Close()
	var data addData
	json.Unmarshal(buf.Bytes(), &data)

	userData, err := getUser(slingClient, data.UserName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = switchToOrg(slingClient, userData.OrgID)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	isAll := false

	for _, dashboard := range data.Dashboard {
		if dashboard == "all" {
			isAll = true
			break
		}
	}

	if isAll {
		for dashboard, uid := range mapDashboardUID {
			if have, _ := hasDashboard(slingClient, uid); !have {
				_, err := createDashboard(slingClient, dashboard)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(err.Error()))
					return
				}
			}
		}
	} else {

		for _, dashboard := range data.Dashboard {
			if have, _ := hasDashboard(slingClient, mapDashboardUID[dashboard]); !have {
				_, err := createDashboard(slingClient, dashboard)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(err.Error()))
					return
				}
			}
		}
	}

	// reset org
	err = switchToOrg(slingClient, 1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}
