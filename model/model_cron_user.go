package model

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// type Body struct {
// 	User string `json:"user,omitempty"`
// 	IP   string `json:"mac,omitempty"`
// }

func CronUpdateUser() error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	response, err := client.Get("http://172.16.202.100/config/nac-user/listUserIP.json")
	if err != nil {
		return err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var body []User
	err = json.Unmarshal(contents, &body)
	if err != nil {
		return err
	}

	err = InsertManyUser(body)
	if err != nil {
		return err
	}
	err = SetRoleAdmin()
	if err != nil {
		return err
	}

	return nil
}
