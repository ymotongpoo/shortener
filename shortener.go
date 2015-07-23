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
	"math/rand"
	"net/http"
)

const Version = "0.1.0"

var (
	seed    = int64(25) // seed for randomness
	samples = []byte{}
)

func initSample() {
	original := make([]byte, 64) // [0-9A-Za-z_\-] contains 64 chars.
	for i := 0; i < 10; i++ {    // 0-9
		original[i] = byte(48 + i)
	}
	for i := 0; i < 26; i++ { // A-Z
		original[10+i] = byte(65 + i)
	}
	for i := 0; i < 26; i++ { // a-z
		original[36+i] = byte(97 + i)
	}
	original[62] = '_'
	original[63] = '-'
	rand.Seed(seed)
	for i, n := range rand.Perm(len(original)) {
		samples[i] = original[n]
	}
}

func init() {
	initSample()
	router := &RegexpHandler{}
	router.HandleFunc(`/`, top)
	router.HandleFunc(`/[0-9A-Za-z_\-]{6}`, redirect)
	router.HandleFunc(`/version`, version)
	router.HandleFunc(`/shortener/v1`, shortner)
	http.Handle("/", router)
}

func version(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, Version)
}

func top(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello")
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

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://www.google.com", http.StatusFound)
}
