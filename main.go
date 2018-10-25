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

package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/cloudfoundry-community/firehose-to-syslog/logging"

	"github.com/CrowdSurge/banner"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/boltdb/bolt"
	goClient "github.com/cloudfoundry-community/go-cfclient"

	"github.com/cloudfoundry-community/firehose-to-syslog/caching"
	"github.com/cloudfoundry/noaa/consumer"
	"github.com/cloudfoundry/sonde-go/events"
	"app-metrics-nozzle/api"
	"app-metrics-nozzle/service"
	"app-metrics-nozzle/usageevents"
)

var (
	debug             = kingpin.Flag("debug", "Enable debug mode. This disables forwarding to syslog").Default("false").OverrideDefaultFromEnvar("DEBUG").Bool()
	apiEndpoint       = kingpin.Flag("api-endpoint", "Api endpoint address. For bosh-lite installation of CF: https://api.10.244.0.34.xip.io").OverrideDefaultFromEnvar("API_ENDPOINT").Required().String()
	dopplerEndpoint   = kingpin.Flag("doppler-endpoint", "Overwrite default doppler endpoint return by /v2/info").OverrideDefaultFromEnvar("DOPPLER_ENDPOINT").String()
	user              = kingpin.Flag("user", "Admin user.").Default("admin").OverrideDefaultFromEnvar("FIREHOSE_USER").String()
	password          = kingpin.Flag("password", "Admin password.").Default("admin").OverrideDefaultFromEnvar("FIREHOSE_PASSWORD").String()
	skipSSLValidation = kingpin.Flag("skip-ssl-validation", "Please don't").Default("false").OverrideDefaultFromEnvar("SKIP_SSL_VALIDATION").Bool()
	boltDatabasePath  = kingpin.Flag("boltdb-path", "Bolt Database path ").Default("my.db").OverrideDefaultFromEnvar("BOLTDB_PATH").String()
	tickerTime        = kingpin.Flag("cc-pull-time", "CloudController Polling time in sec").Default("60s").OverrideDefaultFromEnvar("CF_PULL_TIME").Duration()
)

const (
	version = "0.0.1"
)

type apiClient interface {
	OrgsDetailsFromCloudController() []goClient.Org
	UsersForOrganization(guid string) []goClient.User
	SpacesDetailsFromCloudController() []goClient.Space
	UsersForSpace(guid string) []goClient.User
}

var logger = log.New(os.Stdout, "", 0)

func main() {
	var wg sync.WaitGroup

	banner.Print("metrics usage nozzle")
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	// Start web server
	go func() {
		server := service.NewServer()
		server.Run(":" + port)
	}()

	kingpin.Version(version)
	kingpin.Parse()

	logger.Println(fmt.Sprintf("Starting app-metrics-nozzle %s ", version))

	c := goClient.Config{
		ApiAddress:        *apiEndpoint,
		Username:          *user,
		Password:          *password,
		SkipSslValidation: *skipSSLValidation,
	}
	cfClient, err := goClient.NewClient(&c)

	if err != nil {
		logger.Fatal("Error connecting to CF API", err)
	}

	if len(*dopplerEndpoint) > 0 {
		cfClient.Endpoint.DopplerEndpoint = *dopplerEndpoint
	}
	logger.Println(fmt.Sprintf("Using %s as doppler endpoint", cfClient.Endpoint.DopplerEndpoint))

	//Use bolt for in-memory  - file caching
	db, err := bolt.Open(*boltDatabasePath, 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		logger.Fatal("Error opening bolt db: ", err)
		os.Exit(1)

	}

	defer db.Close()

	caching.SetCfClient(cfClient)
	caching.SetAppDb(db)
	caching.CreateBucket()

	apiClient := api.NewApiClient(cfClient)

	//Let's Update the database the first time
	usageevents.ReloadApps(caching.GetAllApp(), apiClient)
	reloadEnvDetails(apiClient)
	lastReloaded := time.Now()
	fmt.Println("Reloaded first time:", lastReloaded)

	// Ticker Polling the CC every X sec
	ccPolling := time.NewTicker(*tickerTime)

	go func() {
		for range ccPolling.C {
			now := time.Now()
			logger.Print(" ---> " + now.Format(time.RFC3339))
			usageevents.ReloadApps(caching.GetAllApp(), apiClient)
			reloadEnvDetails(apiClient)
		}
	}()

	token, _ := cfClient.GetToken()

	for _, application := range usageevents.AppDbCache.GetAllApp() {
		firehose := createFirehoseChan(cfClient.Endpoint.DopplerEndpoint, token, application.Guid, *skipSSLValidation, consumer.KeepAlive)
		if firehose != nil {
			wg.Add(1)
			go func() {
				usageevents.ProcessEvents(firehose)
				logger.Fatal("Lost connection to Firehose...Please check settings and try again!")
			}()
			logger.Println(fmt.Sprintf("Firehose Subscription Succesfull for %s! Routing events...", application.Guid))
		} else {
			logger.Fatal("Failed connecting to Firehose...Please check settings and try again!")
		}

	}
	wg.Wait()
}

func reloadEnvDetails(client apiClient) {
	usageevents.Orgs = client.OrgsDetailsFromCloudController()

	for idx := range usageevents.Orgs {
		users := client.UsersForOrganization(usageevents.Orgs[idx].Guid)
		usageevents.OrganizationUsers[usageevents.Orgs[idx].Name] = users
	}

	usageevents.Spaces = client.SpacesDetailsFromCloudController()

	for idx := range usageevents.Spaces {
		users := client.UsersForSpace(usageevents.Spaces[idx].Guid)
		usageevents.SpacesUsers[usageevents.Spaces[idx].Name] = users
	}
}

func createFirehoseChan(dopplerEndpoint, token, appID string, skipSSLValidation bool, keepAlive time.Duration) <-chan *events.Envelope {
	consumer.KeepAlive = keepAlive
	connection := consumer.New(dopplerEndpoint, &tls.Config{InsecureSkipVerify: skipSSLValidation}, nil)
	connection.SetDebugPrinter(ConsoleDebugPrinter{})
	msgChan, errorChan := connection.Stream(appID, token)
	go func() {
		for err := range errorChan {
			logging.LogError("Firehose Error!", err.Error())
		}
	}()
	return msgChan
}

// ConsoleDebugPrinter used for the firehoseconnection
type ConsoleDebugPrinter struct{}

// Print function for the ConsoleDebugPrinter for the firehoseconnection
func (c ConsoleDebugPrinter) Print(title, dump string) {
	logging.LogStd(title, false)
	logging.LogStd(dump, false)
}
