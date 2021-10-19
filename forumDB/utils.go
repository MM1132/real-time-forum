package forumDB

import "log"

// Helper function to print errors
func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func fatalErr(err error) {
	if err != nil {
		panic(err)
	}
}
