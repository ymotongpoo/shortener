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
	"net/url"
	"path"
	"time"

	"appengine"
	"appengine/datastore"
)

const (
	Version            = "0.1.0"
	starttime          = int64(1437625938157481) // around 2015 Jul 23 13:34
	timestampLeftShift = 22
)

var chars []byte // used for unique id generation.

// initChars initialize chars as sequence of 0-9A-Za-z
func initChars() {
	chars = make([]byte, 62)
	for i := 0; i < 10; i++ { // 0-9
		chars[i] = byte(48 + i)
	}
	for i := 0; i < 26; i++ { // A-Z
		chars[i+10] = byte(65 + i)
	}
	for i := 0; i < 26; i++ { // a-z
		chars[i+36] = byte(97 + i)
	}
}

// init setup chars and URL routers.
func init() {
	initChars()
	router := &RegexpHandler{}
	router.HandleFunc(`/`, top)
	router.HandleFunc(`/[0-9A-Za-z_\-]{10,}`, redirect)
	router.HandleFunc(`/version`, version)
	router.HandleFunc(`/shortener/v1`, shortener)
	http.Handle("/", router)
}

// version returns application version.
func version(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, Version)
}

// top returns UI front page.
func top(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "hello")
}

// shortener
func shortener(w http.ResponseWriter, r *http.Request) {
	req := URLRequest{}
	if r.Method != "POST" {
		http.Error(w, fmt.Sprintf("Methods but for POST are not allowed"), http.StatusMethodNotAllowed)
		return
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unexpected payload: %v", err), http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(data, &req)
	if err != nil {
		http.Error(w, fmt.Sprintf("JSON decode error: %v", err), http.StatusInternalServerError)
		return
	}
	url, err := url.Parse(req.URL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid URL: %v", req.URL), http.StatusInternalServerError)
		return
	}
	switch url.Scheme {
	case "https", "http", "ftp":
		break
	default:
		http.Error(w, fmt.Sprintf("Scheme is not supported: %v", req.URL), http.StatusInternalServerError)
		return
	}

	c := appengine.NewContext(r)
	id := uniqueid()
	e := &URLEntity{
		ID:  id,
		URL: req.URL,
	}
	key := datastore.NewIncompleteKey(c, "URL", nil)
	_, err = datastore.Put(c, key, e)
	if err != nil {
		http.Error(w, fmt.Sprintf("datastore put error: %v", err), http.StatusInternalServerError)
		return
	}
	entry, err := json.Marshal(e)
	if err != nil {
		http.Error(w, fmt.Sprintf("JSON encode error: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%v", string(entry))
}

// redirect find specified shortened URL path from datastore and redirect to original URL.
func redirect(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	id := path.Base(r.URL.Path)
	es := []URLEntity{}
	keys, err := datastore.NewQuery("URL").Filter("ID=", id).GetAll(c, &es)
	if err != nil {
		http.Error(w, fmt.Sprintf("datastore get error: %v", err), http.StatusInternalServerError)
		return
	}
	if len(es) == 0 {
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}
	original := es[0].URL
	es[0].Count += 1
	_, err = datastore.Put(c, keys[0], &es[0])
	if err != nil {
		http.Error(w, fmt.Sprintf("datastore put error: %v", err), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, original, http.StatusFound)
}

// uniqueid generates unique id from unix time in microsecond and randomnumber based on it,
// then convert it to base 62 number
func uniqueid() string {
	now := time.Now().UnixNano() / 1000
	delta := now - starttime
	rand.Seed(delta)
	n := rand.Intn(2 ^ timestampLeftShift)
	id := delta<<timestampLeftShift | int64(n)

	size := int64(len(chars))
	result := make([]byte, 36)
	i := 0
	for id > 0 {
		rem := id % size
		id = id / size
		result[i] = chars[rem]
		i++
	}
	return string(result[:i])
}
