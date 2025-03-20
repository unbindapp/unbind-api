package sourceanalyzer

import "github.com/unbindapp/unbind-api/internal/sourceanalyzer/enum"

// Assume default port based on framework defaults
func inferPortFromFramework(framework enum.Framework) *int {
	var port int
	switch framework {
	// Node frameworks
	case enum.Next, enum.CRA, enum.Remix, enum.Bun:
		port = 3000
	case enum.Astro:
		port = 4321
	case enum.Vite:
		port = 5173
	case enum.Angular:
		port = 4200

	// Python frameworks
	case enum.Django, enum.FastAPI, enum.FastHTML:
		port = 8000
	case enum.Flask:
		port = 5000

	// Go frameworks
	case enum.Gin:
		port = 8080

	// Java frameworks
	case enum.SpringBoot:
		port = 8080

	// PHP frameworks
	case enum.Laravel:
		port = 8000
	}

	if port == 0 {
		return nil
	}
	return &port
}
