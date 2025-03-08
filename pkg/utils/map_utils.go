package utils

import (
	"fmt"
	"maps"
)

func MapInterfaceInterfaceToStringInterface(input map[any]any) map[string]any {
	casted := make(map[string]any)
	for s, i := range input {
		casted[fmt.Sprint(s)] = i
	}

	return casted
}

func GetInMapValueOrDefault(keyToFind []string, myMap map[string]any, defaultValue any) any {
	current := myMap
	for i := 0; i < len(keyToFind); i++ {
		if _, found := current[keyToFind[i]]; !found {
			return defaultValue
		}

		if i+1 == len(keyToFind) {
			return current[keyToFind[i]]
		}

		current = current[keyToFind[i]].(map[string]any)
	}

	return current
}

func MapStringStructToSlice(input map[string]struct{}) []string {
	keys := make([]string, 0, len(input))
	for key := range maps.Keys(input) {
		keys = append(keys, key)
	}

	return keys
}
