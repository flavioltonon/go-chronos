package main

import (
	"flavioltonon/go-chronos/chronos"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	var chronos chronos.Chronos

	log.Println("Chronos is observing...")

	http.HandleFunc("/chronos", chronos.ListenToGitHub)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("NGROK_STANDARD_PORT")), nil))
}
