"""
Add watches to the consul client config
"""

import sys
import json
import copy

import boto3
from pprint import pprint
from botocore.exceptions import *
from tenacity import retry, wait_fixed, stop_after_delay, retry_if_exception_type


HELP = """
    Run this script as:
        python add_watches.py <template> <KEYPREFIXES>
        template    : config template for the consul client
        KEYPREFIXES : comma seperated keyprefixes to watch in consul kv store
"""

HOTLOAD = "http://localhost:7001"

WATCH = {
    "type": "keyprefix",
    "prefix": "key",
    "handler_type": "http",
    "http_handler_config": {
        "path": HOTLOAD,
        "method": "POST",
        "tls_skip_verify": False
    }
}

def get_credentials():
    """
    Get temporary credentials for the ec2 role
    """
    session = boto3.Session()
    credentials = session.get_credentials()

    credentials = credentials.get_frozen_credentials()
    access_key = credentials.access_key
    secret_key = credentials.secret_key
    print(credentials)


def get_consul_encrypt_key(parameter_key):
    ssm = boto3.client('ssm', region_name='us-east-1')
    response = ssm.get_parameter(Name=parameter_key, WithDecryption=True)
    print(response)
    return response['Parameter']['Value']


def construct_watch(prefix):
    """ Construct a watch from WATCH and prefix """
    watch = copy.deepcopy(WATCH)
    watch['prefix'] = prefix
    return watch

def retry_on_nocreds(exc):
    return isinstance(exc, NoCredentialsError)

@retry(retry=retry_if_exception_type(NoCredentialsError),
        wait=wait_fixed(2),
        stop=stop_after_delay(10))
def add_encrypt_key(config, parameter_key):
    """
    Add consul gossip encryption key from AWS parameter store
    """
    try:
        config['encrypt'] = get_consul_encrypt_key(parameter_key=parameter_key)

    except NoCredentialsError as err:
        print(str(err))
        raise NoCredentialsError

    except ClientError as err:
        print(str(err))
        sys.exit(1)


def add_watches(config, keyprefixes):
    """
    Adds watches for keyprefixes in the config dictionary

    config: a consul client config json as a dictionary
    keyprefixes: a comma seperated string of keyprefixes
    """
    keyprefixes = keyprefixes.split(', ')
    print(keyprefixes)
    for prefix in keyprefixes:
        config['watches'].append(construct_watch(prefix))


def read_config(jsonfile):
    """ 
    Reads a json into dictionary 
    jsonfile: a file handle for a json file
    """
    return json.loads(jsonfile.read())


def generate_config(template, keyprefixes):
    """
    Adds a watch for each keyprefix in the input string <keyprefixes>
    keyprefixes: comma seperated keyprefixes to watch in consul kv store

    A 'config.json' file is created with the updated config.
    """
    with open(template) as templ:
        print('IN: ' + template)
        config = read_config(templ)

    with open('config.json', 'w') as conf:
        add_watches(config, keyprefixes)
        add_encrypt_key(config, 'consul_gossip_encrypt_key')
        
        pprint(config)

        conf.write(json.dumps(config, indent=4))
        print('OUT: ' + 'config.json')
        

if __name__ == '__main__':
    # get_credentials()
    try:
        template = sys.argv[1]
        keyprefixes = sys.argv[2]
        generate_config(template, keyprefixes)
    except IndexError:
        print(HELP)
        sys.exit(1)
    except IOError:
        print('file "' + template + '" does not exist!')
        sys.exit(1)

