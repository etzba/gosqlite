#!/bin/bash
curl -X POST http://localhost:8080/signup -d '{ "username": "Etz", "password": "Pass1234" }'
COOKIE=$(curl -X POST http://localhost:8080/login -d '{ "username": "Etz", "password": "Pass1234" }' | jq -r '.[0]' )
printf "\nheaders\n"
curl http://localhost:8080/headers \
   -H "Accept: application/json" \
   -H "Authorization: Bearer $COOKIE"
printf "\nalpha\n"
curl http://localhost:8080/headers/alpha \
   -H "Accept: application/json" \
   -H "Authorization: Bearer $COOKIE" \
   --cookie "session=$COOKIE"
printf "\naahlp\n"
curl http://localhost:8080/headers/aahlp \
   -H "Accept: application/json" \
   -H "Authorization: Bearer $COOKIE" \
   --cookie "session=$COOKIE"