#!/bin/sh

# Generates config for consul client 
# Mounted as k8s volume  to the sidecar consul client

config_path="/config/config.json"
cd /scripts/add_watches

# Sleep for kiam to acknowledge the annotation
# sleep 5s

python3 add_watches.py \
    config.template.json \
    'config/app/, default/global'

mv config.json $config_path
ls $config_path
cat $config_path
