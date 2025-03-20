package sourceanalyzer

import (
	"encoding/json"

	core "github.com/railwayapp/railpack/core"
	a "github.com/railwayapp/railpack/core/app"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/sourceanalyzer/enum"
)

// The root result
type AnalysisResult struct {
	Provider  enum.Provider  `json:"provider"` // Railpack provider (node, go, deno, python, java, etc.)
	Framework enum.Framework `json:"framework"`
	Port      *int           `json:"port,omitempty"`
}

func AnalyzeSourceCode(sourceDir string) (*AnalysisResult, error) {
	app, err := a.NewApp(sourceDir)
	if err != nil {
		return nil, err
	}

	// Generate railpack plan
	plan := core.GenerateBuildPlan(app, &a.Environment{}, &core.GenerateBuildPlanOptions{
		RailpackVersion: "unbind",
	})

	//Serialize to debug
	marshalled, _ := json.Marshal(plan)

	log.Infof("Plan data: %v", string(marshalled))

	detectedProvider := enum.ParseProvider(plan.DetectedProviders)
	detectedFramework := enum.DetectFramework(detectedProvider, plan)
	detectedPort := inferPortFromFramework(detectedFramework)

	// Railpack only returns gin, so see if this is a web API or not
	if detectedProvider == enum.Go && detectedFramework == enum.UnknownFramework {
		if hasListenAndServe(sourceDir) {
			port := 8080
			detectedPort = &port
		}
	}

	// Check for express as railpack doesn't return it
	if detectedProvider == enum.Node && detectedFramework == enum.UnknownFramework {
		if isExpressApp(sourceDir) {
			port := 3000
			detectedPort = &port
			detectedFramework = enum.Express
		}
	}

	return &AnalysisResult{
		Provider:  detectedProvider,
		Framework: detectedFramework,
		Port:      detectedPort,
	}, nil
}
