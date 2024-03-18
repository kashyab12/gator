package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
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
	appRouter.Use(legler.CorsMiddleware)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/readiness", legler.GetReadinessLegler)
	apiRouter.Get("/err", legler.GetErrorLegler)
	apiRouter.Post("/users", handlerConfig.PostUsersLegler)
	apiRouter.Get("/users", handlerConfig.AuthMiddleware(handlerConfig.GetUsersLegler))
	apiRouter.Post("/feeds", handlerConfig.AuthMiddleware(handlerConfig.PostFeedsLegler))
	apiRouter.Get("/feeds", handlerConfig.GetFeedsLegler)
	apiRouter.Post("/feed_follows", handlerConfig.AuthMiddleware(handlerConfig.PostFeedFollowLegler))
	apiRouter.Delete("/feed_follows/{feedFollowID}", handlerConfig.DeleteFeedFollow)
	appRouter.Mount("/v1", apiRouter)

	server := http.Server{
		Handler: appRouter,
		Addr:    fmt.Sprintf(":%v", serverPort),
	}
	err := server.ListenAndServe()
	if err != nil {
		return
	}
}
