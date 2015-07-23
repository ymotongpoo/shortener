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
	"net/http"
	"regexp"
)

// URLRequest is a struct to request URL shortening to the server.
type URLRequest struct {
	URL string `json:"url"`
}

// route holds URL path regexp pattern and corresponding handler.
type route struct {
	pattern *regexp.Regexp
	handler http.Handler
}

// RegexpHandler holds a set of routing pattern.
type RegexpHandler struct {
	routes []*route
}

func (rh *RegexpHandler) Handler(pattern string, handler http.Handler) {
	p := regexp.MustCompile(`^` + pattern + `$`)
	rh.routes = append(rh.routes, &route{p, handler})
}

func (rh *RegexpHandler) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	p := regexp.MustCompile(`^` + pattern + `$`)
	rh.routes = append(rh.routes, &route{p, http.HandlerFunc(handler)})
}

func (rh *RegexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range rh.routes {
		if route.pattern.MatchString(r.URL.Path) {
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	http.NotFound(w, r)
}
