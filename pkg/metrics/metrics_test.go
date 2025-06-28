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

package metrics

import (
	"github.com/go-faker/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/utils/ptr"

	"github.com/snowy-jaguar/adguardhomesync-swarm/pkg/client/model"
)

var _ = Describe("Metrics", func() {
	BeforeEach(func() {
		stats = make(OverallStats)
	})
	Context("UpdateInstances / getStats", func() {
		It("generate correct stats", func() {
			UpdateInstances(InstanceMetricsList{[]InstanceMetrics{
				{HostName: "foo", Status: &model.ServerStatus{}, Stats: &model.Stats{
					NumDnsQueries: ptr.To(100),
					DnsQueries: ptr.To(
						[]int{10, 20, 30, 40, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					),
				}},
				{HostName: "bar", Status: &model.ServerStatus{}, Stats: &model.Stats{
					NumDnsQueries: ptr.To(200),
					DnsQueries: ptr.To(
						[]int{20, 40, 60, 80, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					),
				}},
				{HostName: "aaa", Status: &model.ServerStatus{}, Stats: &model.Stats{
					NumDnsQueries: ptr.To(300),
					DnsQueries: ptr.To(
						[]int{30, 60, 90, 120, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					),
				}},
			}})
			Ω(stats).Should(HaveKey("foo"))
			Ω(stats["foo"].NumDnsQueries).Should(Equal(ptr.To(100)))
			Ω(stats).Should(HaveKey("bar"))
			Ω(stats["bar"].NumDnsQueries).Should(Equal(ptr.To(200)))
			Ω(stats).Should(HaveKey("aaa"))
			Ω(stats["aaa"].NumDnsQueries).Should(Equal(ptr.To(300)))

			os := getStats()
			tot := os.Total()
			Ω(*tot.NumDnsQueries).Should(Equal(600))
			Ω(
				*tot.DnsQueries,
			).Should(Equal([]int{60, 120, 180, 240, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}))

			foo := os["foo"]
			bar := os["bar"]
			aaa := os["aaa"]

			Ω(*foo.NumDnsQueries).Should(Equal(100))
			Ω(*bar.NumDnsQueries).Should(Equal(200))
			Ω(*aaa.NumDnsQueries).Should(Equal(300))
			Ω(
				*foo.DnsQueries,
			).Should(Equal([]int{10, 20, 30, 40, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}))
			Ω(
				*bar.DnsQueries,
			).Should(Equal([]int{20, 40, 60, 80, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}))
			Ω(
				*aaa.DnsQueries,
			).Should(Equal([]int{30, 60, 90, 120, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}))
		})
	})
	Context("StatsGraph", func() {
		var metrics InstanceMetricsList
		BeforeEach(func() {
			err := faker.FakeData(&metrics)
			Ω(err).ShouldNot(HaveOccurred())
		})
		It("should provide correct results with faked values", func() {
			UpdateInstances(metrics)

			_, dns, blocked, malware, adult := StatsGraph()

			verifyStats(dns)
			verifyStats(blocked)
			verifyStats(malware)
			verifyStats(adult)
		})
	})
})

func verifyStats(lines []Line) {
	var total Line
	sum := make([]int, len(lines[0].Data))
	for _, l := range lines {
		if l.Title == labelTotal {
			total = l
		} else {
			for i, d := range l.Data {
				sum[i] += d
			}
		}
	}
	Ω(sum).Should(Equal(total.Data))
}
