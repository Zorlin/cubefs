// Copyright 2018 The CubeFS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package master

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cubefs/cubefs/proto"
)

type request struct {
	method string
	path   string
	params map[string]string
	header map[string]string
	body   []byte
	err    error
}

type anyParam struct {
	key string
	val interface{}
}

var ReqHeaderUA = fmt.Sprintf("cubefs-sdk/%v (commit %v)", proto.Version, proto.CommitID)

func (r *request) addParamAny(key string, value interface{}) *request {
	var val string
	switch v := value.(type) {
	case string:
		val = v
	case bool:
		val = strconv.FormatBool(v)

	case int8:
		val = strconv.FormatInt(int64(v), 10)
	case int16:
		val = strconv.FormatInt(int64(v), 10)
	case int32:
		val = strconv.FormatInt(int64(v), 10)
	case int64:
		val = strconv.FormatInt(int64(v), 10)
	case int:
		val = strconv.FormatInt(int64(v), 10)

	case uint8:
		val = strconv.FormatUint(uint64(v), 10)
	case uint16:
		val = strconv.FormatUint(uint64(v), 10)
	case uint32:
		val = strconv.FormatUint(uint64(v), 10)
	case uint64:
		val = strconv.FormatUint(uint64(v), 10)
	case uint:
		val = strconv.FormatUint(uint64(v), 10)

	case float32:
		val = strconv.FormatFloat(float64(v), 'f', 6, 64)
	case float64:
		val = strconv.FormatFloat(float64(v), 'f', 6, 64)

	default:
		r.err = fmt.Errorf("unknown param type %v", value)
		return r
	}
	r.params[key] = val
	return r
}

func (r *request) addParam(key, value string) *request {
	r.params[key] = value
	return r
}

func (r *request) addHeader(key, value string) *request {
	r.header[key] = value
	return r
}

func (r *request) setBody(body []byte) *request {
	r.body = body
	return r
}

func (r *request) Param(params ...anyParam) *request {
	for _, param := range params {
		r.addParamAny(param.key, param.val)
	}
	return r
}

func (r *request) Header(headers map[string]string, added ...string) *request {
	if len(added)%2 == 1 {
		added = added[:len(added)-1]
	}
	for k, v := range headers {
		r.header[k] = v
	}
	for idx := 0; idx < len(added); idx += 2 {
		r.header[added[idx]] = added[idx+1]
	}
	return r
}

func (r *request) Body(body interface{}) *request {
	reqBody, ok := body.([]byte)
	if !ok {
		var err error
		if reqBody, err = json.Marshal(body); err != nil {
			r.err = fmt.Errorf("body json marshal %s", err.Error())
			return r
		}
	}
	r.body = reqBody
	return r
}

func newRequest(method string, path string) *request {
	req := &request{
		method: method,
		path:   path,
		params: make(map[string]string),
		header: make(map[string]string),
	}
	req.header["User-Agent"] = ReqHeaderUA
	return req
}

func mergeHeader(headers map[string]string, added ...string) map[string]string {
	if len(added)%2 == 1 {
		added = added[:len(added)-1]
	}
	copied := make(map[string]string, len(headers)+len(added)/2)
	for k, v := range headers {
		copied[k] = v
	}
	for idx := 0; idx < len(added); idx += 2 {
		copied[added[idx]] = added[idx+1]
	}
	return copied
}
