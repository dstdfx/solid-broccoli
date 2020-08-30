package solidbroccoli

import (
	"fmt"

	"github.com/dstdfx/solid-broccoli/internal/pkg/config"
	"go.uber.org/zap"
)

// StartService runs main service's goroutine.
func StartService(log *zap.Logger) error {
	if err := config.CheckConfig(); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	log.Debug("start service")

	// TODO: run interfaces
	// TODO: gsh

	return nil
}
