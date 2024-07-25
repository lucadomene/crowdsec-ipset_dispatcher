package core

import (
	"fmt"
	"log"
	"lucadomeneghetti/ipset_dispatcher/models"
	"os"
	"regexp"

	"github.com/janeczku/go-ipset/ipset"
	"gopkg.in/yaml.v2"
)

type Filter struct {
	Name           string   `yaml:"name"`
	Scenarios      []string `yaml:"scenarios"`
	MatchScenarios []string `yaml:"match-scenarios"`
	Ipset          string   `yaml:"ipset"`
	Type           string   `yaml:"type"`
}

func InitializeFilters(path string) {
	filters := parseFiltersConfig(path)
	spawnFilters(filters)
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

func spawnFilters(filters []Filter) {
	channels := make([]chan models.Decision, len(filters))
	for i := range channels {
		channels[i] = make(chan models.Decision)
	}
	for i, filter := range filters {
		go runningFilter(filter, channels[i])
	}
}

func filterDecision(filter Filter, decision models.Decision) bool {
	matches := make([]*regexp.Regexp, len(filter.MatchScenarios))
	for _, match := range matches {
		if match.MatchString(decision.Scenario) {
			return true
		}
	}

	for _, scenario := range filter.Scenarios {
		if scenario == decision.Scenario {
			return true
		}
	}
	return false
}

func runningFilter(filter Filter, ch <-chan models.Decision) {
	filterIpset, err := ipset.New(filter.Ipset, filter.Type, &ipset.Params{})
	if err != nil {
		log.Fatalf("failed to create ipset %v: %v", filter.Ipset, err)
	}

	for decision := range ch {
		if filterDecision(filter, decision) {
			switch decsion.Action {
			case "add":
				filterIpset.Add(decision.IPAddress, decision.Duration)
			case "del":
				filterIpset.Del(decision.IPAddress)
			}
		}
	}
}
