package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Configuration struct {
	EmailAddr string
	EmailPass string
}

func LoadConfig(filename string) (c Configuration) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		log.Fatal(err)
	}
	return
}
