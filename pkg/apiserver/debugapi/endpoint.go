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

package debugapi

import (
	"net/http"
	"net/url"
	"regexp"

	"github.com/pingcap/tidb-dashboard/pkg/apiserver/model"
)

var (
	ErrMissingRequiredParam = ErrNS.NewType("missing_require_parameter")
	ErrInvalidParam         = ErrNS.NewType("invalid_parameter")
)

type EndpointAPIModel struct {
	ID          string             `json:"id"`
	Component   model.NodeKind     `json:"component"`
	Path        string             `json:"path"`
	Method      EndpointMethod     `json:"method"`
	PathParams  []EndpointAPIParam `json:"path_params"`  // e.g. /stats/dump/{db}/{table} -> db, table
	QueryParams []EndpointAPIParam `json:"query_params"` // e.g. /debug/pprof?seconds=1 -> seconds
}

type EndpointMethod string

const (
	EndpointMethodGet EndpointMethod = http.MethodGet
)

type Request struct {
	Method EndpointMethod
	Host   string
	Port   int
	Path   string
	Query  string
}

func (e *EndpointAPIModel) NewRequest(host string, port int, value map[string]string) (*Request, error) {
	req := &Request{
		Method: e.Method,
		Host:   host,
		Port:   port,
	}

	pathValues, err := transformValues(e.PathParams, value)
	if err != nil {
		return nil, err
	}
	path, err := e.PopulatePath(pathValues)
	if err != nil {
		return nil, err
	}
	req.Path = path

	queryValues, err := transformValues(e.QueryParams, value)
	if err != nil {
		return nil, err
	}
	query, err := e.EncodeQuery(queryValues)
	if err != nil {
		return nil, err
	}
	req.Query = query

	return req, nil
}

var ParamRegexp *regexp.Regexp = regexp.MustCompile(`\{(\w+)\}`)

func (e *EndpointAPIModel) PopulatePath(valMap map[string]string) (string, error) {
	var returnErr error
	replacedPath := e.Path
	replacedPath = ParamRegexp.ReplaceAllStringFunc(replacedPath, func(s string) string {
		if returnErr != nil {
			return s
		}

		key := ParamRegexp.ReplaceAllString(s, "${1}")
		val := valMap[key]
		// means the param can be found in the endpoint path, but not in the param value map
		if val == "" {
			returnErr = ErrMissingRequiredParam.New("path: %s, param: %s", e.Path, key)
			return s
		}

		return val
	})
	return replacedPath, returnErr
}

func (e *EndpointAPIModel) EncodeQuery(valMap map[string]string) (string, error) {
	query := url.Values{}
	for _, q := range e.QueryParams {
		// cuz paramValues was generated by e.Query, paramValues can always find param by q.Name
		val := valMap[q.Name]
		if q.Required && val == "" {
			return "", ErrMissingRequiredParam.New("query param: %s", q.Name)
		}
		query.Add(q.Name, val)
	}
	return query.Encode(), nil
}

func transformValues(params []EndpointAPIParam, values map[string]string) (map[string]string, error) {
	pvMap := map[string]string{}
	for _, p := range params {
		v, err := p.Transform(values[p.Name])
		if err != nil {
			return nil, ErrInvalidParam.WrapWithNoMessage(err)
		}
		pvMap[p.Name] = v
	}
	return pvMap, nil
}

type EndpointAPIParam struct {
	Name     string `json:"name"`
	Required bool   `json:"required"`
	// represents what param is
	Model                EndpointAPIParamModel `json:"model"`
	PreModelTransformer  ModelTransformer      `json:"-"`
	PostModelTransformer ModelTransformer      `json:"-"`
}

// Transform incoming param's value by transformer at endpoint / model definition
func (p *EndpointAPIParam) Transform(value string) (string, error) {
	transfomers := []ModelTransformer{
		p.PreModelTransformer,
		p.Model.Transformer,
		p.PostModelTransformer,
	}

	for _, t := range transfomers {
		if t == nil {
			continue
		}
		v, err := t(value)
		if err != nil {
			return "", err
		}
		value = v
	}

	return value, nil
}

// Transformer can transform the incoming param's value in special scenarios
// Also, now are used as validation function
type ModelTransformer func(value string) (string, error)

type EndpointAPIParamModel struct {
	Type        string           `json:"type"`
	Transformer ModelTransformer `json:"-"`
}