package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// TODO return array instead
func loadEnvironment() (string, string) {
	err_load_env := godotenv.Load()
	if err_load_env != nil {
		log.Fatal("error load env")
	}
	dbURL := os.Getenv("DB_URL")

	return dbURL, os.Getenv("TOKEN_STRING")
}
