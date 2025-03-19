package sourceanalyzer

import (
	"regexp"

	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/common/utils"
)

// hasListenAndServe checks if any Go source file in the directory
// contains the ListenAndServe function call
func hasListenAndServe(sourceDir string) bool {
	// Regular expression to detect ListenAndServe usage
	listenAndServeRegex := regexp.MustCompile(`\bListenAndServe\s*\(`)

	// Get all Go files recursively
	goFiles, err := utils.FindFiles(sourceDir, "*.go")
	if err != nil {
		log.Warnf("Error finding Go files: %v", err)
		return false
	}

	// Check each file for ListenAndServe usage
	for _, file := range goFiles {
		content, err := utils.ReadFile(file)
		if err != nil {
			log.Warnf("Error reading Go file %s: %v", file, err)
			continue
		}

		if listenAndServeRegex.Match([]byte(content)) {
			return true
		}
	}

	return false
}
