package usageevents

import (
	"fmt"
	"time"

	"app-metrics-nozzle/domain"

	"github.com/cloudfoundry-community/firehose-to-syslog/caching"
)

type apiClient interface {
	AnnotateWithCloudControllerData(app *domain.App)
}

// ReloadApps responsilbe for refreshing apps in the cache
func ReloadApps(cachedApps []caching.App, client apiClient) {
	logger.Println("Start filling app/space/org cache.")
	for idx := range cachedApps {

		org := cachedApps[idx].OrgName
		space := cachedApps[idx].SpaceName
		app := cachedApps[idx].Name
		key := GetMapKeyFromAppData(org, space, app)

		appId := cachedApps[idx].Guid
		name := cachedApps[idx].Name

		appDetail := &domain.App{GUID: appId, Name: name}
		client.AnnotateWithCloudControllerData(appDetail)
		appDetail.FetchTime = time.Now().String()

		a, _ := AppDetails.Get(key)
		if a != nil {
			appDetails := a.(domain.App)
			// Do our best to copy over existing Cell IP's for instances
			for idx, eachInstance := range appDetails.Instances {
				if idx < len(appDetail.Instances) {
					if appDetail.Instances[idx].InstanceIndex == eachInstance.InstanceIndex {
						appDetail.Instances[idx].CellIP = eachInstance.CellIP
					}
				}
			}
		}

		AppDetails.Set(key, *appDetail)
		logger.Println(fmt.Sprintf("Registered [%s]", key))
	}

	logger.Println(fmt.Sprintf("Done filling cache! Found [%d] Apps", len(cachedApps)))
}
