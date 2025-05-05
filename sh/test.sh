#!/bin/sh

# todo validate response
curl -d "@./mock/replace.json" -X POST http://localhost:8080/api/validate_chirp -H "Content-Type: application/json"
echo "\n"

curl -d "@./mock/long.json" -X POST http://localhost:8080/api/validate_chirp -H "Content-Type: application/json"
echo "\n"

curl -d "@./mock/simple.json" -X POST http://localhost:8080/api/validate_chirp -H "Content-Type: application/json"
echo "\n"
