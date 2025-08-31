#!/bin/sh
chmod -R o+r /var/lib/docker/containers
find /var/lib/docker/containers -type d -exec chmod o+x {} \;
exec filebeat
