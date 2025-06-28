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
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/snowy-jaguar/adguardhomesync-swarm/pkg/types"
)

func readFile(cfg *types.Config, path string) (string, error) {
	var content string
	if _, err := os.Stat(path); err == nil {
		b, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		content = string(b)
		if err := yaml.Unmarshal(b, cfg); err != nil {
			return "", err
		}
	}
	return content, nil
}

func configFilePath(configFile string) (string, error) {
	if configFile == "" {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, ".adguardhome-sync.yaml"), nil
	}
	return configFile, nil
}
