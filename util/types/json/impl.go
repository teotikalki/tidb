// Copyright 2017 PingCAP, Inc.
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

package json

import (
	"encoding/json"
	"strings"
)

func parseFromString(s string) (JSON, error) {
	if len(s) == 0 {
		return nil, ErrInvalidJSONText.GenByArgs("The document is empty")
	}
	var in interface{}
	if err := json.Unmarshal([]byte(s), &in); err != nil {
		return nil, ErrInvalidJSONText.GenByArgs(err)
	}
	return normalize(in), nil
}

func dumpToString(j JSON) string {
	bytes, _ := json.Marshal(j)
	return strings.Trim(string(bytes), "\n")
}

func normalize(in interface{}) JSON {
	switch t := in.(type) {
	case bool:
		var literal = new(jsonLiteral)
		if t {
			*literal = jsonLiteral(0x01)
		} else {
			*literal = jsonLiteral(0x02)
		}
		return literal
	case nil:
		var literal = new(jsonLiteral)
		*literal = 0x00
		return literal
	case float64:
		var f64 = new(jsonDouble)
		*f64 = jsonDouble(t)
		return f64
	case string:
		var s = new(jsonString)
		*s = jsonString(t)
		return s
	case map[string]interface{}:
		var object = new(jsonObject)
		*object = make(map[string]JSON, len(t))
		for key, value := range t {
			(*object)[key] = normalize(value)
		}
		return object
	case []interface{}:
		var array = new(jsonArray)
		*array = make([]JSON, len(t))
		for i, elem := range t {
			(*array)[i] = normalize(elem)
		}
		return array
	}
	return nil
}

// MarshalJSON implements RawMessage.
func (u jsonLiteral) MarshalJSON() ([]byte, error) {
	switch u {
	case 0x00:
		return []byte("null"), nil
	case 0x01:
		return []byte("true"), nil
	default:
		return []byte("false"), nil
	}
}
