package utils

import "log"

func FatalErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
