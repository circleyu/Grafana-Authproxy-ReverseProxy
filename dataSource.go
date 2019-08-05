package main

import (
	"log"

	"github.com/dghubble/sling"
)

type createDataSourceReq struct {
	Name           string         `json:"name"`
	Type           string         `json:"type"`
	Access         string         `json:"access"`
	JSONData       jsonData       `json:"jsonData"`
	SecureJSONData secureJSONData `json:"secureJsonData"`
}
type jsonData struct {
	AuthType      string `json:"authType"`
	DefaultRegion string `json:"defaultRegion"`
}

type secureJSONData struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

type createDataSourceResp struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
	Name    string `json:"name"`
}

func createDataSource(slingClient *sling.Sling, accessKey string, secretKey string) error {
	path := "/api/datasources"

	body := &createDataSourceReq{
		Name:   "CloudWatch",
		Type:   "cloudwatch",
		Access: "proxy",
		JSONData: jsonData{
			AuthType:      "keys",
			DefaultRegion: "us-west-1",
		},
		SecureJSONData: secureJSONData{
			AccessKey: accessKey,
			SecretKey: secretKey,
		},
	}

	var data createDataSourceResp
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

	return nil
}
