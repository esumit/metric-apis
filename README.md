#### Metric APIs
Metric logging and reporting service that sums metrics by time window for 
the most recent collection time.  

#### Demo Recording

https://tinyurl.com/trace-fruits-demo 

#### How conceptually it works ?

This conceptual metric apis have implementation of two metric apis:

On post request, its get an integer metric value for a specified metric key e.g. active_visitors.
It gets the metric data, prepare a data model as metric data (CreatedAt, Key, Value), 
provide it an Id and store in metric list data structure. 

On every insert it checks the configured COLLECTION_TIMEOUT value. It traverse
all list's and check CreatedAt time should be in between current_time 
and (current_time - COLLECTION_TIMEOUT), means most recent COLLECTION_TIMEOUT
as time window. If list's elements not matches to time window then it removes it
from the list. 

on success, it returns the id of created metric-data in the list. 


````
Request: post metric data for key - active_visitors
POST localhost:9000/metrics/active_visitors
Response:
{
    "value": 20
}
````

![Screenshot Request: post metric data for key - active_visitors ](/images/Metric-POST-ByKey-2021-07-15.png)


On get request, its get sum of metric values for a specified metric key e.g. active_visitors.

On every get it checks the configured COLLECTION_TIMEOUT value. It traverse
all list's data and check its CreatedAt time should be in between current_time 
and (current_time - COLLECTION_TIMEOUT), means most recent COLLECTION_TIMEOUT
time window.

It picks the list elements which matches to the time window, it the calculate sum of its metric data Values, count number of metrics,
collection start time, collection end time and return as json.
If list's elements not matches to time window then it removes it from the list. 

It uses sync.mutex and apply lock on every insert/get and unlock it after.

````
Request : get metric data for key - active_visitors
GET localhost:9000/metrics/active-visitors/sum
Response :
{
    "start_time": "2021-07-14 23:44:10",
    "end_time": "2021-07-14 23:47:30",
    "metric_count": 6,
    "metric_key": "active_visitors",
    "metric_sum": 64
}

Refer start_time and end_time :
- Its configured based COLLECTION_TIMEOUT env variable 
- On this example its 200 seconds 

````
![Screenshot Request: get metric data for key - active_visitors ](/images/Metric-GET-ByKey-2021-07-15.png)




### How to

######  Step-1 : Build Docker  
````
make docker
````

######  step-2 : Refer env variables e.g. configure COLLECTION_TIMEOUT in seconds

````
- SERVER_PORT=9000
- SERVER_IP_ADDRESS=0.0.0.0
- HTTP_WRITE_TIMEOUT=15
- HTTP_READ_TIMEOUT=15
- HTTP_IDLE_TIMEOUT=60
- COLLECTION_TIMEOUT=200
````
 
######  step-3 : run docker-compose up  

It will run build image from step-1

######  step-4 - run metric_post.sh from a terminal 

Refer demo video 

######  step-5 - run metric_get.sh from a terminal 

Refer demo video

######  step-6 - run metric_benchmark.sh from a terminal 

Refer demo video

#### How To Test 

###### Time Window Test

Run: ./metric_post.sh 

Run: ./metric_get.sh 

Observe: metric_post.sh sending one post metric request in every 3 seconds, also 
metric_get.sh calling get metric request in every 3 seconds.

Observe:  metric_count value in metric_get.sh output :

- It increases based on the configured COLLECTION_TIMEOUT e.g. if its configured 
  as 30 seconds then metric_count value will be shown 10 as maximum 
 
Stop: metric_post.sh 

Observe: metric_count value in metric_get.sh output :
  
- It decreases in every 3 seconds and reduced to 0  

###### Concurrency test 

Run: ./metric_post.sh ( set each post request to 1 second or remove sleep)

Run: ./metric_benchmark.sh (It sends 2000 with concurrent count as 10 )

Observe: terminal logs 

- Everything works as expected

(Refer terminal logs)


#### Further work items (not limited to)

Its just a short demo, many further possibilities are with this :

###### 
* Separate http request/response library 
* Separate error related library 
* Various code optimisations 
* Containerize whole thing to docker for testing and implemenation
* Test cases 
* Validations on reuqest input, and wherever applicable by using e.g. ozzo  
* Revisit the logical flow and add cases
* Additional apis
* Code refactoring 
* Makefile with additional command e.g. docker push to docker
* Single script to run whole deployment 
* Additional error cases 
* OpenAPI specifications 
* Enhanced Logging 
* Swagger Integrations 
* Connected with Kibana etc for logging
* ... and so on 

###### go, docker, mac os - versions  
````
➜  metric-apis git:(main) ✗ go version
go version go1.16.5 darwin/amd64
➜  metric-apis git:(main) ✗ 

➜  metric-apis git:(main) ✗ docker --version
Docker version 19.03.8, build afacb8b
➜  metric-apis git:(main) ✗ 

➜  metric-apis git:(main) ✗ uname -a
Darwin Newyork.local 20.5.0 Darwin Kernel Version 20.5.0: Sat May  8 05:10:33 PDT 2021; root:xnu-7195.121.3~9/RELEASE_X86_64 x86_64
➜  metric-apis git:(main) ✗ 

````
###### Logs on terminal 

###### - Make
````
➜  metric-apis git:(main) ✗ make clean
rm -rf ./dist/*
➜  metric-apis git:(main) ✗ make build
rm -rf ./dist
mkdir dist
CGO_ENABLED=0 go build -o dist/metric-apis
cp .env dist/.env
➜  metric-apis git:(main) ✗ make run
./dist/metric-apis
INFO[0000] Config Applied:                              
INFO[0000] Port:  9000                                  
INFO[0000] IPAddress:  0.0.0.0                          
INFO[0000] HTTP WriteTimeout:  15                       
INFO[0000] HTTP ReadTimeout:  15                        
INFO[0000] HTTP IdleTimeout:  60                        
INFO[0000] Collection Timeout:  15                      
INFO[0000] All configs loaded  
````

````
➜  metric-apis git:(main) ✗ make docker
docker build -f Dockerfile -t esumit/metric-apis .
Sending build context to Docker daemon  9.963MB
Step 1/15 : FROM golang:1.12-alpine AS build_base
 ---> 76bddfb5e55e
Step 2/15 : RUN apk add --no-cache git
 ---> Using cache
 ---> 54d5d5637991
Step 3/15 : WORKDIR /app
 ---> Using cache
 ---> 604bff6f1803
Step 4/15 : COPY go.mod ./
 ---> Using cache
 ---> cc8d7841e086
Step 5/15 : COPY go.sum ./
 ---> Using cache
 ---> f4942387cb49
Step 6/15 : RUN go mod download
 ---> Using cache
 ---> 11d4051c8ed7
Step 7/15 : COPY . .
 ---> 03d469381a4c
Step 8/15 : RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o metric-apis .
 ---> Running in fe343c93ce5f
Removing intermediate container fe343c93ce5f
 ---> ced5c61a6bfb
Step 9/15 : FROM alpine:3.9
 ---> 78a2ce922f86
Step 10/15 : RUN apk add ca-certificates
 ---> Using cache
 ---> 6d0cc6a42b4b
Step 11/15 : WORKDIR /root/
 ---> Using cache
 ---> 256d200d3756
Step 12/15 : COPY --from=build_base /app/metric-apis .
 ---> Using cache
 ---> 5611a00f5c72
Step 13/15 : COPY --from=build_base /app/.env .
 ---> Using cache
 ---> 26e7cbbbf549
Step 14/15 : EXPOSE 9000
 ---> Using cache
 ---> 5e036e9ac923
Step 15/15 : ENTRYPOINT ["./metric-apis"]
 ---> Using cache
 ---> 7ddd4e95acdd
Successfully built 7ddd4e95acdd
Successfully tagged esumit/metric-apis:latest
➜  metric-apis git:(main) ✗ 
````

###### - docker

````
➜  metric-apis git:(main) ✗ docker run -p 9000:9000 esumit/metric-apis:latest
time="2021-07-15T00:03:34Z" level=info msg="Config Applied:"
time="2021-07-15T00:03:34Z" level=info msg="Port:  9000"
time="2021-07-15T00:03:34Z" level=info msg="IPAddress:  0.0.0.0"
time="2021-07-15T00:03:34Z" level=info msg="HTTP WriteTimeout:  15"
time="2021-07-15T00:03:34Z" level=info msg="HTTP ReadTimeout:  15"
time="2021-07-15T00:03:34Z" level=info msg="HTTP IdleTimeout:  60"
time="2021-07-15T00:03:34Z" level=info msg="Collection Timeout:  15"
time="2021-07-15T00:03:34Z" level=info msg="All configs loaded"
````

###### - scripts
````
➜  metric-apis git:(main) ✗ docker-compose up --remove-orphan
Removing orphan container "metric-apis_metrics-api_1"
Starting metrics-api.local.com ... done
Attaching to metrics-api.local.com
metrics-api.local.com | time="2021-07-15T00:07:54Z" level=info msg="Config Applied:"
metrics-api.local.com | time="2021-07-15T00:07:54Z" level=info msg="Port:  9000"
metrics-api.local.com | time="2021-07-15T00:07:54Z" level=info msg="IPAddress:  0.0.0.0"
metrics-api.local.com | time="2021-07-15T00:07:54Z" level=info msg="HTTP WriteTimeout:  15"
metrics-api.local.com | time="2021-07-15T00:07:54Z" level=info msg="HTTP ReadTimeout:  15"
metrics-api.local.com | time="2021-07-15T00:07:54Z" level=info msg="HTTP IdleTimeout:  60"
metrics-api.local.com | time="2021-07-15T00:07:54Z" level=info msg="Collection Timeout:  200"
metrics-api.local.com | time="2021-07-15T00:07:54Z" level=info msg="All configs loaded"
````

````
➜  metric-apis git:(main) ✗ ./metric_post.sh 

Metric Logging : 1 
Press [CTRL+C] to stop..
-----------------------------------------------

POST >>>
URL: http://localhost:9000/metrics/active_visitors

{ "value": 6}

Response:

 { "metric_id": "ad169b41-4845-46d2-a404-9229f1b6a0f5" }

Metric Logging : 2 
Press [CTRL+C] to stop..
-----------------------------------------------

POST >>>
URL: http://localhost:9000/metrics/active_visitors

{ "value": 19}

Response:

 { "metric_id": "e95d2668-75d2-44aa-bb30-8e450009c0a0" }
````

````
➜  metric-apis git:(main) ✗ ./metric_get.sh 

Metric Reporting : 1 
Press [CTRL+C] to stop..
-----------------------------------------------

GET <<<
URL: http://localhost:9000/metrics/active_visitors/sum


Response:

 { "start_time": "2021-07-15 00:05:49", "end_time": "2021-07-15 00:09:09", "metric_count": 10, "metric_key": "active_visitors", "metric_sum": 111 }

Metric Reporting : 2 
Press [CTRL+C] to stop..
-----------------------------------------------

GET <<<
URL: http://localhost:9000/metrics/active_visitors/sum


Response:

 { "start_time": "2021-07-15 00:05:52", "end_time": "2021-07-15 00:09:12", "metric_count": 11, "metric_key": "active_visitors", "metric_sum": 118 }

Metric Reporting : 3 
Press [CTRL+C] to stop..
-----------------------------------------------

GET <<<
URL: http://localhost:9000/metrics/active_visitors/sum


Response:

 { "start_time": "2021-07-15 00:05:55", "end_time": "2021-07-15 00:09:15", "metric_count": 12, "metric_key": "active_visitors", "metric_sum": 136 }

````

````
➜  metric-apis git:(main) ✗ ab -n 20000 -c 100 "127.0.0.1:9000/active_visitors/sum"
This is ApacheBench, Version 2.3 <$Revision: 1879490 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking 127.0.0.1 (be patient)
Completed 2000 requests
Completed 4000 requests
Completed 6000 requests
Completed 8000 requests
Completed 10000 requests
Completed 12000 requests
Completed 14000 requests
Completed 16000 requests
Completed 18000 requests
Completed 20000 requests
Finished 20000 requests


Server Software:        
Server Hostname:        127.0.0.1
Server Port:            9000

Document Path:          /active_visitors/sum
Document Length:        51 bytes

Concurrency Level:      100
Time taken for tests:   38.532 seconds
Complete requests:      20000
Failed requests:        0
Non-2xx responses:      20000
Total transferred:      4540000 bytes
HTML transferred:       1020000 bytes
Requests per second:    519.05 [#/sec] (mean)
Time per request:       192.661 [ms] (mean)
Time per request:       1.927 [ms] (mean, across all concurrent requests)
Transfer rate:          115.06 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0  102 1288.2      1   20116
Processing:     2   90  64.0     88    1016
Waiting:        2   88  62.7     86    1016
Total:          3  192 1287.3     91   20238

Percentage of the requests served within a certain time (ms)
  50%     91
  66%    105
  75%    114
  80%    120
  90%    136
  95%    154
  98%    191
  99%    964
 100%  20238 (longest request)
➜  metric-apis git:(main) ✗ 
````

````
➜  metric-apis git:(main) ✗ ./metric_benchmark.sh 

ab -n 2000 -c 10 'localhost:9000/active_visitors/sum'
GET <<<
URL: http://localhost:9000/metrics/'active_visitors'/sum

Completed 200 requests
Completed 400 requests
Completed 600 requests
Completed 800 requests
Completed 1000 requests
Completed 1200 requests
Completed 1400 requests
Completed 1600 requests
Completed 1800 requests
Completed 2000 requests
Finished 2000 requests


This is ApacheBench, Version 2.3 <$Revision: 1879490 $> Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/ Licensed to The Apache Software Foundation, http://www.apache.org/ Benchmarking localhost (be patient) Server Software: Server Hostname: localhost Server Port: 9000 Document Path: /active_visitors/sum Document Length: 51 bytes Concurrency Level: 10 Time taken for tests: 3.329 seconds Complete requests: 2000 Failed requests: 0 Non-2xx responses: 2000 Total transferred: 454000 bytes HTML transferred: 102000 bytes Requests per second: 600.81 [#/sec] (mean) Time per request: 16.644 [ms] (mean) Time per request: 1.664 [ms] (mean, across all concurrent requests) Transfer rate: 133.19 [Kbytes/sec] received Connection Times (ms) min mean[+/-sd] median max Connect: 0 1 3.1 0 139 Processing: 3 16 13.5 14 163 Waiting: 3 15 12.9 14 163 Total: 3 16 13.8 15 164 Percentage of the requests served within a certain time (ms) 50% 15 66% 17 75% 19 80% 20 90% 22 95% 25 98% 36 99% 95 100% 164 (longest request)
➜  metric-apis git:(main) ✗ 
````