package enum

import (
	"reflect"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/railwayapp/railpack/core/generate"
)

// Detect framework based on provider and plan
func DetectFramework(provider Provider, ctx *generate.GenerateContext) Framework {
	switch provider {
	// Node frameworks
	case Node:
		framework, ok := ctx.Metadata.Properties["nodeRuntime"]
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
		framework, ok := ctx.Metadata.Properties["pythonRuntime"]
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
		gin, ok := ctx.Metadata.Properties["goGin"]
		if ok && gin == "true" {
			return Gin
		}
	// Java frameworks
	case Java:
		javaFramework, ok := ctx.Metadata.Properties["javaFramework"]
		if ok {
			switch javaFramework {
			case "spring-boot":
				return SpringBoot
			}
		}
	// PHP frameworks
	case PHP:
		for _, log := range ctx.Logger.Logs {
			if strings.Contains(strings.ToLower(log.Msg), "laravel") {
				return Laravel
			}
		}
	// Ruby frameworks
	case Ruby:
		railsFramework, ok := ctx.Metadata.Properties["rubyRails"]
		if ok && railsFramework == "true" {
			return Rails
		}
		// Rust frameworks
	case Rust:
		if ctx.Deploy != nil {
			for _, variable := range ctx.Deploy.Variables {
				if strings.EqualFold(variable, "rocket_address") {
					return Rocket
				}
			}
		}
	}

	return UnknownFramework
}

// Framework represents the detected framework
type Framework string

const (
	// * node frameworks
	Next          Framework = "next"
	Astro         Framework = "astro"
	Vite          Framework = "vite"
	CRA           Framework = "cra"
	Angular       Framework = "angular"
	Remix         Framework = "remix"
	Bun           Framework = "bun"
	Express       Framework = "express"
	Sveltekit     Framework = "sveltekit"
	Svelte        Framework = "svelte"
	Solid         Framework = "solid"
	Hono          Framework = "hono"
	TanstackStart Framework = "tanstack-start"
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
	// * Ruby frameworks
	Rails Framework = "rails"
	// * Rust frameworks
	Rocket Framework = "rocket"
	// * not detected
	UnknownFramework Framework = "unknown"
)

var allFrameworks = []Framework{
	Next, Astro, Vite, CRA, Angular, Remix, Bun, Express, Sveltekit, Svelte, Solid, Hono, TanstackStart,
	PythonFramework, Django, Flask, FastAPI, FastHTML,
	Gin,
	SpringBoot,
	Laravel,
	Rails,
	Rocket,
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
