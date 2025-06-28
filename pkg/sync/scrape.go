// Copyright 2025 snowy-jaguar
// Contact: @snowyjaguar (Discord)
// Contact: contact@snowyjaguar.xyz (Email)
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
//     http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sync

import (
	"time"

	"github.com/snowy-jaguar/adguardhomesync-swarm/pkg/metrics"
	"github.com/snowy-jaguar/adguardhomesync-swarm/pkg/types"
)

func (w *worker) startScraping() {
	metrics.Init()
	if w.cfg.API.Metrics.ScrapeInterval == 0 {
		w.cfg.API.Metrics.ScrapeInterval = 30 * time.Second
	}
	if w.cfg.API.Metrics.QueryLogLimit == 0 {
		w.cfg.API.Metrics.QueryLogLimit = 10_000
	}
	l.With(
		"scrape-interval", w.cfg.API.Metrics.ScrapeInterval,
		"query-log-limit", w.cfg.API.Metrics.QueryLogLimit,
	).Info("setup metrics")
	w.scrape()
	for range time.Tick(w.cfg.API.Metrics.ScrapeInterval) {
		w.scrape()
	}
}

func (w *worker) scrape() {
	var iml metrics.InstanceMetricsList

	iml.Metrics = append(iml.Metrics, w.getMetrics(*w.cfg.Origin))
	for _, replica := range w.cfg.Replicas {
		iml.Metrics = append(iml.Metrics, w.getMetrics(replica))
	}
	metrics.UpdateInstances(iml)
}

func (w *worker) getMetrics(inst types.AdGuardInstance) metrics.InstanceMetrics {
	var im metrics.InstanceMetrics
	client, err := w.createClient(inst)
	if err != nil {
		l.With("error", err, "url", w.cfg.Origin.URL).Error("Error creating origin client")
		return im
	}

	im.HostName = inst.Host
	im.Status, _ = client.Status()
	im.Stats, _ = client.Stats()
	im.QueryLog, _ = client.QueryLog(w.cfg.API.Metrics.QueryLogLimit)
	return im
}
