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

	cfclient "github.com/cloudfoundry-community/go-cfclient"

	"github.com/gorilla/mux"
	"app-metrics-nozzle/usageevents"
	"github.com/unrolled/render"
)

func spacesUsersHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", req.Header.Get("Origin"))
		w.Header().Add("Access-Control-Allow-Methods", "GET")

		vars := mux.Vars(req)
		space := vars["space"]

		if 0 < len(usageevents.Spaces) {
			formatter.JSON(w, http.StatusOK, usageevents.SpacesUsers[space])
		} else {
			formatter.JSON(w, http.StatusNotFound, "No spaces found.")
		}
	}
}

func spacesUsersByRoleHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", req.Header.Get("Origin"))
		w.Header().Add("Access-Control-Allow-Methods", "GET")

		vars := mux.Vars(req)
		space := vars["space"]
		role := vars["role"]

		if 0 < len(usageevents.SpacesUsers[space]) {
			foundUsers := make([]cfclient.User, 0)
			for idx := range usageevents.SpacesUsers[space] {
				if stringInSlice(role, usageevents.SpacesUsers[space][idx].SpaceRoles) {
					foundUsers = append(foundUsers, usageevents.SpacesUsers[space][idx])
				}
			}
			if 0 < len(foundUsers) {
				formatter.JSON(w, http.StatusOK, foundUsers)
			} else {
				formatter.JSON(w, http.StatusNotFound, "No users found for specified role.")
			}
		} else {
			formatter.JSON(w, http.StatusNotFound, "No spaces found.")
		}
	}
}

func orgsUsersByRoleHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", req.Header.Get("Origin"))
		w.Header().Add("Access-Control-Allow-Methods", "GET")

		vars := mux.Vars(req)
		org := vars["org"]
		role := vars["role"]

		if 0 < len(usageevents.OrganizationUsers[org]) {
			foundUsers := make([]cfclient.User, 0)
			for idx := range usageevents.OrganizationUsers[org] {
				if stringInSlice(role, usageevents.OrganizationUsers[org][idx].OrganizationRoles) {
					foundUsers = append(foundUsers, usageevents.OrganizationUsers[org][idx])
				}
			}
			if 0 < len(foundUsers) {
				formatter.JSON(w, http.StatusOK, foundUsers)
			} else {
				formatter.JSON(w, http.StatusNotFound, "No users found for specified role.")
			}
		} else {
			formatter.JSON(w, http.StatusNotFound, "No organizations found.")
		}
	}
}

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func orgsUsersHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", req.Header.Get("Origin"))
		w.Header().Add("Access-Control-Allow-Methods", "GET")

		vars := mux.Vars(req)
		org := vars["org"]

		if 0 < len(usageevents.Orgs) {
			formatter.JSON(w, http.StatusOK, usageevents.OrganizationUsers[org])
		} else {
			formatter.JSON(w, http.StatusNotFound, "No organizations found.")
		}
	}
}
