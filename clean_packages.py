#!/usr/bin/python2
# -*- coding: utf-8 -*-

import os
import requests

token = os.getenv("PACKAGECLOUD_TOKEN")
repository = "jollheef/henhouse"
api_url = "https://%s:@packagecloud.io/api/v1/repos/%s/" % (token, repository)
name = 'henhouse'

def delete_package(filename):
    response = requests.delete(api_url+filename)
    print("Delete package status: %d" % response.status_code)

packages = requests.get(api_url+"/packages.json").json()

for pkg in packages[:-9]:
    if pkg['name'] == name:
        if pkg['release'] != '0':
            delete_package(pkg['distro_version']+"/"+pkg['filename'])
