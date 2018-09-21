package usageevents

import (
	"app-metrics-nozzle/api"
	"app-metrics-nozzle/domain"
	"fmt"

	"github.com/cloudfoundry-community/firehose-to-syslog/caching"
)

func ReloadApps(cachedApps []caching.App) {
	logger.Println("Start filling app/space/org cache.")
	for idx := range cachedApps {

		org := cachedApps[idx].OrgName
		space := cachedApps[idx].SpaceName
		app := cachedApps[idx].Name
		key := GetMapKeyFromAppData(org, space, app)

		appId := cachedApps[idx].Guid
		name := cachedApps[idx].Name

		appDetail := domain.App{GUID: appId, Name: name}
		api.AnnotateWithCloudControllerData(&appDetail)

		// Do our best to copy over existing Cell IP's for instances
		mutex.RLock()
		for idx, eachInstance := range AppDetails[key].Instances {
			if eachInstance.InstanceIndex == appDetail.Instances[idx].InstanceIndex {
				appDetail.Instances[idx].CellIP = eachInstance.CellIP
			}
		}
		mutex.RUnlock()

		mutex.Lock()
		AppDetails[key] = appDetail
		mutex.Unlock()
		logger.Println(fmt.Sprintf("Registered [%s]", key))
	}

	logger.Println(fmt.Sprintf("Done filling cache! Found [%d] Apps", len(cachedApps)))
}
