package models

type Decision struct {
	UUID      string `json:"uuid"`
	Scenario  string `json:"scenario"`
	IPAddress string `json:"ip"`
	Type      string `json:"type"`
	Duration  string `json:"duration"`
	Action    string
}

type DecisionArray []Decision
