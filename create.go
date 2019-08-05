package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dghubble/sling"
)

type createData struct {
	UserName  string `json:"userName"`
	Email     string `json:"email"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

func createHandler(w http.ResponseWriter, r *http.Request) {

	s := fmt.Sprintf("admin:%s", adminPassWord)
	base64String := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(s)))

	slingClient := sling.New().Base(grafanaURL).Client(nil)
	slingClient.Set("Accept", "application/json")
	slingClient.Set("Content-Type", "application/json")
	slingClient.Set("Authorization", base64String)

	// 這邊傳進來的時候會帶著accessKey、secretKey、userName、email
	buf := new(bytes.Buffer)
	io.Copy(buf, r.Body)
	r.Body.Close()
	var data createData
	json.Unmarshal(buf.Bytes(), &data)

	userID, err := createUser(slingClient, data.UserName, data.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	orgID, err := createOrg(slingClient, data.UserName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = addUserToOrg(slingClient, orgID, data.UserName, userID)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	err = switchToOrg(slingClient, orgID)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	err = createDataSource(slingClient, data.AccessKey, data.SecretKey)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	respData, err := createDashboard(slingClient, "ec2")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// reset org
	err = switchToOrg(slingClient, 1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	b, err := json.Marshal(respData)
	w.Write(b)
}
