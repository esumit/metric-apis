package metric

import (
	"context"
)

const DATETIME_FORMAT = "2006-01-02 15:04:05"
type DataManager interface {
	Save(ctx context.Context, rq *MetricCreateRq) (*MetricCreateRs, error)
	Get(ctx context.Context, key string) (*MetricSumRs, error)
}

type MetricCreateRq struct {
	Key   string   `json:"key,omitempty"`
	Value int16    `json:"value,omitempty"`
}

type MetricCreateRs struct {
	ID    string   `json:"metric_id"`
}

type MetricSumRs struct {
	StartTime   string  `json:"start_time"`
	EndTime     string  `json:"end_time"`
	MetricCount int16     `json:"metric_count"`
	Key         string  `json:"metric_key"`
	Value       int16   `json:"metric_sum"`
}