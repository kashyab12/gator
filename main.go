package main

import (
	"database/sql"
	"fmt"
	chi "github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/kashyab12/gator/internal/database"
	"github.com/kashyab12/gator/legler"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func main() {
	if loadErr := godotenv.Load(); loadErr != nil {
		log.Fatalln(loadErr)
	}
	serverPort := os.Getenv("PORT")
	dbURL := os.Getenv("DB_CONN")
	db, openErr := sql.Open("postgres", dbURL)
	if openErr != nil {
		fmt.Println("Unable to open psql db connection")
	}
	dbQueries := database.New(db)
	handlerConfig := legler.ApiConfig{
		DB: dbQueries,
	}
	appRouter := chi.NewRouter()
	appRouter.Get("/v1/readiness", legler.GetReadinessLegler)
	appRouter.Get("/v1/err", legler.GetErrorLegler)
	appRouter.Post("/v1/users", handlerConfig.PostUsersLegler)
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
