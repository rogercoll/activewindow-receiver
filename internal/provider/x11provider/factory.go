package x11provider

import (
	"context"

	"github.com/rogercoll/activewindowreceiver/internal/provider"
	"go.opentelemetry.io/collector/receiver"
)

const (
	// TypeStr the value of "type" key in configuration.
	TypeStr = "x11"
)

type Factory struct{}

var _ provider.ActiveWindowProviderFactory = (*Factory)(nil)

// CreateDefaultConfig creates the default configuration for the Provider.
func (f *Factory) CreateDefaultConfig() provider.Config {
	return &Config{}
}

func (f *Factory) CreateActiveWindowProvider(_ context.Context, _ receiver.Settings, cfg provider.Config) (provider.ActiveWindowProvider, error) {
	x11Config := cfg.(*Config)
	return newX11Provider(x11Config)
}
