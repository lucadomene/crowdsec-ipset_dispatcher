package main

import (
	"log"
	"lucadomeneghetti/ipset_dispatcher/core"
	"lucadomeneghetti/ipset_dispatcher/models"
	"lucadomeneghetti/ipset_dispatcher/utils"
	"os"
	"os/signal"
	"time"
)

func terminateExecution(channels []chan models.Decision) {
	log.Println("received SIGINT, terminating")
	for _, ch := range channels {
		close(ch)
	}
	log.Println("closed all filter channels")
	os.Exit(0)
}

func main() {
	err := utils.ImportConfig("config.yml")
	if err != nil {
		log.Fatalln(err)
	}

	sigintc := make(chan os.Signal, 1)
	signal.Notify(sigintc, os.Interrupt)

	duration, err := time.ParseDuration(utils.GetUpdateTime())
	if err != nil {
		log.Fatalln(err)
	}
	ticker := time.NewTicker(duration)

	var startup bool = true
	channels := core.InitializeFilters("filters.yml")
	for {
		select {
		case <-sigintc:
			terminateExecution(channels)
		case <-ticker.C:
			decisions, err := core.QueryUpdateDecisions(startup, utils.GetRetries())
			if err != nil {
				log.Fatalln(err)
			}
			startup = false
			core.ForwardDecisions(decisions, channels)
		}
	}
}
