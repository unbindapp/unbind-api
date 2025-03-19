package sourceanalyzer

// Assume default port based on framework defaults
func inferPortFromFramework(framework Framework) *int {
	var port int
	switch framework {
	// Node frameworks
	case Next, CRA, Remix, Bun:
		port = 3000
	case Astro:
		port = 4321
	case Vite:
		port = 5173
	case Angular:
		port = 4200

	// Python frameworks
	case Django, FastAPI, FastHTML:
		port = 8000
	case Flask:
		port = 5000

	// Go frameworks
	case Gin:
		port = 8080

	// Java frameworks
	case SpringBoot:
		port = 8080

	// PHP frameworks
	case Laravel:
		port = 8000
	}

	if port == 0 {
		return nil
	}
	return &port
}
