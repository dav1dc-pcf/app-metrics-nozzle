# App Metrics Nozzle *(for Vizceral App Metrics UI)*

This project was forked from [pivotalservices/app-metrics-nozzle](https://github.com/pivotalservices/app-metrics-nozzle).

The original **README** can be viewed [here](README.orig.md).

## Future Improvements

1. This project uses a very old version of the **firehose2syslog** library, and unfortunately it contains a bug that causes it to crash under high message load from the firehose with the error `ERR fatal error: concurrent map read and map write`. Updating to the latest f2s library is non-trivial. We're planning to address this problem in the next iteration of this project, likely by re-writing the core.

2. The current behaviour of this project drops all accumulated statistical data when the list of applications is refreshed from the Cloud Controller API (resulting in a drop out of the data). We would like to change this behaviour to use a Ring Buffer pattern so that a rolling history of stats can be maintained between refreshes of the Application list.

*As always, PR's welcome!*
