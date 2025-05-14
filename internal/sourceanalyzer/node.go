package sourceanalyzer

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	"github.com/unbindapp/unbind-api/internal/sourceanalyzer/enum"
)

// DetectNodeFramework inspects the project at sourceDir and returns the detected
// framework in priority order, falling back to Express and finally UnknownFramework.
//
// Detection now uses **layered heuristics**:
//  1. Primary dependencies (package.json dependencies or devDependencies)
//  2. Scripts section keywords (e.g. "svelte-kit", "vite", "solid-start")
//  3. Presence of canonical config files (e.g. svelte.config.js, vite.config.ts)
//  4. Fallback source‑code regexes (Express + Hono)
//
// Order is important – the first match wins.
func DetectNodeFramework(sourceDir string) enum.Framework {
	switch {
	case isSvelteKitApp(sourceDir):
		return enum.Sveltekit
	case isSvelteApp(sourceDir):
		return enum.Svelte
	case isSolidApp(sourceDir):
		return enum.Solid
	case isHonoApp(sourceDir):
		return enum.Hono
	case isTanStackStartApp(sourceDir):
		return enum.TanstackStart
	case isViteApp(sourceDir):
		return enum.Vite
	case isExpressApp(sourceDir):
		return enum.Express
	default:
		return enum.UnknownFramework
	}
}

// ---------------------------------------------------------------------------
// Individual framework detectors
// ---------------------------------------------------------------------------

func isTanStackStartApp(sourceDir string) bool {
	// Primary indicator
	if hasDependency(sourceDir, "@tanstack/react-start") && hasScriptKeyword(sourceDir, "vinxi") {
		return true
	}
	return false
}

func isViteApp(sourceDir string) bool {
	if hasDependency(sourceDir, "vite") || hasScriptKeyword(sourceDir, "vite") {
		return true
	}
	return hasAnyFile(sourceDir, []string{"vite.config.js", "vite.config.ts", "vite.config.mjs", "vite.config.cjs"})
}

func isSvelteKitApp(sourceDir string) bool {
	if hasDependency(sourceDir, "@sveltejs/kit") || hasScriptKeyword(sourceDir, "svelte-kit") {
		return true
	}
	// SvelteKit always has a svelte.config.* file in root
	return hasAnyFile(sourceDir, []string{"svelte.config.js", "svelte.config.ts", "svelte.config.cjs", "svelte.config.mjs"})
}

func isSvelteApp(sourceDir string) bool {
	// Must contain svelte, but we already checked for kit above
	if hasDependency(sourceDir, "svelte") || hasScriptKeyword(sourceDir, "svelte") {
		return true
	}
	// Svelte projects may also ship a svelte.config.*
	return hasAnyFile(sourceDir, []string{"svelte.config.js", "svelte.config.ts", "svelte.config.cjs", "svelte.config.mjs"})
}

func isSolidApp(sourceDir string) bool {
	if hasDependency(sourceDir, "solid-js") || hasDependency(sourceDir, "solid-start") {
		return true
	}
	return hasScriptKeyword(sourceDir, "solid-start") || hasAnyFile(sourceDir, []string{"solid.config.ts", "solid.config.js"})
}

func isHonoApp(sourceDir string) bool {
	if hasDependency(sourceDir, "hono") {
		return true
	}
	// Fast source scan: look for `new Hono()` in TypeScript / JavaScript files
	tsFiles, _ := utils.FindFilesWithExclusions(sourceDir, "*.{ts,js,tsx,jsx}", []string{"node_modules", ".git"})
	honoInit := regexp.MustCompile(`(?m)new\s+Hono\s*\(`)
	for _, f := range tsFiles {
		if content, err := utils.ReadFile(f); err == nil && honoInit.Match([]byte(content)) {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// Express detector (existing logic kept, now internal helper)
// ---------------------------------------------------------------------------

// isExpressApp checks if the Node.js app uses Express.js
func isExpressApp(sourceDir string) bool {
	// Check for express in package.json dependencies
	packageJSONPath := filepath.Join(sourceDir, "package.json")
	if utils.FileExists(packageJSONPath) {
		content, err := utils.ReadFile(packageJSONPath)
		if err == nil {
			// Check if express is listed as a dependency
			if strings.Contains(content, "\"express\":") || strings.Contains(content, "'express':") {
				// Look for express app initialization pattern in JavaScript files
				jsFiles, err := utils.FindFilesWithExclusions(sourceDir, "*.js", []string{"node_modules", ".git"})
				if err == nil {
					// Patterns to look for in source code
					expressInitRegex := regexp.MustCompile(`(?m)(const|let|var)\s+(\w+)\s*=\s*require\s*\(\s*['"]express['"]\s*\)`)
					expressListenRegex := regexp.MustCompile(`(?m)\.listen\s*\(\s*(\d+|process\.env\.[A-Z_]+)`)

					for _, file := range jsFiles {
						content, err := utils.ReadFile(file)
						if err == nil {
							// Check if file imports express and has listen method
							if expressInitRegex.MatchString(content) && expressListenRegex.MatchString(content) {
								log.Infof("Found Express app initialization in %s", file)
								return true
							}
						}
					}
				}

				// If we found express in dependencies but couldn't confirm with code pattern,
				// still return true since it's likely an Express app
				return true
			}
		}
	}

	return false
}

// ---------------------------------------------------------------------------
// Shared helpers
// ---------------------------------------------------------------------------

// hasDependency opens package.json and checks for the specified dependency string.
func hasDependency(sourceDir, dep string) bool {
	packageJSONPath := filepath.Join(sourceDir, "package.json")
	if !utils.FileExists(packageJSONPath) {
		return false
	}

	content, err := utils.ReadFile(packageJSONPath)
	if err != nil {
		return false
	}
	// naïve but fast scan of dependencies/devDependencies
	return strings.Contains(content, "\""+dep+"\":") || strings.Contains(content, "'"+dep+"':")
}

// hasScriptKeyword scans the package.json "scripts" section for the given keyword.
func hasScriptKeyword(sourceDir, kw string) bool {
	packageJSONPath := filepath.Join(sourceDir, "package.json")
	if !utils.FileExists(packageJSONPath) {
		return false
	}

	content, err := utils.ReadFile(packageJSONPath)
	if err != nil {
		return false
	}
	// A quick heuristic – any occurrence within scripts is good enough
	scriptRegex := regexp.MustCompile(`"scripts"\s*:\s*{[\s\S]*?}`)
	m := scriptRegex.Find([]byte(content))
	if m == nil {
		return false
	}
	return strings.Contains(string(m), kw)
}

// hasAnyFile returns true if any file in files exists in sourceDir.
func hasAnyFile(sourceDir string, files []string) bool {
	for _, f := range files {
		if utils.FileExists(filepath.Join(sourceDir, f)) {
			return true
		}
	}
	return false
}
