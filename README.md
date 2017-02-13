![fetch-proxy](https://raw.githubusercontent.com/kcmerrill/fetch-proxy/master/assets/fetch.png "fetch-proxy")

[![Build Status](https://travis-ci.org/kcmerrill/fetch-proxy.svg?branch=master)](https://travis-ci.org/kcmerrill/fetch-proxy) [![Join the chat at https://gitter.im/fetch-proxy/Lobby](https://badges.gitter.im/fetch-proxy/Lobby.svg)](https://gitter.im/fetch-proxy/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Fetch is a simple proxy that automagically routes web traffic to running docker containers to host ports. Great for dev/ci environments. Works great in production for hosts that have web sites running on one machine.

 * Healthchecks for services
 * Secure connections using lets encrypt
 * Default service if mapping not found
 * Response timeouts
 * Automagically maps new containers
 * Works using localhost for dev environments

![fetch-proxy](https://raw.githubusercontent.com/kcmerrill/fetch-proxy/master/assets/fetch-proxy.gif "fetch-proxy gif")

## Getting Started
` $ go get github.com/kcmerrill/fetch-proxy`

Use the `--insecure` flag if you are working on localhost, else enjoy `https://` urls from letsencrypt.org

or try it out in docker:

` $ docker run -d -p 80:80 -p 443:443 -v /var/run/docker.sock:/var/run/docker.sock --restart=always --name fetch-proxy kcmerrill/fetch-proxy --containerized --insecure`

Remove the `--containerized --insecure` if you want to use `https://` urls from letsencrypt.org

#### Questions/Comments/Feedback?
Would love to hear it. Email me at kcmerrill [at] gmail [dot] com
