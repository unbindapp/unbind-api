package models

import (
	"math"
	"reflect"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/infrastructure/prometheus"
)

// Add this to your existing models package:
type NodeMetricsType string

const (
	NodeMetricsTypeNode    NodeMetricsType = "node"
	NodeMetricsTypeZone    NodeMetricsType = "zone"
	NodeMetricsTypeRegion  NodeMetricsType = "region"
	NodeMetricsTypeCluster NodeMetricsType = "cluster"
)

var NodeMetricsTypeValues = []NodeMetricsType{
	NodeMetricsTypeNode,
	NodeMetricsTypeZone,
	NodeMetricsTypeRegion,
	NodeMetricsTypeCluster,
}

// Register enum in OpenAPI specification
func (u NodeMetricsType) Schema(r huma.Registry) *huma.Schema {
	if r.Map()["NodeMetricsType"] == nil {
		schemaRef := r.Schema(reflect.TypeOf(""), true, "NodeMetricsType")
		schemaRef.Title = "NodeMetricsType"
		for _, v := range NodeMetricsTypeValues {
			schemaRef.Enum = append(schemaRef.Enum, string(v))
		}
		r.Map()["NodeMetricsType"] = schemaRef
	}
	return &huma.Schema{Ref: "#/components/schemas/NodeMetricsType"}
}

// NodeMetricsQueryInput defines the query parameters for node prometheus metrics
type NodeMetricsQueryInput struct {
	Type        NodeMetricsType `query:"type" required:"true"`
	NodeName    string          `query:"node_name" required:"false"`
	Zone        string          `query:"zone" required:"false"`
	Region      string          `query:"region" required:"false"`
	ClusterName string          `query:"cluster_name" required:"false"`
	Start       time.Time       `query:"start" required:"false" doc:"Start time for the query, defaults to 24 hours ago"`
	End         time.Time       `query:"end" required:"false" doc:"End time for the query, defaults to now"`
}

// NodeMetricsResult holds the transformed metrics data for nodes
type NodeMetricsResult struct {
	Labels  []string                    `json:"labels"`
	Step    time.Duration               `json:"step"`
	Metrics map[string]*NodeMetricsData `json:"metrics"`
}

// NodeMetricsData contains the time series data for a specific node metric
type NodeMetricsData struct {
	CPU        []float64 `json:"cpu"`        // CPU usage in cores
	RAM        []float64 `json:"ram"`        // RAM usage in bytes
	Network    []float64 `json:"network"`    // Network I/O in bytes/sec
	Disk       []float64 `json:"disk"`       // Disk I/O in bytes/sec
	FileSystem []float64 `json:"filesystem"` // Filesystem usage in bytes
	Load       []float64 `json:"load"`       // Load average (1min)
}

// TransformNodeMetricsEntity converts raw prometheus metrics to our API response format
func TransformNodeMetricsEntity(
	rawMetrics map[string]*prometheus.NodeMetrics,
	step time.Duration,
	sumBy prometheus.NodeMetricsFilterSumBy,
) *NodeMetricsResult {
	result := &NodeMetricsResult{
		Labels:  make([]string, 0),
		Step:    step,
		Metrics: make(map[string]*NodeMetricsData),
	}

	// Process time labels (use the first series if available)
	for _, metrics := range rawMetrics {
		if len(metrics.CPU) > 0 {
			for _, pair := range metrics.CPU {
				result.Labels = append(result.Labels, pair.Timestamp.Time().Format(time.RFC3339))
			}
			break
		}
	}

	// Process metrics
	for nodeID, metrics := range rawMetrics {
		nodeData := &NodeMetricsData{
			CPU:        make([]float64, len(result.Labels)),
			RAM:        make([]float64, len(result.Labels)),
			Network:    make([]float64, len(result.Labels)),
			Disk:       make([]float64, len(result.Labels)),
			FileSystem: make([]float64, len(result.Labels)),
			Load:       make([]float64, len(result.Labels)),
		}

		// Initialize with NaN to indicate missing data
		for i := range nodeData.CPU {
			nodeData.CPU[i] = math.NaN()
			nodeData.RAM[i] = math.NaN()
			nodeData.Network[i] = math.NaN()
			nodeData.Disk[i] = math.NaN()
			nodeData.FileSystem[i] = math.NaN()
			nodeData.Load[i] = math.NaN()
		}

		// Fill in CPU values
		for i, pair := range metrics.CPU {
			if i < len(nodeData.CPU) {
				nodeData.CPU[i] = float64(pair.Value)
			}
		}

		// Fill in RAM values
		for i, pair := range metrics.RAM {
			if i < len(nodeData.RAM) {
				nodeData.RAM[i] = float64(pair.Value)
			}
		}

		// Fill in Network values
		for i, pair := range metrics.Network {
			if i < len(nodeData.Network) {
				nodeData.Network[i] = float64(pair.Value)
			}
		}

		// Fill in Disk values
		for i, pair := range metrics.Disk {
			if i < len(nodeData.Disk) {
				nodeData.Disk[i] = float64(pair.Value)
			}
		}

		// Fill in FileSystem values
		for i, pair := range metrics.FileSystem {
			if i < len(nodeData.FileSystem) {
				nodeData.FileSystem[i] = float64(pair.Value)
			}
		}

		// Fill in Load values
		for i, pair := range metrics.Load {
			if i < len(nodeData.Load) {
				nodeData.Load[i] = float64(pair.Value)
			}
		}

		result.Metrics[nodeID] = nodeData
	}

	return result
}
