//go:build tools
// +build tools

package main

//! TODO - go 1.24 fixes this so IDE won't whine - keep that in mind, when we update to 1.24
import (
	_ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
)
