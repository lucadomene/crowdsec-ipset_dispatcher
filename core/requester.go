package core

import (
	"encoding/json"
	"fmt"
	"lucadomeneghetti/ipset_dispatcher/models"
	"lucadomeneghetti/ipset_dispatcher/utils"
	"net/http"
)

func QueryUpdateDecisions(startup bool, retry int) (models.DecisionArray, []string, error) {

	base_url := utils.GetBaseURL()
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/decisions/stream?startup=%v", base_url, startup), nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", utils.GetAPI())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	} else if res.StatusCode > 300 && retry > 0 {
		return QueryUpdateDecisions(startup, retry-1)
	} else if retry <= 0 {
		http_err := fmt.Errorf("%v", res.Status)
		return nil, nil, http_err
	}
	defer res.Body.Close()

	result := make(map[string]interface{})
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, nil, err
	}

	new_decisions_raw := result["new"].([]interface{})
	var new_decisions models.DecisionArray
	for _, v := range new_decisions_raw {
		var decision models.Decision

		v_parsed := v.(map[string]interface{})

		uuid, ok := v_parsed["uuid"].(string)
		if ok {
			decision.UUID = uuid
		}

		scenario, ok := v_parsed["scenario"].(string)
		if ok {
			decision.Scenario = scenario
		}

		ipaddr, ok := v_parsed["value"].(string)
		if ok {
			decision.IPAddress = ipaddr
		}

		dec_type, ok := v_parsed["type"].(string)
		if ok {
			decision.Type = dec_type
		}

		new_decisions = append(new_decisions, decision)
	}

	deleted_decisions_raw := result["deleted"].([]interface{})
	var deleted_decisions []string
	for _, w := range deleted_decisions_raw {
		var decisionUUID string

		w_parsed := w.(map[string]interface{})

		uuid, ok := w_parsed["uuid"].(string)
		if ok {
			decisionUUID = uuid
		}

		deleted_decisions = append(deleted_decisions, decisionUUID)
	}

	return new_decisions, deleted_decisions, nil
}
