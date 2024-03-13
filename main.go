package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

const (
	DefaultPort = "8080"
)

func main() {
	if loadErr := godotenv.Load(); loadErr != nil {
		log.Fatalln(loadErr)
	}
	if serverPort := os.Getenv("PORT"); serverPort == "" {
		serverPort = DefaultPort
	}
	fmt.Println("gator gator")
}
