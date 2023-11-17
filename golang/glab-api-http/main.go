package main

import (
	"log"
	"os"
)

func main() {
	token := os.Getenv("GITLAB_TOKEN")
	if token == "" {
		log.Fatal("Missing token")
	}
	fmt.Println(searchProjectID(251, token))
}
