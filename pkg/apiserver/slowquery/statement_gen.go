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

package slowquery

import (
	"fmt"
	"strings"

	"github.com/thoas/go-funk"
)

func (s *Service) genSelectStmt(tableColumns []string, reqJSONColumns []string) (string, error) {
	fields := getFieldsAndTags()

	// use required fields filter when not all fields are requested
	if reqJSONColumns[0] != "*" {
		// These three fields are the most basic information of a slow query record and should contain them
		requiredFields := funk.UniqString(append(reqJSONColumns, "digest", "connection_id", "timestamp"))
		fields = funk.Filter(fields, func(f Field) bool {
			return funk.Contains(requiredFields, f.JSONName)
		}).([]Field)
	}

	// We have both TiDB 4.x and TiDB 5.x columns listed in the model. Filter out columns that do not exist in current version TiDB schema.
	fields = funk.Filter(fields, func(f Field) bool {
		hasProjection := f.Projection != ""
		isTableColumnValid := funk.Contains(tableColumns, f.ColumnName)
		return hasProjection || isTableColumnValid
	}).([]Field)

	if len(fields) == 0 {
		return "", fmt.Errorf("unknown request fields: %q", reqJSONColumns)
	}

	stmt := funk.Map(fields, func(f Field) string {
		if f.Projection == "" {
			return f.ColumnName
		}
		return fmt.Sprintf("%s AS %s", f.Projection, f.ColumnName)
	}).([]string)
	return strings.Join(stmt, ", "), nil
}

func (s *Service) genOrderStmt(tableColumns []string, orderBy string, isDesc bool) (string, error) {
	var order string
	// to handle the special case: timestamp
	// Order by column instead of expression, see related optimization in TiDB: https://github.com/pingcap/tidb/pull/20750
	if orderBy == "timestamp" {
		order = "Time"
	} else {
		// We have both TiDB 4.x and TiDB 5.x columns listed in the model. Filter out columns that do not exist in current version TiDB schema.
		fields := funk.Filter(getFieldsAndTags(), func(f Field) bool {
			return funk.Contains(tableColumns, f.ColumnName)
		}).([]Field)
		orderField := funk.Find(fields, func(f Field) bool {
			return f.JSONName == orderBy
		})
		if orderField == nil {
			return "", fmt.Errorf("unknown order by %s", orderBy)
		}

		order = orderField.(Field).ColumnName
	}

	if isDesc {
		order = fmt.Sprintf("%s DESC", order)
	} else {
		order = fmt.Sprintf("%s ASC", order)
	}

	return order, nil
}
