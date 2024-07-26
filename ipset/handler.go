package ipset

import (
	"log"
	"log/slog"
	"net"

	"github.com/lrh3321/ipset-go"
)

type IpSet struct {
	name string
}

func checkSetExists(name string) bool {
	_, err := ipset.List(name)
	if err != nil {
		return false
	} else {
		return true
	}
}

func CreateSet(name string, set_type string) (set *IpSet) {
	if checkSetExists(name) {
		log.Printf("ipset %v already exists, overwriting", name)
		ipset.Flush(name)
		return &IpSet{name: name}
	}
	err := ipset.Create(name, set_type, ipset.CreateOptions{Timeout: 300})
	if err != nil {
		log.Fatalf("failed to create ipset %v: %v", name, err)
	}
	log.Printf("successfully created ipset %v", name)
	return &IpSet{name: name}
}

func (set *IpSet) AddEntry(entry string, duration uint32) {
	err := ipset.Add(set.name, &ipset.Entry{IP: net.ParseIP(entry), Timeout: &duration})
	if err != nil {
		slog.Debug("failed to add entry " + entry + " to ipset " + set.name + ": " + err.Error())
		return
	}
	log.Printf("successfully added entry %v to ipset %v", entry, set.name)
}

func (set *IpSet) DeleteEntry(entry string) {
	err := ipset.Del(set.name, &ipset.Entry{IP: net.ParseIP(entry)})
	if err != nil {
		slog.Debug("failed to delete entry " + entry + " from ipset " + set.name + ": " + err.Error())
		return
	}
	log.Printf("successfully deleted entry %v from ipset %v", entry, set.name)
}
