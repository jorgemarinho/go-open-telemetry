#!/bin/sh
chown -R nobody:nogroup /prometheus
exec prometheus --config.file=/etc/prometheus/prometheus.yml --storage.tsdb.path=/prometheus
