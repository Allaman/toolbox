# Prometheus Pushgateway

How to push metrics to [Prometheus Pushgateway](https://github.com/prometheus/pushgateway)

## Go version

Requires a Pushgateway running on `localhost:9091`.

## Sh version

Env variables:

-   `pUSHGATEWAY` url of the Pushgateway
-   `JOB` the name of the job
-   `LABELS` a label string
-   `HELP_MESSAGE` the help message
-   `METRIC_NAME` name of the metric
-   `METRIC_TYPE` type of the metric
-   `METRIC_VALUE` value of the metric

Refer to the Github repository for more details about the variables.
