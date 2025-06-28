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
	"bytes"
	_ "embed"
	"os"
	"runtime"
	"sort"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"

	"github.com/snowy-jaguar/adguardhomesync-swarm/pkg/client"
	"github.com/snowy-jaguar/adguardhomesync-swarm/pkg/types"
	"github.com/snowy-jaguar/adguardhomesync-swarm/version"
)

//go:embed print-config.md
var printConfigTemplate string

func (ac *AppConfig) Print() error {
	originVersion := aghVersion(*ac.cfg.Origin)
	var replicaVersions []string
	for _, replica := range ac.cfg.Replicas {
		replicaVersions = append(replicaVersions, aghVersion(replica))
	}

	out, err := ac.printInternal(os.Environ(), originVersion, replicaVersions)
	if err != nil {
		return err
	}

	logger.Infof(
		"Printing adguardhome-sync aggregated config (THE APPLICATION WILL NOT START IN THIS MODE):\n%s",
		out,
	)

	return nil
}

func aghVersion(i types.AdGuardInstance) string {
	cl, err := client.New(i)
	if err != nil {
		return "N/A"
	}
	stats, err := cl.Status()
	if err != nil {
		return "N/A"
	}
	return stats.Version
}

func (ac *AppConfig) printInternal(env []string, originVersion string, replicaVersions []string) (string, error) {
	config, err := yaml.Marshal(ac.Get())
	if err != nil {
		return "", err
	}

	funcMap := template.FuncMap{
		// The name "inc" is what the function will be called in the template text.
		"inc": func(i int) int {
			return i + 1
		},
	}

	t, err := template.New("printConfigTemplate").Funcs(funcMap).Parse(printConfigTemplate)
	if err != nil {
		return "", err
	}

	sort.Strings(env)

	var buf bytes.Buffer

	err = t.Execute(&buf, map[string]any{
		"Version":              version.Version,
		"Build":                version.Build,
		"OperatingSystem":      runtime.GOOS,
		"Architecture":         runtime.GOARCH,
		"AggregatedConfig":     string(config),
		"ConfigFilePath":       ac.filePath,
		"ConfigFileContent":    ac.content,
		"EnvironmentVariables": strings.Join(env, "\n"),
		"OriginVersion":        originVersion,
		"ReplicaVersions":      replicaVersions,
	})
	return buf.String(), err
}
