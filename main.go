package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kashyab12/gator/legler"
	"log"
	"net/http"
	"os"
)

func main() {
	if loadErr := godotenv.Load(); loadErr != nil {
		log.Fatalln(loadErr)
	}
	serverPort := os.Getenv("PORT")
	appRouter := http.NewServeMux()
	// Add handlers for route
	corsMux := legler.CorsMiddleware(appRouter)
	server := http.Server{
		Handler: corsMux,
		Addr:    fmt.Sprintf(":%v", serverPort),
	}
	err := server.ListenAndServe()
	if err != nil {
		return
	}
}
