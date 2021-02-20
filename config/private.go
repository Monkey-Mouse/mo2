package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func SaveConfig(filename string) {
	config := Configuration{
		EmailAddr: "easilylazy@qq.com",
		EmailPass: "xmxviunohwbphadc",
	}
	bytes, err := yaml.Marshal(config)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile(filename, bytes, 0644)
}
