package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-yaml/yaml"
)

type Translation = map[string]map[string]string

type Fragment struct {
	Name        string      `yaml:"name"`
	Translation Translation `yaml:"translations"`
}

type translationList struct {
	Version   string     `yaml:"version"`
	Fragments []Fragment `yaml:"fragments"`
}

func main() {
	translationFile, err := os.Open("translation.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer translationFile.Close()
	content, err := ioutil.ReadAll(translationFile)
	if err != nil {
		log.Fatal(err)
	}
	var list translationList
	err = yaml.Unmarshal(content, &list)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v", list)
}
