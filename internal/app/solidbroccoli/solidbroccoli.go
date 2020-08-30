package solidbroccoli

import (
	"fmt"

	"github.com/dstdfx/solid-broccoli/internal/pkg/config"
)

// StartService runs main service's goroutine.
func StartService() error {
	if err := config.CheckConfig(); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	// TODO: run interfaces
	// TODO: gsh

	return nil
}
