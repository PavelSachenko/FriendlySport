#!/bin/bash
source .env

#curl --header "Content-Type: application/json" \
#  --request PUT \
#  --data '{"mappings": {"properties": {"title": {"type": "text", "fields": {"keyword": {"type": "keyword"}}}}}}' \
#  http://$ELASTIC_HOST$ELASTIC_PORT/friendly_sport_workout_recommendation

#for str in "arms" "legs" "body" "running" "yoga" "crossfit" "back with chest" "legs and back" \
#            "shoulder" "dance" "sprint" "back and shoulder" "arms with chest";\
# do
#   echo '{"title":"'$str'"}'
#  curl --header "Content-Type: application/json" \
#    --request POST \
#    --data '{"title":"'$str'"}' \
#    http://$ELASTIC_HOST$ELASTIC_PORT/friendly_sport_workout_recommendation/_doc
#done
#
#echo "\r\n"

  curl --header "Content-Type: application/json" \
    --request POST \
    --data '{"title":"test and test"}' \
    http://$ELASTIC_HOST$ELASTIC_PORT/friendly_sport_workout_recommendation/_doc