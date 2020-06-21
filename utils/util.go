package utils

import (
	"os"
)

// FileExists ensure if a file exists on determinated path
func FileExists(path string) (b bool, err error) {
	if _, err := os.Stat(path); err == nil {
		b = true
	} else if os.IsNotExist(err) {
		b = false
	} else {
		b = false
	}
	return
}
