// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package activewindowreceiver

import (
	"time"

	"github.com/rogercoll/activewindowreceiver/internal/metadata"
	"github.com/rogercoll/activewindowreceiver/internal/provider"
	"github.com/rogercoll/activewindowreceiver/internal/provider/x11provider"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
)

var providerFactories = map[string]provider.ActiveWindowProviderFactory{
	x11provider.TypeStr: &x11provider.Factory{},
}

// CreateDefaultConfig creates the default configuration for the Scraper.
func createDefaultReceiverConfig() *Config {
	return &Config{
		ControllerConfig: scraperhelper.NewDefaultControllerConfig(),
		Precision:        1 * time.Second,
	}
}

func getProviderFactory(key string) (provider.ActiveWindowProviderFactory, bool) {
	if factory, ok := providerFactories[key]; ok {
		return factory, true
	}

	return nil, false
}

func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		metadata.Type,
		func() component.Config {
			return createDefaultReceiverConfig()
		},
		receiver.WithMetrics(createMetricsReceiver, metadata.MetricsStability))
}
