package main

import (
	"flavioltonon/go-chronos/chronos"
	"log"
	"net/http"
)

func main() {
	var chronos chronos.Chronos
	log.Println("Chronos is observing...")
	http.HandleFunc("/chronos", chronos.ListenToGitHub)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
