package utils

import (
	"fmt"
	"os"
)

func StringInSlice(list []string, q string) bool {
	for i := 0; i < len(list); i++ {
		if list[i] == q {
			return true
		}
	}
	return false
}

func ExitIfErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
