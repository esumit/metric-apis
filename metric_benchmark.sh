#!/bin/bash

METRIC_KEY="active_visitors"

RANDOM_GETPOST=600
REQUESTS=2000
CONCURRENCY=10

get_metric_sum_benchamrk(){
  echo
  echo "ab -n $REQUESTS -c $CONCURRENCY 'localhost:9000/$METRIC_KEY/sum'"
  echo "GET <<<"
  echo "URL: http://localhost:9000/metrics/'$METRIC_KEY'/sum"
  echo
  result=$(ab -n $REQUESTS -c $CONCURRENCY 'localhost:9000/'$METRIC_KEY'/sum')
  echo
  echo
  echo $result
}


get_metric_sum_benchamrk

