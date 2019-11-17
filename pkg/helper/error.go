package helper

import "fmt"

// CheckErr print a friendly error message
func CheckErr(err error) {
	switch {
	case err == nil:
		return
	default:
		fmt.Println(err)
	}
}
