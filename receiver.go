// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package activewindowreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/activewindowreceiver"

import (
	"context"
	"sync"
	"time"

	"github.com/jezek/xgb"
	"github.com/rogercoll/activewindowreceiver/internal/metadata"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
)

const entityType = "host"

type activewindowReceiver struct {
	cfg *Config

	connection *xgb.Conn
	cancel     context.CancelFunc

	mb *metadata.MetricsBuilder

	windows sync.Map
}

func createMetricsReceiver(
	_ context.Context,
	params receiver.Settings,
	config component.Config,
	consumer consumer.Metrics,
) (receiver.Metrics, error) {
	X, err := xgb.NewConn()
	if err != nil {
		return nil, err
	}

	activeConfig := config.(*Config)
	recv := activewindowReceiver{
		cfg:        activeConfig,
		connection: X,
		mb:         metadata.NewMetricsBuilder(metadata.DefaultMetricsBuilderConfig(), params),
	}

	scrp, err := scraperhelper.NewScraperWithoutType(recv.scrape, scraperhelper.WithStart(recv.start), scraperhelper.WithShutdown(recv.shutdown))
	if err != nil {
		return nil, err
	}
	return scraperhelper.NewScraperControllerReceiver(&recv.cfg.ControllerConfig, params, consumer, scraperhelper.AddScraperWithType(metadata.Type, scrp))
}

func (ar *activewindowReceiver) start(ctx context.Context, _ component.Host) error {
	ctx, ar.cancel = context.WithCancel(ctx)
	ticker := time.NewTicker(ar.cfg.Precision)
	go func() {
		for {
			select {
			case <-ticker.C:
				windowId, windowName := activeWindow(ar.connection)
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
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()

	return nil
}

func (ar *activewindowReceiver) scrape(ctx context.Context) (pmetric.Metrics, error) {
	now := pcommon.NewTimestampFromTime(time.Now())
	ar.windows.Range(func(windowId, value any) bool {
		windows := value.(*sync.Map)
		windows.Range(func(windowName, value any) bool {
			time := value.(float64)
			windowIdStr := windowId.(string)
			windowNameStr := windowName.(string)
			ar.mb.RecordSystemGuiWindowTimeDataPoint(now, time, windowIdStr, windowNameStr)
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
