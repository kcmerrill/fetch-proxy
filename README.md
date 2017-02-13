![fetch-proxy](https://raw.githubusercontent.com/kcmerrill/fetch-proxy/master/assets/fetch.png "fetch-proxy")

[![Build Status](https://travis-ci.org/kcmerrill/fetch-proxy.svg?branch=master)](https://travis-ci.org/kcmerrill/fetch-proxy) [![Join the chat at https://gitter.im/kcmerrill/fetch-proxy](https://badges.gitter.im/kcmerrill/fetch-proxy.svg)](https://gitter.im/kcmerrill/fetch-proxy?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

## What is it
A simple proxy that automagically routes web traffic to running docker containers to host ports. Great for dev/ci environments. Works great in production for hosts that have web sites running on one machine.

## Getting Started
` $ go get github.com/kcmerrill/fetch-proxy`

or via docker:

` $ docker run -d -p 80:80 --restart=always --name=fetch-proxy kcmerrill/fetch-proxy`

## How it works
Normally if you have multiple running docker containers on a host, you'd need a proxy to configure the ports and map between the host machine and the containers. With `fetch-proxy` you simply start coding. 
