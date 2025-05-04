package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func loadEnvironment() string {
	err_load_env := godotenv.Load()
	if err_load_env != nil {
		log.Fatal("error load env")
	}
	dbURL := os.Getenv("DB_URL")

	return dbURL
}
