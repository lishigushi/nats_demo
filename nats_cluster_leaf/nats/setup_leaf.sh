#!/usr/bin/env bash

set -e

nats --server nats://localhost:4222 stream add \
     --js-domain leaf teststream \
     --mirror teststream \
     --storage=file \
     --replicas=1 \
     --retention=limits \
     --discard=old \
     --max-msgs=-1 \
     --max-msgs-per-subject=-1 \
     --max-bytes=-1 \
     --max-age=-1 \
     --max-msg-size=-1 \
     --dupe-window=2m0s \
     --no-allow-rollup \
     --no-deny-delete \
     --no-deny-purge
