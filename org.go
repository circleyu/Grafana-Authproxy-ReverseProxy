package main

import (
	"fmt"
	"log"

	"github.com/dghubble/sling"
)

type createOrgReq struct {
	OrgName string `json:"name"`
}

type createOrgResp struct {
	OrgID   int    `json:"orgId"`
	Message string `json:"message"`
}

func createOrg(slingClient *sling.Sling, orgName string) (int, error) {
	path := "/api/orgs"
	body := &createOrgReq{
		OrgName: orgName,
	}
	var data createOrgResp
	var err2 error
	_, err := slingClient.New().Post(path).BodyJSON(body).Receive(&data, err2)
	if err == nil {
		err = err2
	}
	if err != nil {
		log.Print(err)
		return -1, err
	}
	log.Print(data.Message)

	return data.OrgID, nil
}

type addUserToOrgReq struct {
	LoginOrEmail string `json:"loginOrEmail"`
	Role         string `json:"role"`
}

type addUserToOrgResp struct {
	Message string `json:"message"`
}

func addUserToOrg(slingClient *sling.Sling, orgID int, userName string, userID int) error {
	path := fmt.Sprintf("/api/orgs/%d/users", orgID)
	body := &addUserToOrgReq{
		LoginOrEmail: userName,
		// 幫客戶開通的權限只會給Viewer
		Role: "Viewer",
	}
	var data addUserToOrgResp
	var err2 error
	_, err := slingClient.New().Post(path).BodyJSON(body).Receive(&data, err2)
	if err == nil {
		err = err2
	}
	if err != nil {
		log.Print(err)
		return err
	}
	log.Print(data.Message)

	// 將user從main org移除
	path = fmt.Sprintf("/api/orgs/1/users/%d", userID)
	_, err = slingClient.New().Delete(path).Receive(&data, err2)
	if err == nil {
		err = err2
	}
	if err != nil {
		log.Print(err)
		return err
	}
	log.Print(data.Message)
	return nil
}

type switchToOrgResp struct {
	Message string `json:"message"`
}

func switchToOrg(slingClient *sling.Sling, orgID int) error {
	path := fmt.Sprintf("/api/user/using/%d", orgID)
	var data switchToOrgResp
	var err2 error

	_, err := slingClient.New().Post(path).Receive(&data, err2)

	if err == nil {
		err = err2
	}
	if err != nil {
		log.Print(err)
		return err
	}

	log.Print(data.Message)

	return nil
}
