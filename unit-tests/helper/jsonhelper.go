// Copyright 2021 Edgecast Inc. Licensed under the terms of the Apache 2.0 license.
// See LICENSE file in project root for terms.

package helper

import (
	"encoding/json"
	"io/ioutil"
	"terraform-provider-ec/unit-tests/model"
)

// Users struct which contains
// an array of users
type Users struct {
	Users []User `json:"users"`
}

// User struct which contains a name
// a type and a list of social links
type User struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Age  int    `json:"Age"`
}

func ReadCredentialJsonfile(path string, credential *model.Credentials) error {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(data), credential)
	if err != nil {
		panic(err)
	}
	return nil
}
