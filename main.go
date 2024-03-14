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
	// TODO: add chi and define separate routers and mount
	appRouter := http.NewServeMux()
	appRouter.HandleFunc("/v1/readiness", legler.GetReadinessLegler)
	appRouter.HandleFunc("/v1/err", legler.GetErrorLegler)
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
