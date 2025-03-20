package enum

import (
	"reflect"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	core "github.com/railwayapp/railpack/core"
)

// Detect framework based on provider and plan
func DetectFramework(provider Provider, plan *core.BuildResult) Framework {
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
type Framework string

const (
	// * node frameworks
	Next    Framework = "next"
	Astro   Framework = "astro"
	Vite    Framework = "vite"
	CRA     Framework = "cra"
	Angular Framework = "angular"
	Remix   Framework = "remix"
	Bun     Framework = "bun"
	Express Framework = "express"
	// * python frameworks
	PythonFramework Framework = "python" // Sometimes detected as a framework
	Django          Framework = "django"
	Flask           Framework = "flask"
	FastAPI         Framework = "fastapi"
	FastHTML        Framework = "fasthtml"
	// * GO frameworks
	Gin Framework = "gin"
	// * Java frameworks
	SpringBoot Framework = "spring-boot"
	// * PHP frameworks
	Laravel Framework = "laravel"
	// * not detected
	UnknownFramework Framework = "unknown"
)

var allFrameworks = []Framework{
	Next, Astro, Vite, CRA, Angular, Remix, Bun, Express,
	PythonFramework, Django, Flask, FastAPI, FastHTML,
	Gin,
	SpringBoot,
	Laravel,
	UnknownFramework,
}

// Values provides list valid values for Enum.
func (Framework) Values() (kinds []string) {
	for _, s := range allFrameworks {
		kinds = append(kinds, string(s))
	}
	return
}

// Register enum in OpenAPI specification
// https://github.com/danielgtaylor/huma/issues/621
func (u Framework) Schema(r huma.Registry) *huma.Schema {
	if r.Map()["Framework"] == nil {
		schemaRef := r.Schema(reflect.TypeOf(""), true, "Framework")
		schemaRef.Title = "Framework"
		for _, v := range allFrameworks {
			schemaRef.Enum = append(schemaRef.Enum, string(v))
		}
		r.Map()["Framework"] = schemaRef
	}
	return &huma.Schema{Ref: "#/components/schemas/Framework"}
}
