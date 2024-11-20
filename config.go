// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
package activewindowreceiver

import (
	"errors"
	"fmt"
	"time"

	"github.com/rogercoll/activewindowreceiver/internal/provider"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
)

const (
	providersKey = "providers"
)

// Config defines configuration for HostMetrics receiver.
type Config struct {
	scraperhelper.ControllerConfig `mapstructure:",squash"`

	// Providers specifies the sources to get the active window from.
	Providers map[string]provider.Config `mapstructure:"-"`

	Precision time.Duration `mapstructure:"precision"`
}

var _ component.Config = (*Config)(nil)

// Validate checks the receiver configuration is valid
func (cfg *Config) Validate() error {
	if len(cfg.Providers) == 0 {
		return errors.New("must specify at least one active window provider")
	}
	return nil
}

// Unmarshal a config.Parser into the config struct.
func (cfg *Config) Unmarshal(componentParser *confmap.Conf) error {
	if componentParser == nil {
		return nil
	}

	// load the non-dynamic config normally
	err := componentParser.Unmarshal(cfg, confmap.WithIgnoreUnused())
	if err != nil {
		return err
	}

	// dynamically load the individual providers configs based on the key name
	cfg.Providers = map[string]provider.Config{}

	// retrieve `providers` configuration section
	providersSection, err := componentParser.Sub(providersKey)
	if err != nil {
		return err
	}

	// loop through all defined providers and load their configuration
	for key := range providersSection.ToStringMap() {
		factory, ok := getProviderFactory(key)
		if !ok {
			return fmt.Errorf("invalid provider key: %s", key)
		}

		providerCfg := factory.CreateDefaultConfig()
		providerSection, err := providersSection.Sub(key)
		if err != nil {
			return err
		}
		err = providerSection.Unmarshal(providerCfg)
		if err != nil {
			return fmt.Errorf("error reading settings for provider type %q: %w", key, err)
		}

		cfg.Providers[key] = providerCfg
	}

	return nil
}
