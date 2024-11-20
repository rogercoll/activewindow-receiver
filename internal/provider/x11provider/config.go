package x11provider

import (
	"github.com/rogercoll/activewindowreceiver/internal/provider"
)

type Config struct{}

var _ provider.Config = (*Config)(nil)

// Validate implements provider.Config.
func (c *Config) Validate() error {
	return nil
}
