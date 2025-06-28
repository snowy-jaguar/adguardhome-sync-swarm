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

package client

import (
	"context"

	"github.com/snowy-jaguar/adguardhomesync-swarm/pkg/client/model"
)

type Client interface {
	Host(ctx context.Context) string
	GetServerStatus(ctx context.Context) (*model.ServerStatus, error)

	GetFilteringStatus(ctx context.Context) (*model.FilterStatus, error)
	SetFilteringConfig(ctx context.Context, config model.FilterConfig) error
}
