/*
Copyright 2016 Pivotal

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package service

import (
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

// NewServer configures and returns a Server.
func NewServer() *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})

	n := negroni.Classic()
	mx := mux.NewRouter()

	initRoutes(mx, formatter)

	n.UseHandler(mx)
	return n
}

func initRoutes(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/api/apps/{org}/{space}/{app}/{instance_id}", appInstanceHandler(formatter)).Methods("GET")
	mx.HandleFunc("/api/apps/{org}/{space}/{app}", appHandler(formatter)).Methods("GET")
	mx.HandleFunc("/api/apps/{org}/{space}", appSpaceHandler(formatter)).Methods("GET")
	mx.HandleFunc("/api/apps/{org}", appOrgHandler(formatter)).Methods("GET")
	mx.HandleFunc("/api/apps", appAllHandler(formatter)).Methods("GET")
	mx.HandleFunc("/api/orgs/{org}", orgDetailsHandler(formatter)).Methods("GET")
	mx.HandleFunc("/api/orgs", orgsHandler(formatter)).Methods("GET")
	mx.HandleFunc("/api/orgs/{org}/users", orgsUsersHandler(formatter)).Methods("GET")
	mx.HandleFunc("/api/orgs/{org}/{role}/users", orgsUsersByRoleHandler(formatter)).Methods("GET")
	mx.HandleFunc("/api/spaces/{space}", spaceDetailsHandler(formatter)).Methods("GET")
	mx.HandleFunc("/api/spaces/{space}/users", spacesUsersHandler(formatter)).Methods("GET")
	mx.HandleFunc("/api/spaces/{space}/{role}/users", spacesUsersByRoleHandler(formatter)).Methods("GET")
	mx.HandleFunc("/api/spaces", spaceHandler(formatter)).Methods("GET")
}
