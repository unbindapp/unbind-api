package prometheus

import "github.com/prometheus/common/model"

type NodeMetrics struct {
	CPU        []model.SamplePair
	RAM        []model.SamplePair
	Network    []model.SamplePair
	Disk       []model.SamplePair
	FileSystem []model.SamplePair
	Load       []model.SamplePair
}

// NodeMetricsFilterSumBy defines the possible grouping options for node metrics
type NodeMetricsFilterSumBy string

const (
	NodeSumByName    NodeMetricsFilterSumBy = "node"
	NodeSumByZone    NodeMetricsFilterSumBy = "zone"
	NodeSumByRegion  NodeMetricsFilterSumBy = "region"
	NodeSumByCluster NodeMetricsFilterSumBy = "cluster"
)

// Label returns the label name for use in Prometheus queries
func (s NodeMetricsFilterSumBy) Label() string {
	return string(s)
}

// NodeMetricsFilter contains filtering options for node metrics
type NodeMetricsFilter struct {
	Name    []string
	Zone    []string
	Region  []string
	Cluster []string
}
