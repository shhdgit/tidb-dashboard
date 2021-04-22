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

import (
	"fmt"
	"net/http"

	"github.com/pingcap/tidb-dashboard/pkg/apiserver/model"
)

var tidbIPParam EndpointAPIParam = EndpointAPIParam{
	Name:   "tidb_ip",
	Prefix: "{",
	Suffix: "}:10080",
	Model:  EndpointAPIModelIP,
	PostModelTransformer: func(value string) (string, error) {
		return fmt.Sprintf("%s:10080", value), nil
	},
}

var tidbStatsDump EndpointAPI = EndpointAPI{
	ID:        "tidb_stats_dump",
	Component: model.NodeKindTiDB,
	Path:      "/stats/dump/{db}/{table}",
	Method:    http.MethodGet,
	Host:      tidbIPParam,
	Segment: []EndpointAPISegmentParam{
		NewEndpointAPISegmentParam(EndpointAPIParam{
			Name:  "db",
			Model: EndpointAPIModelText,
		}),
		NewEndpointAPISegmentParam(EndpointAPIParam{
			Name:  "table",
			Model: EndpointAPIModelText,
		}),
	},
	Query: []EndpointAPIParam{},
}

var tidbStatsDumpWithTimestamp EndpointAPI = EndpointAPI{
	ID:        "tidb_stats_dump_timestamp",
	Component: model.NodeKindTiDB,
	Path:      "/stats/dump/{db}/{table}/{yyyyMMddHHmmss}",
	Method:    http.MethodGet,
	Host:      tidbIPParam,
	Segment: []EndpointAPISegmentParam{
		NewEndpointAPISegmentParam(EndpointAPIParam{
			Name:  "db",
			Model: EndpointAPIModelText,
		}),
		NewEndpointAPISegmentParam(EndpointAPIParam{
			Name:  "table",
			Model: EndpointAPIModelText,
		}),
		NewEndpointAPISegmentParam(EndpointAPIParam{
			Name:  "yyyyMMddHHmmss",
			Model: EndpointAPIModelText,
		}),
	},
	Query: []EndpointAPIParam{},
}

var tidbConfig EndpointAPI = EndpointAPI{
	ID:        "tidb_config",
	Component: model.NodeKindTiDB,
	Path:      "/settings",
	Method:    http.MethodGet,
	Host:      tidbIPParam,
	Segment:   []EndpointAPISegmentParam{},
	Query:     []EndpointAPIParam{},
}

var tidbSchema EndpointAPI = EndpointAPI{
	ID:        "tidb_schema",
	Component: model.NodeKindTiDB,
	Path:      "/schema",
	Method:    http.MethodGet,
	Host:      tidbIPParam,
}

var tidbSchemaWithDB EndpointAPI = EndpointAPI{
	ID:        "tidb_schema_db",
	Component: model.NodeKindTiDB,
	Path:      "/schema/{db}",
	Method:    http.MethodGet,
	Host:      tidbIPParam,
	Segment: []EndpointAPISegmentParam{
		NewEndpointAPISegmentParam(EndpointAPIParam{
			Name:  "db",
			Model: EndpointAPIModelText,
		}),
	},
	Query: []EndpointAPIParam{},
}

var tidbSchemaWithDBTable EndpointAPI = EndpointAPI{
	ID:        "tidb_schema_db_table",
	Component: model.NodeKindTiDB,
	Path:      "/schema/{db}/{table}",
	Method:    http.MethodGet,
	Host:      tidbIPParam,
	Segment: []EndpointAPISegmentParam{
		NewEndpointAPISegmentParam(EndpointAPIParam{
			Name:  "db",
			Model: EndpointAPIModelText,
		}),
		NewEndpointAPISegmentParam(EndpointAPIParam{
			Name:  "table",
			Model: EndpointAPIModelText,
		}),
	},
	Query: []EndpointAPIParam{},
}

var tidbSchemaWithTableID EndpointAPI = EndpointAPI{
	ID:        "tidb_schema_tableid",
	Component: model.NodeKindTiDB,
	Path:      "/db-table/{tableID}",
	Method:    http.MethodGet,
	Host:      tidbIPParam,
	Segment: []EndpointAPISegmentParam{
		NewEndpointAPISegmentParam(EndpointAPIParam{
			Name:  "tableID",
			Model: EndpointAPIModelText,
		}),
	},
	Query: []EndpointAPIParam{},
}

var tidbDDLHistory EndpointAPI = EndpointAPI{
	ID:        "tidb_ddl_history",
	Component: model.NodeKindTiDB,
	Path:      "/ddl/history",
	Method:    http.MethodGet,
	Host:      tidbIPParam,
	Segment:   []EndpointAPISegmentParam{},
	Query:     []EndpointAPIParam{},
}

var tidbInfo EndpointAPI = EndpointAPI{
	ID:        "tidb_info",
	Component: model.NodeKindTiDB,
	Path:      "/info",
	Method:    http.MethodGet,
	Host:      tidbIPParam,
	Segment:   []EndpointAPISegmentParam{},
	Query:     []EndpointAPIParam{},
}

var tidbInfoAll EndpointAPI = EndpointAPI{
	ID:        "tidb_info_all",
	Component: model.NodeKindTiDB,
	Path:      "/info/all",
	Method:    http.MethodGet,
	Host:      tidbIPParam,
	Segment:   []EndpointAPISegmentParam{},
	Query:     []EndpointAPIParam{},
}

var tidbRegionsMeta EndpointAPI = EndpointAPI{
	ID:        "tidb_regions_meta",
	Component: model.NodeKindTiDB,
	Path:      "/regions/meta",
	Method:    http.MethodGet,
	Host:      tidbIPParam,
	Segment:   []EndpointAPISegmentParam{},
	Query:     []EndpointAPIParam{},
}

var EndpointAPIList []EndpointAPI = []EndpointAPI{
	tidbStatsDump,
	tidbStatsDumpWithTimestamp,
	tidbConfig,
	tidbSchema,
	tidbSchemaWithDB,
	tidbSchemaWithTableID,
	tidbDDLHistory,
	tidbInfo,
	tidbInfoAll,
	tidbRegionsMeta,
}