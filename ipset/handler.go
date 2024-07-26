package ipset

import (
	"log"
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
		log.Printf("ipset %v already exists", name)
		return &IpSet{name: name}
	}
	err := ipset.Create(name, set_type, ipset.CreateOptions{})
	if err != nil {
		log.Fatalf("failed to create ipset %v: %v", name, err)
	}
	log.Printf("successfully created ipset %v", name)
	return &IpSet{name: name}
}

func (set *IpSet) AddEntry(entry string, duration uint32) error {
	err := ipset.Add(set.name, &ipset.Entry{IP: net.ParseIP(entry), Timeout: &duration})
	if err != nil {
		log.Fatalf("failed to add entry %v to ipset %v", entry, set.name)
	}
	log.Printf("successfully added entry %v to ipset %v", entry, set.name)
	return nil
}

func (set *IpSet) DeleteEntry(entry string) error {
	err := ipset.Del(set.name, &ipset.Entry{IP: net.ParseIP(entry)})
	if err != nil {
		log.Fatalf("failed to delete entry %v from ipset %v", entry, set.name)
	}
	log.Printf("successfully deleted entry %v from ipset %v", entry, set.name)
	return nil
}
