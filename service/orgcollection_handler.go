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
	"net/http"

	"app-metrics-nozzle/usageevents"
	"github.com/unrolled/render"

	"strings"

	"github.com/gorilla/mux"
)

func spaceDetailsHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", req.Header.Get("Origin"))
		w.Header().Add("Access-Control-Allow-Methods", "GET")

		vars := mux.Vars(req)
		space := vars["space"]

		found := false
		for idx := range usageevents.Spaces {
			if 0 == strings.Compare(space, usageevents.Spaces[idx].Name) {
				found = true
				formatter.JSON(w, http.StatusOK, usageevents.Spaces[idx])
			}
		}
		if !found {
			formatter.JSON(w, http.StatusNotFound, "Space not found.")
		}
	}
}

func spaceHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", req.Header.Get("Origin"))
		w.Header().Add("Access-Control-Allow-Methods", "GET")

		if 0 < len(usageevents.Spaces) {
			formatter.JSON(w, http.StatusOK, usageevents.Spaces)
		} else {
			formatter.JSON(w, http.StatusNotFound, "No spaces found.")
		}

	}
}

func orgDetailsHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", req.Header.Get("Origin"))
		w.Header().Add("Access-Control-Allow-Methods", "GET")

		vars := mux.Vars(req)
		org := vars["org"]

		found := false
		for idx := range usageevents.Orgs {
			if 0 == strings.Compare(org, usageevents.Orgs[idx].Name) {
				found = true
				formatter.JSON(w, http.StatusOK, usageevents.Orgs[idx])
			}
		}
		if !found {
			formatter.JSON(w, http.StatusNotFound, "Org not found.")
		}
	}
}

func orgsHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", req.Header.Get("Origin"))
		w.Header().Add("Access-Control-Allow-Methods", "GET")

		if 0 < len(usageevents.Orgs) {
			formatter.JSON(w, http.StatusOK, usageevents.Orgs)
		} else {
			formatter.JSON(w, http.StatusNotFound, "No organizations found.")
		}

	}
}
