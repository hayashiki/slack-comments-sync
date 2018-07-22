package main

import (
	"encoding/json"
	"io/ioutil"
)

var (
	Version = "0.0.1"
)

type Config struct {
	Accounts     map[string]string `json:"accounts"`
	Repositories map[string]string `json:"repositories"`
}

func ParseFile(filename string) (*Config, error) {
	c := Config{}

	jsonString, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonString, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
