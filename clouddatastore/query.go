// Copyright 2023 Stichting ThingsIX Foundation
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
//
// SPDX-License-Identifier: Apache-2.0

package clouddatastore

import (
	"cloud.google.com/go/datastore"
)

func QueryBeginsWith(query *datastore.Query, field, beginsWith string) *datastore.Query {
	start := beginsWith
	endB := []byte(beginsWith)
	endB[len(endB)-1] = endB[len(endB)-1] + 1
	end := string(endB)

	query = query.FilterField(field, ">=", start)
	query = query.FilterField(field, "<", end)

	return query
}
