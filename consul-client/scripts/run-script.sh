#!/bin/sh

# Generated config by the initcontainer 
echo "/consul/config/*"
ls /consul/config
config_path="/consul/config/config.json"
# consul agent --config-dir=$config_path -retry-join myconsul-consul-server.default.svc.cluster.local

# consul agent \
#     --config-dir=$config_path \
#     -retry-join 'provider=k8s label_selector="app=consul,component=server"'

datacenter='kvdc1-consul-qa-travel-nextgen-nv-aws'
consul agent \
    --config-dir=$config_path \
    -retry-join "provider=aws tag_key=DataCenterName tag_value=$datacenter region=us-east-1 addr_type=private_v4"
