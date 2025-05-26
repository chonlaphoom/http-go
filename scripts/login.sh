#!/bin/sh

curl -d "./json/login.json" -X POST http://localhost:8080/api/users -H "Content-Type: application/json"
echo "\n"
