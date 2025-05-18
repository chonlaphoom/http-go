#!/bin/sh

# todo update test 
curl -d "./json/prohibit.json" -X POST http://localhost:8080/api/validate_chirp -H "Content-Type: application/json"
echo "\n"

curl -d "./json/long-body.json" -X POST http://localhost:8080/api/validate_chirp -H "Content-Type: application/json"
echo "\n"

curl -d "./json/simple.json" -X POST http://localhost:8080/api/validate_chirp -H "Content-Type: application/json"
echo "\n"
