#!/bin/sh

consul agent --config-dir=/consul/config -retry-join myconsul-consul-server.default.svc.cluster.local
