#!/bin/bash
source .env

curl --header "Content-Type: application/json" \
  --request PUT \
  --data '{"mappings": {"properties": {"title": {"type": "text", "fields": {"keyword": {"type": "keyword"}}}}}}' \
  http://0.0.0.0:9200/friendly_sport_workout_recommendation
