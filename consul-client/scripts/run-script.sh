#!/bin/sh

#consul agent --config-dir=/consul/config -retry-join myconsul-consul-server.default.svc.cluster.local   // This command can also be used for joining to Consul Server
consul agent --config-dir=/consul/config -retry-join "provider=k8s label_selector=\"app=consul,component=server\""