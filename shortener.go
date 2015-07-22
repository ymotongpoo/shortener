// Copyright 2015 Yoshi Yamaguchi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shortener

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const Version = "0.1.0"

func init() {
	http.HandleFunc("/shortener/v1", shortner)
	http.HandleFunc("/version", version)
}

func version(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, Version)
}

func shortner(w http.ResponseWriter, r *http.Request) {
	req := URLRequest{}
	if r.Method != "POST" {
		http.Error(w, fmt.Sprintf("Methods but for POST are not allowed"), http.StatusMethodNotAllowed)
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unexpected payload: %v", err), http.StatusInternalServerError)
	}
	err = json.Unmarshal(data, &req)
	if err != nil {
		http.Error(w, fmt.Sprintf("JSON decode error: %v", err), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, req.URL)
}
