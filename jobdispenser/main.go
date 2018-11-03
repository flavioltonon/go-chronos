package main

import (
	"flavioltonon/go-chronos/chronos"
	"log"
)

func main() {
	var chronos chronos.Chronos

	err := chronos.UpdateIssuesDeadlines()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("All deadlines have been updated successfully")
}
