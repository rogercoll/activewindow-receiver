// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package activewindowreceiver

import (
	"time"

	"github.com/rogercoll/activewindowreceiver/internal/metadata"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
)

// This file implements Factory for CPU scraper.

// CreateDefaultConfig creates the default configuration for the Scraper.
func createDefaultReceiverConfig() *Config {
	return &Config{
		ControllerConfig: scraperhelper.NewDefaultControllerConfig(),
		Precision:        1 * time.Second,
	}
}

func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		metadata.Type,
		func() component.Config {
			return createDefaultReceiverConfig()
		},
		receiver.WithMetrics(createMetricsReceiver, metadata.MetricsStability))
}
