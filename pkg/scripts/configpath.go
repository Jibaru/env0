package scripts

import (
	"context"

	"github.com/Jibaru/env0/pkg/auth"
	"github.com/Jibaru/env0/pkg/client"
	"github.com/Jibaru/env0/pkg/logger"
)

// ConfigPathInput represents the input parameters for the config path operation
type ConfigPathInput struct {
	// Empty since we don't need any input parameters
}

// ConfigPathFn represents a function that performs the config path operation
type ConfigPathFn func(context.Context, ConfigPathInput) error

// NewConfigPath creates a new config path function with injected dependencies
func NewConfigPath(c client.Client, logger logger.Logger) ConfigPathFn {
	return func(ctx context.Context, input ConfigPathInput) error {
		cfgPath, err := auth.GetConfigDir()
		if err != nil {
			return err
		}

		logger.Printf(cfgPath)
		return nil
	}
}
