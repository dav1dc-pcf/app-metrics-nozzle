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
package api

import (
	"app-metrics-nozzle/domain"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
)

var logger = log.New(os.Stdout, "", 0)
var Client CFClientCaller

type CFClientCaller interface {
	AppByGuid(guid string) (cfclient.App, error)
	GetAppInstances(guid string) (map[string]cfclient.AppInstance, error)
	UsersBy(guid string, entity string) ([]cfclient.User, error)
	ListSpaces() ([]cfclient.Space, error)
	ListOrgs() ([]cfclient.Org, error)
	ListApps() ([]cfclient.App, error)
	AppSpace(app cfclient.App) (cfclient.Space, error)
	SpaceOrg(space cfclient.Space) (cfclient.Org, error)
}

func AppByGuidVerify(guid string) cfclient.App {
	app, _ := Client.AppByGuid(guid)
	return app
}

func AppInstancesByGuidVerify(guid string) map[string]cfclient.AppInstance {
	app, _ := Client.GetAppInstances(guid)
	return app
}

func UsersByOrgVerify(guid string) []cfclient.User {
	app, _ := Client.UsersBy(guid, "organizations")
	return app
}

func UsersBySpaceVerify(guid string) []cfclient.User {
	app, _ := Client.UsersBy(guid, "spaces")
	return app
}

func AnnotateWithCloudControllerData(app *domain.App) {

	ccAppDetails, _ := Client.AppByGuid(app.GUID)

	instances, _ := Client.GetAppInstances(app.GUID)
	runnintCount := 0
	instanceUp := "RUNNING"

	space, _ := Client.AppSpace(ccAppDetails)
	org, _ := Client.SpaceOrg(space)

	app.Diego = ccAppDetails.Diego
	app.Buildpack = ccAppDetails.Buildpack
	app.Instances = make([]domain.Instances, int64(len(instances)))

	for idx, eachInstance := range instances {
		if strings.Compare(instanceUp, eachInstance.State) == 0 {
			runnintCount++
		}
		i, _ := strconv.ParseInt(idx, 10, 32)
		app.Instances[i].InstanceIndex = i
		app.Instances[i].State = eachInstance.State
		app.Instances[i].Since = eachInstance.Since
		app.Instances[i].Uptime = eachInstance.Uptime
	}

	if len(app.Buildpack) == 0 {
		app.Buildpack = ccAppDetails.DetectedBP
	}

	// Omit ENVIRONMENT information, as this often leaks sensitive information about the applications runtime
	//app.Environment = ccAppDetails.Environment
	app.Environment = nil

	app.Organization.ID = org.Guid
	app.Organization.Name = org.Name

	app.Space.ID = space.Guid
	app.Space.Name = space.Name

	app.InstanceCount.Configured = len(instances)
	app.InstanceCount.Running = runnintCount

	app.EnvironmentSummary.TotalDiskConfigured = ccAppDetails.DiskQuota * 1024 * 1024
	app.EnvironmentSummary.TotalMemoryConfigured = ccAppDetails.MemQuota * 1024 * 1024

	app.EnvironmentSummary.TotalDiskProvisioned = ccAppDetails.DiskQuota * 1024 * 1024 * int32(len(instances))
	app.EnvironmentSummary.TotalMemoryProvisioned = ccAppDetails.MemQuota * 1024 * 1024 * int32(len(instances))

	now := time.Now()
	app.StatsSince = now.Unix()

	if 0 < len(ccAppDetails.RouteData) {
		app.Routes = make([]string, len(ccAppDetails.RouteData))
		for i := 0; i < len(ccAppDetails.RouteData); i++ {
			app.Routes[i] = ccAppDetails.RouteData[i].Entity.Host + "." + ccAppDetails.RouteData[i].Entity.DomainData.Entity.Name
		}
	}

	app.State = ccAppDetails.State
}

func UsersForSpace(guid string) (Users []cfclient.User) {
	users, _ := Client.UsersBy(guid, "spaces")
	return users
}

func UsersForOrganization(guid string) (Users []cfclient.User) {
	users, _ := Client.UsersBy(guid, "organizations")
	return users
}

func SpacesDetailsFromCloudController() (Spaces []cfclient.Space) {
	spaces, _ := Client.ListSpaces()
	return spaces
}

func OrgsDetailsFromCloudController() (Orgs []cfclient.Org) {
	orgs, _ := Client.ListOrgs()
	return orgs
}
