![fetch-proxy](https://raw.githubusercontent.com/kcmerrill/fetch-proxy/master/assets/fetch.png "fetch-proxy")

[![Build Status](https://travis-ci.org/kcmerrill/fetch-proxy.svg?branch=master)](https://travis-ci.org/kcmerrill/fetch-proxy) [![Join the chat at https://gitter.im/fetch-proxy/Lobby](https://badges.gitter.im/fetch-proxy/Lobby.svg)](https://gitter.im/fetch-proxy/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Fetch is a simple proxy that automagically routes web traffic to running docker containers to host ports. Great for dev/ci environments. Works great in production for hosts that have web sites running on one machine.

 * Healthchecks for services
 * Secure connections using lets encrypt(on by default, disable by using `--insecure` flag)
 * Default service if mapping not found
 * Response timeouts
 * Automagically maps new containers
 * Works using localhost for dev environments

## Getting Started
` $ go get github.com/kcmerrill/fetch-proxy`

or via docker:

` $ docker run -d -p 80:80 -p 443:443 --restart=always --name=fetch-proxy kcmerrill/fetch-proxy`

or via docker on localhost:

` $ docker run -d -p 80:80 -p 443:443 --restart=always --name=fetch-proxy kcmerrill/fetch-proxy --containerized --insecure`


## How it works
Normally if you have multiple running docker containers on a host, you'd need a proxy to configure the ports and map between the host machine and the containers. With `fetch-proxy` you simply start coding. 

## Demo
