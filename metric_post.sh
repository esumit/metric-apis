#!/bin/bash

METRIC_KEY="active_visitors"

RANDOM_GETPOST=600

post_metric_sum(){
  data="{ \"value\": $((RANDOM % 20))}"
  echo
  echo "POST >>>"
  echo  "URL: http://localhost:9000/metrics/$METRIC_KEY"
  echo
  echo "$data"
  result=$(curl -s -i --location --request POST 'http://localhost:9000/metrics/'$METRIC_KEY'' \
--header 'Content-Type: application/json' \
--data "$data" )

  echo
  echo "Response:"
  echo
  echo $result

}


for ((i = 1; i < RANDOM_GETPOST; ++i)); do
	echo
	echo "Metric Logging : $i "
	echo "Press [CTRL+C] to stop.."
  echo "-----------------------------------------------"
  post_metric_sum
  sleep 3

done
