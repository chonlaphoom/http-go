package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// TODO return array instead
func loadEnvironment() (string, string, string) {
	err_load_env := godotenv.Load()
	if err_load_env != nil {
		log.Fatal("error load env")
	}
	dbURL := os.Getenv("DB_URL")
	polka_key := os.Getenv("POLKA_KEY")

	return dbURL, os.Getenv("TOKEN_STRING"), polka_key
}
