package sourceanalyzer

import (
	"encoding/json"

	core "github.com/railwayapp/railpack/core"
	a "github.com/railwayapp/railpack/core/app"
	"github.com/unbindapp/unbind-api/internal/common/log"
)

// The root result
type AnalysisResult struct {
	Provider  Provider  `json:"provider"` // Railpack provider (node, go, deno, python, java, etc.)
	Framework Framework `json:"framework"`
	Port      *int      `json:"port,omitempty"`
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

	detectedProvider := parseProvider(plan.DetectedProviders)
	detectedFramework := detectFramework(detectedProvider, plan)
	detectedPort := inferPortFromFramework(detectedFramework)

	// Railpack only returns gin, so see if this is a web API or not
	if detectedProvider == Go && detectedFramework == UnknownFramework {
		if hasListenAndServe(sourceDir) {
			port := 8080
			detectedPort = &port
		}
	}

	// Check for express as railpack doesn't return it
	if detectedProvider == Node && detectedFramework == UnknownFramework {
		if isExpressApp(sourceDir) {
			port := 3000
			detectedPort = &port
			detectedFramework = Express
		}
	}

	return &AnalysisResult{
		Provider:  detectedProvider,
		Framework: detectedFramework,
		Port:      detectedPort,
	}, nil
}
