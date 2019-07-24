package nexusrm

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const restMaintenanceDBCheck = "service/rest/v1/maintenance/%s/check"

func CheckDatabase(rm RM, dbName string) (DatabaseState, error) {
	doError := func(err error) error {
		return fmt.Errorf("error checking status of database '%s': %v", dbName, err)
	}

	var state DatabaseState

	url := fmt.Sprintf(restMaintenanceDBCheck, dbName)
	body, resp, err := rm.Put(url, nil)
	if err != nil || resp.StatusCode != http.StatusOK {
		return state, doError(err)
	}

	if err := json.Unmarshal(body, &state); err != nil {
		return state, doError(err)
	}

	return state, nil
}

func CheckAllDatabases(rm RM) (states map[string]DatabaseState, err error) {
	states = make(map[string]DatabaseState)

	check := func(dbName string) {
		if err != nil {
			return
		}

		if state, er := CheckDatabase(rm, dbName); er != nil {
			err = fmt.Errorf("error with '%s' database when all states: %v", dbName, er)
		} else {
			states[dbName] = state
		}
	}

	check(AccessLogDB)
	check(ComponentDB)
	check(ConfigDB)
	check(SecurityDB)

	return
}
