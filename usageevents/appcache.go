package usageevents

import "github.com/cloudfoundry-community/firehose-to-syslog/caching"

// AppCache allows retrieving information from the cache
type AppCache struct{}

// CachedApp interface for retrieving information from the cache
type CachedApp interface {
	GetAppByGUID(appGUID string) []caching.App
	GetAppInfo(appGUID string) caching.App
	GetAllApp() []caching.App
}

// GetAppByGUID get Applications by guid
func (c *AppCache) GetAppByGUID(appGUID string) []caching.App {
	return caching.GetAppByGuid(appGUID)
}

// GetAppInfo get App by appGUID
func (c *AppCache) GetAppInfo(appGUID string) caching.App {
	return caching.GetAppInfo(appGUID)
}

// GetAllApp get all caching.App
func (c *AppCache) GetAllApp() []caching.App {
	return caching.GetAllApp()
}
