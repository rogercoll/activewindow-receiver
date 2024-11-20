// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/receiver"
)

type Config interface {
	component.ConfigValidator
}

type ActiveWindowProvider interface {
	ActiveWindow(context.Context) (string, string, error)
}

type ActiveWindowProviderFactory interface {
	CreateDefaultConfig() Config

	CreateActiveWindowProvider(ctx context.Context, settings receiver.Settings, cfg Config) (ActiveWindowProvider, error)
}
