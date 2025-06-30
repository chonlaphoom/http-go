#!/bin/sh

curl -d @login-2.json -X POST http://localhost:8080/api/login -H "Content-Type: application/json"
echo "\n"
