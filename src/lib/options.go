package lib

import (
	"os"
	"slices"
)

// Extracts `option` from os.Args.
// Assumes that option requires a value and return it if exists.
func GetOption(option string) (string, bool) {
	optionIndex := slices.Index(os.Args, option)
	if optionIndex != -1 {
		if optionIndex+1 >= len(os.Args) {
			return "", true
		}
		optionValue := os.Args[optionIndex+1]
		os.Args = slices.Delete(os.Args, optionIndex, optionIndex+2)
		return optionValue, true
	}
	return "", false
}
