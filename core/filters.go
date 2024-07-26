package core

import (
	"log"
	"log/slog"
	"lucadomeneghetti/ipset_dispatcher/ipset"
	"lucadomeneghetti/ipset_dispatcher/models"
	"math"
	"os"
	"regexp"
	"time"

	"gopkg.in/yaml.v2"
)

type Filter struct {
	Name           string   `yaml:"name"`
	Scenarios      []string `yaml:"scenarios"`
	MatchScenarios []string `yaml:"match-scenarios"`
	Type           string   `yaml:"type"`
	Ipset          string   `yaml:"ipset"`
	IpsetType      string   `yaml:"ipset-type"`
}

func InitializeFilters(path string) []chan models.Decision {
	filters := parseFiltersConfig(path)
	return spawnFilters(filters)
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
	return
}

func spawnFilters(filters []Filter) []chan models.Decision {
	channels := make([]chan models.Decision, len(filters))
	for i := range channels {
		channels[i] = make(chan models.Decision)
	}
	for i, filter := range filters {
		runningFilter(filter, channels[i])
	}
	return channels
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

	if len(filter.Type) > 0 {
		if decision.Type == filter.Type {
			return true
		}
	}
	return false
}

func runningFilter(filter Filter, ch <-chan models.Decision) {
	filterIpset := ipset.CreateSet(filter.Ipset, filter.IpsetType)

	for decision := range ch {
		if filterDecision(filter, decision) {
			switch decision.Action {
			case "add":
				duration_int, err := time.ParseDuration(decision.Duration)
				if err != nil {
					slog.Warn(err.Error())
					break
				}
				duration_uint32 := uint32(math.Round(duration_int.Seconds()))
				filterIpset.AddEntry(decision.IPAddress, duration_uint32)
			case "del":
				filterIpset.DeleteEntry(decision.IPAddress)
			}
		}
	}
}

func ForwardDecisions(decisions models.DecisionArray, channels []chan models.Decision) {
	for _, dec := range decisions {
		for _, ch := range channels {
			ch <- dec
		}
	}
}
