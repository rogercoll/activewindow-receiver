// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
package activewindowreceiver

import (
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
)

const (
	scrapersKey = "scrapers"
)

// Config defines configuration for HostMetrics receiver.
type Config struct {
	scraperhelper.ControllerConfig `mapstructure:",squash"`

	Precision time.Duration `mapstructure:"precision"`
}

var _ component.Config = (*Config)(nil)

// Validate checks the receiver configuration is valid
func (cfg *Config) Validate() error {
	return nil
}
