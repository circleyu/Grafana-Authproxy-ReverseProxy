package main

import (
	"log"

	"github.com/dghubble/sling"
)

type createUserReq struct {
	UserName string `json:"name"`
	Email    string `json:"email"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type createUserResp struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

func createUser(slingClient *sling.Sling, userName string, email string) (int, error) {
	path := "/api/admin/users"

	body := &createUserReq{
		UserName: userName,
		Email:    email,
		Login:    userName,
		Password: userName,
	}

	var data createUserResp
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

	return data.ID, nil
}
