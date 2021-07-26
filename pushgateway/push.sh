#!/bin/sh

function pushMetric {
  _checkVars
  echo "`date -Iseconds` Pushing metric to Prometheus Pushgateway"
  cat <<EOF | curl -fsSL -XPUT --data-binary @- "$PUSHGATEWAY/metrics/job/$JOB/$LABELS"
# HELP $HELP_MESSAGE
# TYPE $METRIC_NAME $METRIC_TYPE
$METRIC_NAME $METRIC_VALUE
EOF
}

function _checkVars {
  if [ -z "$PUSHGATEWAY" ]; then
    echo "missing variable `PUSHGATEWAY`"
    exit 1
  fi
  if [ -z "$JOB" ]; then
    echo "missing variable `JOB`"
    exit 1
  fi
  if [ -z "$LABELS" ]; then
    echo "warning: empty `LABELS`"
  fi
  if [ -z "$HELP_MESSAGE" ]; then
    echo "missing variable `HELP_MESSAGE`"
    exit 1
  fi
  if [ -z "$METRIC_NAME" ]; then
    echo "missing variable `METRIC_NAME`"
    exit 1
  fi
  if [ -z "$METRIC_TYPE" ]; then
    echo "missing variable `METRIC_TYPE`"
    exit 1
  fi
  if [ -z "$METRIC_VALUE" ]; then
    echo "missing variable `METRIC_VALUE`"
    exit 1
  fi
}
