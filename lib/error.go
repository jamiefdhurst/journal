package lib

import "log"

// CheckErr Check and fatal if applicable
func CheckErr(err error) {
	if err != nil {
		log.Fatal("Error reported: ", err)
	}
}
