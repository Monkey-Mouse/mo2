package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	EmailAddr string
	EmailPass string
}

func LoadConfig(filename string) (c Configuration) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		panic(err)
	}
	return
}
