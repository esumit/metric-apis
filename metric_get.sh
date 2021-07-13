#!/bin/bash

METRIC_KEY="active_visitors"

RANDOM_GETPOST=600

get_metric_sum(){
  echo
  echo "GET <<<"
  echo "URL: http://localhost:9000/metrics/$METRIC_KEY/sum"
  echo
  result=$(curl -s -i --location --request GET --header 'Content-Type: application/json' 'http://localhost:9000/metrics/'$METRIC_KEY'/sum')
  echo
  echo "Response:"
  echo
  echo $result
}

for ((i = 1; i < RANDOM_GETPOST; ++i)); do
	echo
	echo "Metric Reporting : $i "
	echo "Press [CTRL+C] to stop.."
  echo "-----------------------------------------------"
  get_metric_sum
  sleep 3

done
