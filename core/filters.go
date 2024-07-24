package core

import (
	"fmt"
	"log"
	"os"
	// "sync"
	"gopkg.in/yaml.v2"
)

type Filter struct {
	Name string 			`yaml:"name"`
	Scenarios []string 		`yaml:"scenarios"`
	MatchScenarios []string `yaml:"match-scenarios"`
	Ipset string 			`yaml:"ipset"`
	Type string 			`yaml:"type"`
	Append bool				`yaml:"append"`
}

type DispatchTable []struct {
	Name string
	Scenarios []string
	MatchScenarios []string
	Channel chan byte
}

func parseFiltersConfig(path string) (filters []Filter) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	log.Printf("found filters configuration file %v", path)

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&filters)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("successfully parsed filters configuration from %v", path)
	fmt.Println(filters)
	return
}

func spawnFilters(filters []Filters) (dt DispatchTable) {
	for filter := range filters {

	}
}


