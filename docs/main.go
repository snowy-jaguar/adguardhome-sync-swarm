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

// Print the available environment variables
package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/snowy-jaguar/adguardhomesync-swarm/pkg/types"
)

func main() {
	_, _ = fmt.Println("| Name | Type | Description |")
	_, _ = fmt.Println("| :--- | ---- |:----------- |")
	printEnvTags(reflect.TypeOf(types.Config{}), "")
}

// printEnvTags recursively prints all fields with `env` tags.
func printEnvTags(t reflect.Type, prefix string) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return
	}

	for _, field := range reflect.VisibleFields(t) {
		if field.PkgPath != "" { // unexported field
			continue
		}

		envTag := field.Tag.Get("env")
		if envTag == "" {
			switch field.Name {
			case "Origin":
				envTag = "ORIGIN"
			case "Replica":
				envTag = "REPLICA#"
			}
		}
		combinedTag := envTag
		if prefix != "" && envTag != "" {
			combinedTag = prefix + "_" + envTag
		} else if prefix != "" {
			combinedTag = prefix
		}

		ft := field.Type
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}

		if ft.Kind() == reflect.Struct && ft.Name() != "Time" { // skip time.Time
			printEnvTags(ft, strings.TrimSuffix(combinedTag, "_"))
		} else if envTag != "" {
			envVar := strings.Trim(combinedTag, "_") + " (" + ft.Kind().String() + ")"
			docs := field.Tag.Get("documentation")

			_, _ = fmt.Printf("| %s | %s | %s |\n", envVar, ft.Kind().String(), docs)
		}
	}
}
