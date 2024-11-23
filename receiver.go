// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package activewindowreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/activewindowreceiver"

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/rogercoll/activewindowreceiver/internal/metadata"
	"github.com/rogercoll/activewindowreceiver/internal/provider"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
	"go.uber.org/zap"
)

const entityType = "host"

type activewindowReceiver struct {
	cfg    *Config
	logger *zap.Logger

	providers []provider.ActiveWindowProvider
	cancel    context.CancelFunc

	mb *metadata.MetricsBuilder

	windows sync.Map
}

func createMetricsReceiver(
	ctx context.Context,
	set receiver.Settings,
	config component.Config,
	consumer consumer.Metrics,
) (receiver.Metrics, error) {
	activeConfig := config.(*Config)
	providers := make([]provider.ActiveWindowProvider, 0, len(activeConfig.Providers))

	for key, cfg := range activeConfig.Providers {
		factory := providerFactories[key]
		if factory == nil {
			return nil, fmt.Errorf("activewindow provider factory not found for key: %q", key)
		}

		provider, err := factory.CreateActiveWindowProvider(ctx, set, cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create provider for key %q: %w", key, err)
		}

		providers = append(providers, provider)
	}

	recv := activewindowReceiver{
		cfg:       activeConfig,
		logger:    set.Logger,
		mb:        metadata.NewMetricsBuilder(metadata.DefaultMetricsBuilderConfig(), set),
		providers: providers,
	}

	scrp, err := scraperhelper.NewScraperWithoutType(recv.scrape, scraperhelper.WithStart(recv.start), scraperhelper.WithShutdown(recv.shutdown))
	if err != nil {
		return nil, err
	}
	return scraperhelper.NewScraperControllerReceiver(&recv.cfg.ControllerConfig, set, consumer, scraperhelper.AddScraperWithType(metadata.Type, scrp))
}

func (ar *activewindowReceiver) start(ctx context.Context, _ component.Host) error {
	ctx, ar.cancel = context.WithCancel(ctx)
	ticker := time.NewTicker(ar.cfg.Precision)
	go func() {
		for {
			select {
			case <-ticker.C:
				// windowId, windowName := activeWindow(ar.connection)
				for _, windowProvider := range ar.providers {
					windowId, windowName, err := windowProvider.ActiveWindow(ctx)
					if err != nil {
						ar.logger.Error(err.Error())
						continue
					}
					windowsId, ok := ar.windows.Load(windowId)
					if !ok {
						var windowNames sync.Map
						windowNames.Store(windowName, ar.cfg.Precision.Seconds())
						ar.windows.Store(windowId, &windowNames)
					} else {
						value, ok := windowsId.(*sync.Map).Load(windowName)
						if ok {
							windowsId.(*sync.Map).Store(windowName, value.(float64)+ar.cfg.Precision.Seconds())
						} else {
							windowsId.(*sync.Map).Store(windowName, ar.cfg.Precision.Seconds())
						}
					}
				}
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()

	return nil
}

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

func parseWindowName(windowName string) []string {
	windows := nonAlphanumericRegex.Split(strings.TrimSpace(windowName), -1)
	for i := range windows {
		windows[i] = strings.TrimSpace(windows[i])
	}
	return windows
}

func (ar *activewindowReceiver) scrape(ctx context.Context) (pmetric.Metrics, error) {
	now := pcommon.NewTimestampFromTime(time.Now())
	ar.windows.Range(func(windowId, value any) bool {
		windows := value.(*sync.Map)
		windows.Range(func(windowName, value any) bool {
			time := value.(float64)
			windowIdStr := windowId.(string)
			parsedWindow := parseWindowName(windowName.(string))
			ar.mb.RecordSystemGuiWindowTimeDataPoint(now, time, windowIdStr, windowName.(string), parsedWindow[len(parsedWindow)-1])
			return true
		})
		return true
	})

	return ar.mb.Emit(), nil
}

func (hmr *activewindowReceiver) shutdown(_ context.Context) error {
	if hmr.cancel != nil {
		hmr.cancel()
	}
	return nil
}
