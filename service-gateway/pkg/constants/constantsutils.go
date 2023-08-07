package constants

import (
	"fmt"
	"os"
)

func genericParse [T any](str string) (T, error) {
	var result T
	_, err := fmt.Sscanf(str, "%v", &result)
	return result, err
}

func getGenericEnv[K any](key string, defaultVal K) K {
	if value, exists := os.LookupEnv(key); exists {
		parsed, err := genericParse[K](value)
		if err == nil {
			return parsed
		}
	}

	return defaultVal
}