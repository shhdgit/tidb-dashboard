// Copyright 2021 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package schema

import "net"

var (
	ErrIPFormat = ErrNS.NewType("invalid_ip_format")
)

type ModelTransformer func(value string) (string, error)

type EndpointAPIModel struct {
	Type        string           `json:"type"`
	Transformer ModelTransformer `json:"-"`
}

var EndpointAPIModelText EndpointAPIModel = EndpointAPIModel{
	Type: "text",
}

var EndpointAPIModelIP EndpointAPIModel = EndpointAPIModel{
	Type: "ip",
	Transformer: func(value string) (string, error) {
		ip := net.ParseIP(value)
		if ip == nil {
			return "", ErrIPFormat.New("input: %s", value)
		}
		return value, nil
	},
}
