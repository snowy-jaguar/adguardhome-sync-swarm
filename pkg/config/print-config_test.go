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
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snowy-jaguar/adguardhomesync-swarm/pkg/test/matchers"
	"github.com/snowy-jaguar/adguardhomesync-swarm/pkg/types"
	"github.com/snowy-jaguar/adguardhomesync-swarm/version"
)

var _ = Describe("AppConfig", func() {
	var (
		ac  *AppConfig
		env []string
	)
	BeforeEach(func() {
		ac = &AppConfig{
			cfg: &types.Config{
				Origin: &types.AdGuardInstance{
					URL: "https://ha.xxxx.net:3000",
				},
			},
			content: `
origin:
  url: https://ha.xxxx.net:3000
`,
		}
		env = []string{"FOO=foo", "BAR=bar"}
	})
	Context("printInternal", func() {
		It("should printInternal config without file", func() {
			out, err := ac.printInternal(env, "v0.0.1", []string{"v0.0.2"})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(out).
				Should(matchers.EqualIgnoringLineEndings(fmt.Sprintf(expected(1), version.Version, version.Build, runtime.GOOS, runtime.GOARCH)))
		})
		It("should printInternal config with file", func() {
			ac.filePath = "config.yaml"
			out, err := ac.printInternal(env, "v0.0.1", []string{"v0.0.2"})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(out).
				Should(matchers.EqualIgnoringLineEndings(fmt.Sprintf(expected(2), version.Version, version.Build, runtime.GOOS, runtime.GOARCH)))
		})
	})
})

func expected(id int) string {
	b, err := os.ReadFile(
		filepath.Join("..", "..", "testdata", "config", fmt.Sprintf("print-config_test_expected%d.md", id)),
	)
	Ω(err).ShouldNot(HaveOccurred())
	return string(b)
}
