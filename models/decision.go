package models

import (
)

type Decision struct {
	UUID string `json:"uuid"`
	Scenario  string `json:"scenario"`
	IPAddress string `json:"ip"`
	Type string `json:"type"`
	Until string `json:"until"`
	Duration string `json:"duration"`
}

type DecisionArray []Decision
