package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"devfolio/backend/internal/data"
	"devfolio/backend/internal/domain"
	"devfolio/backend/internal/server"
	"devfolio/backend/internal/store"
)

func main() {
	seedData := data.Seed()
	repository := chooseRepository(seedData)
	apiServer := server.New(repository)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	httpServer := &http.Server{
		Addr:              ":" + port,
		Handler:           apiServer,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("devfolio-api listening on %s", httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func chooseRepository(seed domain.Data) store.Repository {
	mongoURI := os.Getenv("DEVFOLIO_MONGO_URI")
	if mongoURI == "" {
		log.Println("devfolio-api running with in-memory store (set DEVFOLIO_MONGO_URI for MongoDB persistence)")
		return store.New(seed)
	}

	databaseName := os.Getenv("DEVFOLIO_MONGO_DB")
	if databaseName == "" {
		databaseName = "devfolio"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoStore, err := store.NewMongo(ctx, mongoURI, seed, databaseName)
	if err != nil {
		log.Printf("failed to connect to MongoDB, falling back to in-memory store: %v", err)
		return store.New(seed)
	}

	log.Printf("devfolio-api using MongoDB database %q", databaseName)
	return mongoStore
}