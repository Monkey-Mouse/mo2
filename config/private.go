package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func SaveConfig(filename string) {
	config := Configuration{
		EmailAddr: "XXXXXXXXXXXXX",
		EmailPass: "XXXXXXXXXXXXX",
	}
	bytes, err := yaml.Marshal(config)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(filename, bytes, 0644)
}
