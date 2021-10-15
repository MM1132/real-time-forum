package utils

import "log"

func FatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CheckErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
