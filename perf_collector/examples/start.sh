#!/bin/bash

set -e
ROOT=$(unset CDPATH && cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
cd $ROOT

exec $ROOT/perf_collector \
    --listen-port=":1121" \
    --config-path="/var/lib/perf_collector/config.yml" \
    --refresh-interval="60" \
    --worker-threads="5"