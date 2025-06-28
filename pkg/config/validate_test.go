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

package config

import (
	"github.com/go-faker/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"

	"github.com/snowy-jaguar/adguardhomesync-swarm/pkg/types"
)

var _ = Describe("Config", func() {
	Context("validateSchema", func() {
		DescribeTable("validateSchema config",
			func(configFile string, expectFail bool) {
				err := validateSchema(configFile)
				if expectFail {
					Ω(err).Should(HaveOccurred())
				} else {
					Ω(err).ShouldNot(HaveOccurred())
				}
			},
			Entry(`Should be valid`, "../../testdata/config/config-valid.yaml", false),
			Entry(`Should be valid if file doesn't exist`, "../../testdata/config/foo.bar", false),
			Entry(`Should fail if file is not yaml`, "../../go.mod", true),
		)
		It("validate config with all fields randomly populated", func() {
			cfg := &types.Config{}

			err := faker.FakeData(cfg)
			Ω(err).ShouldNot(HaveOccurred())

			data, err := yaml.Marshal(&cfg)
			Ω(err).ShouldNot(HaveOccurred())

			err = validateYAML(data)
			Ω(err).ShouldNot(HaveOccurred())
		})
		It("validate config with empty file", func() {
			var data []byte
			err := validateYAML(data)
			Ω(err).ShouldNot(HaveOccurred())
		})
	})
})
