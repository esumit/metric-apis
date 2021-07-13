package data

import (
	"container/list"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

var mutex sync.Mutex

type listMetric struct {
	sync.Mutex
	lm *list.List
}

func NewListMetric() *listMetric {
	return &listMetric{sync.Mutex{},list.New()}
}

func (l *listMetric) Clears() (err error) {
	if l.lm == nil {
		return
	} else {
		l.lm.Init()
	}
	log.Println("Clearing Metric List ")
	return nil
}

func (l *listMetric) SaveMetricData(ctx context.Context, data MetricData, seconds int) error {
	l.Lock()
	defer l.Unlock()
	l.lm.PushFront(data)
	var StartTime = time.Now().Add(time.Second * time.Duration(-seconds))
	var EndTime = time.Now()
	
	for e := l.lm.Front(); e != nil; e = e.Next() {
		v, ok := e.Value.(MetricData)
		
		if !ok {
			log.Println("Error")
		}
		if !(v.CreatedAt.Before(EndTime) && v.CreatedAt.After(StartTime)) {
			l.lm.Remove(e)
		}
	}
	
	return nil
}

func (l *listMetric) GetMetricSumByKey(ctx context.Context, key string, seconds int) (*MetricRs, error) {
	mutex.Lock()
	defer mutex.Unlock()
	
	var sum int16
	var count int16
	
	var StartTime = time.Now().Add(time.Second * time.Duration(-seconds))
	var EndTime = time.Now()

	for e := l.lm.Front(); e != nil; e = e.Next() {
		v, ok := e.Value.(MetricData)

		if !ok {
			fmt.Println("Error")
		}
		if v.Key == key && v.CreatedAt.Before(EndTime) && v.CreatedAt.After(StartTime) {
			sum = sum + v.Value
			count = count + 1
		} else {
			l.lm.Remove(e)
		}
	}
	
	var rs MetricRs
	
	rs.StartTime = StartTime
	rs.EndTime = EndTime
	rs.Sum = sum
	rs.Count = count
	return &rs, nil
}
