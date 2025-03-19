package sourceanalyzer

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/common/utils"
)

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
