package sourceanalyzer

import (
	"encoding/json"
	"strings"

	core "github.com/railwayapp/railpack/core"
)

// Detect framework based on provider and plan
func detectFramework(provider Provider, plan *core.BuildResult) Framework {
	switch provider {
	// Node frameworks
	case Node:
		framework, ok := plan.Metadata["nodeRuntime"]
		if ok {
			switch framework {
			case "next":
				return Next
			case "astro":
				return Astro
			case "vite":
				return Vite
			case "cra":
				return CRA
			case "angular":
				return Angular
			case "remix":
				return Remix
			case "bun":
				return Bun
			case "express":
				return Express
			}
		}
	// Python frameworks
	case Python:
		framework, ok := plan.Metadata["pythonRuntime"]
		if ok {
			switch framework {
			case "python":
				return PythonFramework
			case "django":
				return Django
			case "flask":
				return Flask
			case "fastapi":
				return FastAPI
			case "fasthtml":
				return FastHTML
			}
		}
	// Go frameworks
	case Go:
		gin, ok := plan.Metadata["goGin"]
		if ok && gin == "true" {
			return Gin
		}
	// Java frameworks
	case Java:
		javaFramework, ok := plan.Metadata["javaFramework"]
		if ok {
			switch javaFramework {
			case "spring-boot":
				return SpringBoot
			}
		}
	// PHP frameworks
	case PHP:
		for _, log := range plan.Logs {
			if strings.Contains(strings.ToLower(log.Msg), "laravel") {
				return Laravel
			}
		}
	}

	return UnknownFramework
}

// Framework represents the detected framework
type Framework int

const (
	// * node frameworks
	Next Framework = iota
	Astro
	Vite
	CRA
	Angular
	Remix
	Bun
	Express
	// * python frameworks
	PythonFramework // Sometimes detected as a framework
	Django
	Flask
	FastAPI
	FastHTML
	// * GO frameworks
	Gin
	// * Java frameworks
	SpringBoot
	// * PHP frameworks
	Laravel
	// * not detected
	UnknownFramework
)

func (f Framework) String() string {
	names := map[Framework]string{
		// * node
		Next:    "next",
		Astro:   "astro",
		Vite:    "vite",
		CRA:     "cra",
		Angular: "angular",
		Remix:   "remix",
		Bun:     "bun",
		Express: "express",
		// * python
		PythonFramework: "python",
		Django:          "django",
		Flask:           "flask",
		FastAPI:         "fastapi",
		FastHTML:        "fasthtml",
		// * go
		Gin: "gin",
		// * java
		SpringBoot: "spring-boot",
		// * php
		Laravel:          "laravel",
		UnknownFramework: "unknown",
	}

	if name, ok := names[f]; ok {
		return name
	}
	return "unknown"
}

func (f Framework) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.String())
}
