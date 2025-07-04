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

package model

import (
	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"go.uber.org/zap"

	"github.com/snowy-jaguar/adguardhomesync-swarm/pkg/log"
	"github.com/snowy-jaguar/adguardhomesync-swarm/pkg/utils"
)

var _ = Describe("Types", func() {
	Context("DhcpConfigV4", func() {
		DescribeTable("DhcpConfigV4 should not be valid",
			func(v4 DhcpConfigV4) {
				gomega.Ω(v4.isValid()).Should(gomega.BeFalse())
			},
			Entry(`When GatewayIp is nil`, DhcpConfigV4{
				GatewayIp:  nil,
				SubnetMask: utils.Ptr("2.2.2.2"),
				RangeStart: utils.Ptr("3.3.3.3"),
				RangeEnd:   utils.Ptr("4.4.4.4"),
			}),
			Entry(`When GatewayIp is ""`, DhcpConfigV4{
				GatewayIp:  utils.Ptr(""),
				SubnetMask: utils.Ptr("2.2.2.2"),
				RangeStart: utils.Ptr("3.3.3.3"),
				RangeEnd:   utils.Ptr("4.4.4.4"),
			}),
			Entry(`When SubnetMask is nil`, DhcpConfigV4{
				GatewayIp:  utils.Ptr("1.1.1.1"),
				SubnetMask: nil,
				RangeStart: utils.Ptr("3.3.3.3"),
				RangeEnd:   utils.Ptr("4.4.4.4"),
			}),
			Entry(`When SubnetMask is ""`, DhcpConfigV4{
				GatewayIp:  utils.Ptr("1.1.1.1"),
				SubnetMask: utils.Ptr(""),
				RangeStart: utils.Ptr("3.3.3.3"),
				RangeEnd:   utils.Ptr("4.4.4.4"),
			}),
			Entry(`When SubnetMask is nil`, DhcpConfigV4{
				GatewayIp:  utils.Ptr("1.1.1.1"),
				SubnetMask: utils.Ptr("2.2.2.2"),
				RangeStart: nil,
				RangeEnd:   utils.Ptr("4.4.4.4"),
			}),
			Entry(`When SubnetMask is ""`, DhcpConfigV4{
				GatewayIp:  utils.Ptr("1.1.1.1"),
				SubnetMask: utils.Ptr("2.2.2.2"),
				RangeStart: utils.Ptr(""),
				RangeEnd:   utils.Ptr("4.4.4.4"),
			}),
			Entry(`When RangeEnd is nil`, DhcpConfigV4{
				GatewayIp:  utils.Ptr("1.1.1.1"),
				SubnetMask: utils.Ptr("2.2.2.2"),
				RangeStart: utils.Ptr("3.3.3.3"),
				RangeEnd:   nil,
			}),
			Entry(`When RangeEnd is ""`, DhcpConfigV4{
				GatewayIp:  utils.Ptr("1.1.1.1"),
				SubnetMask: utils.Ptr("2.2.2.2"),
				RangeStart: utils.Ptr("3.3.3.3"),
				RangeEnd:   utils.Ptr(""),
			}),
		)
	})
	Context("DhcpConfigV6", func() {
		DescribeTable("DhcpConfigV6 should not be valid",
			func(v6 DhcpConfigV6) {
				gomega.Ω(v6.isValid()).Should(gomega.BeFalse())
			},
			Entry(`When SubnetMask is nil`, DhcpConfigV6{RangeStart: nil}),
			Entry(`When SubnetMask is ""`, DhcpConfigV6{RangeStart: utils.Ptr("")}),
		)
	})
	Context("DNSConfig", func() {
		var (
			cfg *DNSConfig
			l   *zap.SugaredLogger
		)

		BeforeEach(func() {
			cfg = &DNSConfig{
				UsePrivatePtrResolvers: utils.Ptr(true),
			}
			l = log.GetLogger("test")
		})
		Context("Sanitize", func() {
			It("should disable UsePrivatePtrResolvers resolvers is nil ", func() {
				cfg.LocalPtrUpstreams = nil
				cfg.Sanitize(l)
				gomega.Ω(cfg.UsePrivatePtrResolvers).ShouldNot(gomega.BeNil())
				gomega.Ω(*cfg.UsePrivatePtrResolvers).Should(gomega.Equal(false))
			})
			It("should disable UsePrivatePtrResolvers resolvers is empty ", func() {
				cfg.LocalPtrUpstreams = utils.Ptr([]string{})
				cfg.Sanitize(l)
				gomega.Ω(cfg.UsePrivatePtrResolvers).ShouldNot(gomega.BeNil())
				gomega.Ω(*cfg.UsePrivatePtrResolvers).Should(gomega.Equal(false))
			})
		})
	})
})
