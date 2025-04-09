package utils

import (
	"fmt"
	"time"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/unbindapp/unbind-api/internal/common/log"
)

// GetexportedPortsFromRegistry detects exposed ports from a docker manifest, in 8080/tcp format
func GetExposedPortsFromRegistry(imageName string) ([]string, error) {
	// Parse the image reference
	ref, err := name.ParseReference(imageName)
	if err != nil {
		return nil, fmt.Errorf("invalid image name %s: %v", imageName, err)
	}

	// Get the image configuration
	// ! TODO - support regcred?
	start := time.Now()
	img, err := remote.Image(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return nil, fmt.Errorf("failed to get image from registry: %v", err)
	}

	end := time.Now()
	log.Infof("Docker remote.Image took %s", end.Sub(start).String())

	// Get the image config
	start = time.Now()
	configFile, err := img.ConfigFile()
	if err != nil {
		return nil, fmt.Errorf("failed to get image config: %v", err)
	}
	end = time.Now()
	log.Infof("Docker img.ConfigFile() took %s", end.Sub(start).String())

	// Extract exposed ports
	ports := []string{}
	for port := range configFile.Config.ExposedPorts {
		ports = append(ports, port)
	}

	return ports, nil
}
