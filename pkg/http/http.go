/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.
 * Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var (
	cli *http.Client
)

func init() {
	cli = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

// Get http client get function
func Get(url string, result ...interface{}) (err error) {
	return do(http.MethodGet, url, nil, result...)
}

// Post http client post function
func Post(url string, data interface{}, result ...interface{}) (err error) {
	return do(http.MethodPost, url, data, result...)
}

// Put http client put function
func Put(url string, data interface{}, result ...interface{}) (err error) {
	return do(http.MethodPut, url, data, result...)
}

// Delete http client delete function
func Delete(url string, result ...interface{}) (err error) {
	return do(http.MethodDelete, url, nil, result...)
}

func do(method, url string, data interface{}, result ...interface{}) (err error) {
	var byt []byte
	if data != nil {
		byt, err = json.Marshal(data)
		if err != nil {
			return
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(byt))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error status code: %d", resp.StatusCode)
	}

	if len(result) == 0 {
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if b, ok1 := result[0].(*[]byte); ok1 {
		*b = body
		return
	} else if str, ok2 := result[0].(*string); ok2 {
		*str = string(body)
		return
	}

	return json.Unmarshal(body, result[0])
}

// New http client
func New() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}
