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

package versions_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snowy-jaguar/adguardhomesync-swarm/pkg/versions"
)

var _ = Describe("Versions", func() {
	Context("IsNewerThan", func() {
		It("should correctly parse json", func() {
			Ω(versions.IsNewerThan("v0.106.10", "v0.106.9")).Should(BeTrue())
			Ω(versions.IsNewerThan("v0.106.9", "v0.106.10")).Should(BeFalse())
			Ω(versions.IsNewerThan("v0.106.10", "0.106.9")).Should(BeTrue())
			Ω(versions.IsNewerThan("v0.106.9", "0.106.10")).Should(BeFalse())
		})
	})
	Context("IsSame", func() {
		It("should be the same version", func() {
			Ω(versions.IsSame("v0.106.9", "v0.106.9")).Should(BeTrue())
			Ω(versions.IsSame("0.106.9", "v0.106.9")).Should(BeTrue())
		})
	})
})
