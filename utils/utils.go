package utils

import "fmt"

func BuildVersionString(version, revision, date string) string {
	result := fmt.Sprintf("Version: %s-%s (built at %s)\n", version, revision, date)
	return result
}
