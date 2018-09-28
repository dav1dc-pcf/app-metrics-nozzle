# App Metrics Nozzle *(for Vizceral App Metrics UI)*

This project was forked from [pivotalservices/app-metrics-nozzle](https://github.com/pivotalservices/app-metrics-nozzle).

The original **README** can be viewed [here](README.orig.md).

## Setup

Each Foundation that one wishes to visualize will need an instance of this application pushed to it.
Then your instance of [vz-metrics-adapter](https://github.com/dav1dc-pcf/vz-metrics-adapter) will need to primed (configured) with knowledge about the location of each instance of **app-metrics-nozzle**.

Edit `manifest.yml` and un-comment the following lines:

```
#    API_ENDPOINT: https://api.local.pcfdev.io
#    DOPPLER_ENDPOINT: wss://doppler.local.pcfdev.io:443
#    FIREHOSE_USER: xxxxx
#    FIREHOSE_PASSWORD: xxxxx
```

Next, set values that are appropirate for your environment:

```
    API_ENDPOINT: https://api.sys.your-pcf.com
    DOPPLER_ENDPOINT: wss://doppler.sys.your-pcf.com:443
    FIREHOSE_USER: username
    FIREHOSE_PASSWORD: password
```

Now push the application to your PCF/CF.

If one requires more setup guidance, please refer to the original [README.md](README.orig.md).

## Notes

The **username** & **password** provided requires authorization to both the firehose and the Cloud Controller API *(as the same credentials are used to access both services).*

The environment information has been removed from the output of **app-metrics-nozzle** *(So that credentials passed into the application via the environment are not leaked to the Internet via this application instance).*

## Future Improvements

1. This project uses a very old version of the **firehose2syslog** library, and unfortunately it contains a bug that causes it to crash under high message load from the firehose with the error `ERR fatal error: concurrent map read and map write`. Updating to the latest f2s library is non-trivial. We're planning to address this problem in the next iteration of this project, likely by re-writing the core.

2. The current behaviour of this project drops all accumulated statistical data when the list of applications is refreshed from the Cloud Controller API (resulting in a drop out of the data). We would like to change this behaviour to use a Ring Buffer pattern so that a rolling history of stats can be maintained between refreshes of the Application list.

*As always, PR's welcome!*
