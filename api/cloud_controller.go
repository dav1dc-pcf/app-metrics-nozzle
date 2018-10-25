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
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"app-metrics-nozzle/domain"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
)

type apiClient struct {
	client CFClientCaller
}

var logger = log.New(os.Stdout, "", 0)

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

// NewApiClient returns wrapper for CC Client
func NewApiClient(client CFClientCaller) apiClient {
	return apiClient{client}
}

func (a apiClient) AppByGuidVerify(guid string) cfclient.App {
	app, _ := a.client.AppByGuid(guid)
	return app
}

func (a apiClient) AppInstancesByGuidVerify(guid string) map[string]cfclient.AppInstance {
	app, _ := a.client.GetAppInstances(guid)
	return app
}

func (a apiClient) UsersByOrgVerify(guid string) []cfclient.User {
	app, _ := a.client.UsersBy(guid, "organizations")
	return app
}

func (a apiClient) UsersBySpaceVerify(guid string) []cfclient.User {
	app, _ := a.client.UsersBy(guid, "spaces")
	return app
}

func (a apiClient) AnnotateWithCloudControllerData(app *domain.App) {

	ccAppDetails, _ := a.client.AppByGuid(app.GUID)

	instances, _ := a.client.GetAppInstances(app.GUID)
	runnintCount := 0
	instanceUp := "RUNNING"

	space, _ := a.client.AppSpace(ccAppDetails)
	org, _ := a.client.SpaceOrg(space)

	app.Diego = ccAppDetails.Diego
	app.Buildpack = ccAppDetails.Buildpack
	app.Instances = make([]domain.Instance, int64(len(instances)))

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

func (a apiClient) UsersForSpace(guid string) (Users []cfclient.User) {
	users, _ := a.client.UsersBy(guid, "spaces")
	return users
}

func (a apiClient) UsersForOrganization(guid string) (Users []cfclient.User) {
	users, _ := a.client.UsersBy(guid, "organizations")
	return users
}

func (a apiClient) SpacesDetailsFromCloudController() (Spaces []cfclient.Space) {
	spaces, _ := a.client.ListSpaces()
	return spaces
}

func (a apiClient) OrgsDetailsFromCloudController() (Orgs []cfclient.Org) {
	orgs, _ := a.client.ListOrgs()
	return orgs
}
