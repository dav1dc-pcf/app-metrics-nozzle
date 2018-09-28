# App Metrics Nozzle

This is a nozzle for the Cloud Foundry firehose component. It will ingest router events for every application it can detect and use the timestamps on those events to compute usage metrics. The first usage for this nozzle will be to determine if an application is unused by tracking the last time a request was routed to it.
It will also make REST api calls to Cloud Controller and cache application, and user data for the interval specified in configuration.

## REST API
This application exposes a RESTful API that allows consumers to query application usage statistics. This API is described in the following table:

| Resource        | Method           | Description  |
| --- | --- | --- |
| `/api/apps` | GET | Queries the list of all _deployed_ applications. This is all applications in all organizations, in all spaces. The results will be a map whose key is a string of the format `[org]/[space]/[app name]`. |
| `/api/apps/[org]/[space]/[app]` | GET | Obtains application detail information, including time-based usage statistics as of the time of request, including elapsed time (in seconds) since the last event was received, and the requests per second for the app. |
| `/api/apps/[org]/[space]/[app]/[instance_id]` | GET | Obtains information for the instance of the application including IP, CPI usage and Memory usage. |
| `/api/apps/[org]/[space]` | GET | Obtains application details deployed in specified space. |
| `/api/apps/[org]` | GET | Obtains application details deployed in specified organization. |
| `/api/orgs` | GET | Obtains names and guids of all organizations. |
| `/api/orgs/[org]` | GET | Obtains name and guid of an organization. |
| `/api/orgs/[org]/users` | GET | Returns information about users of specified organization. |
| `/api/orgs/[org]/[role]/users` | GET | Returns users of specified organization by role. |
| `/api/spaces` | GET | Returns a list of spaces. |
| `/api/spaces/[space]` | GET | Returns space details. |
| `/api/spaces/[space]/users` | GET | Returns users of specified space. |
| `/api/spaces/[space]/[role]/users` | GET | Returns users of specified space by role. |

### JSON Payloads
This is a sample of what the JSON response looks like for the app `/api/apps`:

```javascript

{
  "buildpack": "java-buildpack=v3.6-offline-https://github.com/cloudfoundry/java-buildpack.git#5194155 java-main open-jdk-like-jre=1.8.0_71 open-jdk-like-memory-calculator=2.0.1_RELEASE spring-auto-reconfiguration=1.10.0_RELEASE",
  "diego": true,
  "environment": {},
  "environment_summary": {
    "total_cpu": 0.19134577882529913,
    "total_disk_configured": 1073741824,
    "total_disk_provisioned": 1073741824,
    "total_disk_usage": 162381824,
    "total_memory_configured": 1073741824,
    "total_memory_provisioned": 1073741824,
    "total_memory_usage": 744574976
  },
  "guid": "bb7b3c89-0a7f-47f7-9dd3-5e4fbd8ded6c",
  "instance_count": {
    "configured": 1,
    "running": 1
  },
  "instances": [
    {
      "cell_ip": "10.65.201.46",
      "cpu_usage": 0.19134577882529913,
      "disk_usage": 162381824,
      "gc_stats": "",
      "index": 0,
      "memory_usage": 744574976,
      "uptime": 16876,
      "since": 1465999843,
      "state": "RUNNING"
    }
  ],
  "name": "cd-demo-music",
  "organization": {
    "id": "c661e8c6-649a-4fe0-b471-afe5982e4e53",
    "name": "Pivotal"
  },
  "event_count": 55,
  "last_event_time": 1466016785654174635,
  "requests_per_second": 0.025380710659898477,
  "elapsed_since_last_event": 0,
  "routes": [
    "ashumilov.cfapps.haas-41.pez.pivotal.io"
  ],
  "space": {
    "id": "dc4d1d1f-f4b9-4c60-8cbb-5763491d00c1",
    "name": "ashumilov"
  },
  "state": "STARTED"
}
,
"org/space/app" : {},
```

If the `last_event_time` field is `0` that indicates that no _router_ events for that application have been discovered _since the nozzle was started_.

## Installation
Run glide install to pull dependencies into vendor directory.
To install this application, it should be run as an app within Cloud Foundry. So, the first thing you'll need to do is push the app. There is a `manifest.yml` already included in the project, so you can just do:

```
cf push app-metrics-nozzle --no-start
```

The `no-start` is important because we have not yet defined the environment variables that allow the application to connect to the Firehose and begin monitoring router requests. We want to end up with a set of environment variables that looks like this when we issue a `cf env app-metrics-nozzle` command:

```
User-Provided:
API_ENDPOINT: https://api.local.pcfdev.io
DOPPLER_ENDPOINT: wss://doppler.local.pcfdev.io:443
CF_PULL_TIME: 9999s
FIREHOSE_PASSWORD: (this is a secret)
FIREHOSE_SUBSCRIPTION_ID: app-metrics-nozzle
FIREHOSE_USER: (this is also secret)
SKIP_SSL_VALIDATION: true
```
Once you've set these environment variables with `cf set-env (app) (var) (value)` you can just start the application usage nozzle via `cf start`. Make sure the application has come up by hitting the API endpoint. Depending on how large of a foundation in which it was deployed, it can take _several minutes_ for the cache of application metadata to fill up.

DOPPLER_ENDPOINT can be obtained by running
```bash
cf curl /v2/info
```

