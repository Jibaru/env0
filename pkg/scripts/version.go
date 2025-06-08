package scripts

import (
	"context"

	"github.com/Jibaru/env0/pkg/logger"
)

// Version represents the current version of the application
const Version = "v0.2.0"

// VersionFn represents a function that performs the version operation
type VersionFn func(context.Context) error

// NewVersion creates a new version function with injected dependencies
func NewVersion(logger logger.Logger) VersionFn {
	return func(ctx context.Context) error {
		logger.Printf(Version)
		return nil
	}
}
