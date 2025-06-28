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
	"net/http"

	"github.com/go-resty/resty/v2"

	"github.com/snowy-jaguar/adguardhomesync-swarm/pkg/client/model"
)

var _ model.HttpRequestDoer = &adapter{}

func RestyAdapter(r *resty.Client) model.HttpRequestDoer {
	return &adapter{
		client: r,
	}
}

type adapter struct {
	client *resty.Client
}

func (a adapter) Do(req *http.Request) (*http.Response, error) {
	r, err := a.client.R().
		SetHeaderMultiValues(req.Header).
		Execute(req.Method, req.URL.String())
	return r.RawResponse, err
}
