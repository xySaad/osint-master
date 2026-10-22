package lib

import "fmt"

func PrintIfSet(label, value string) {
	if value != "" {
		fmt.Println(label, value)
	}
}
