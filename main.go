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
	apiRouter := chi.NewRouter()
	apiRouter.Get("/readiness", legler.GetReadinessLegler)
	apiRouter.Get("/err", legler.GetErrorLegler)
	apiRouter.Post("/users", handlerConfig.PostUsersLegler)
	apiRouter.Get("/users", handlerConfig.GetUsersLegler)

	appRouter := chi.NewRouter()
	appRouter.Mount("/v1", apiRouter)

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
