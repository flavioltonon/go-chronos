package main

import (
	"flavioltonon/go-chronos/github"
	"log"
	"net/http"
)

func main() {
	var chronos github.Chronos
	log.Println("server started")
	http.HandleFunc("/chronos", chronos.ListenToGitHub)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
