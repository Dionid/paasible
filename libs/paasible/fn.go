package paasible

import (
	"fmt"
	"log"
)

func MergeVariables(
	storage *EntityStorage,
	variables ExtendableByVariables,
) (map[string]string, error) {
	resultVariables := make(map[string]string)

	for _, mapId := range variables.VariablesMapsIds {
		variableMap, ok := storage.VariablesMaps[mapId]
		if !ok {
			return nil, fmt.Errorf("Failed to find variable map with ID %s", mapId)
		}

		for variableKey, variable := range variableMap.VariablesIds {
			variableEnt, ok := storage.Variables[variable]
			if !ok {
				return nil, fmt.Errorf("Failed to find variable with ID %s", variableKey)
			}

			if resultVariables[variableKey] != "" {
				log.Printf("Variable '%s' already exists", variableKey)
			}
			resultVariables[variableKey] = variableEnt.Value
		}

		for variableKey, variable := range variableMap.Variables {
			if resultVariables[variableKey] != "" {
				log.Printf("Variable '%s' already exists", variableKey)
			}
			resultVariables[variableKey] = variable.Value
		}
	}

	for variableKey, variable := range variables.VariablesIds {
		variableEnt, ok := storage.Variables[variable]
		if !ok {
			return nil, fmt.Errorf("Failed to find variable with ID %s", variableKey)
		}

		if resultVariables[variableKey] != "" {
			log.Printf("Variable '%s' already exists", variableKey)
		}
		resultVariables[variableKey] = variableEnt.Value
	}

	for variableKey, variable := range variables.Variables {
		if resultVariables[variableKey] != "" {
			log.Printf("Variable '%s' already exists", variableKey)
		}
		resultVariables[variableKey] = variable.Value
	}

	return resultVariables, nil
}
