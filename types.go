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

	"appengine"
)

type URLRequest struct {
	URL string `json:"url"`
}

type route struct {
	pattern *regexp.Regexp
	handler http.Handler
}

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
	c := appengine.NewContext(r)
	c.Infof("[ServeHTTP] %v, N routers: %v", r.URL.Path, len(rh.routes))
	for _, route := range rh.routes {
		c.Infof("[ServeHTTP] matched %v", route.pattern)
		if route.pattern.MatchString(r.URL.Path) {
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	http.NotFound(w, r)
}
