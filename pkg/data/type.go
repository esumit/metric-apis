package data

import (
	"context"
	"time"
)

type MetricListData interface {
	SaveMetricData(ctx context.Context, data MetricData, seconds int) error
	GetMetricSumByKey(ctx context.Context, key string, seconds int) (*MetricRs ,error)
	Clears() (err error)
}

type MetricData struct {
	ID         string
	CreatedAt  time.Time
	Key        string
	Value      int16
}

 type MetricRs struct{
 	StartTime   time.Time
 	EndTime     time.Time
 	Sum int16
 	Count int16
 }
 