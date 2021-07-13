package metric

import (
	"context"
	"fmt"
	"github.com/esumit/metric-apis/pkg/data"
	"github.com/google/uuid"
	"time"
)

type metricApiDataManager struct {
	metricList data.MetricListData
	seconds int
}

func NewMetricApiDataManager(ml data.MetricListData, seconds int) *metricApiDataManager {
	return &metricApiDataManager{ml,seconds}
}

func (m *metricApiDataManager) Save(ctx context.Context, rq *MetricCreateRq) (*MetricCreateRs, error) {
	var md data.MetricData
	md.CreatedAt = time.Now()
	md.ID = uuid.New().String()
	md.Value = rq.Value
	md.Key = rq.Key

	err := m.metricList.SaveMetricData(ctx, md, m.seconds)

	if err != nil {
		fmt.Println("Error")
	}

	var rs MetricCreateRs

	rs.ID = md.ID
	return &rs, nil
}

func (m *metricApiDataManager) Get(ctx context.Context, key string) (*MetricSumRs, error) {

	sum, err := m.metricList.GetMetricSumByKey(ctx, key, m.seconds)

	if err != nil {
		fmt.Println("Error")
	}

	var rs MetricSumRs

	rs.StartTime = sum.StartTime.Format(DATETIME_FORMAT)
	rs.EndTime = sum.EndTime.Format(DATETIME_FORMAT)
	rs.Key = key
	rs.Value = sum.Sum
	rs.MetricCount = sum.Count

	return &rs, nil
}
