![fetch-proxy](https://raw.githubusercontent.com/kcmerrill/fetch-proxy/master/assets/fetch.png "fetch-proxy")

[![Build Status](https://travis-ci.org/kcmerrill/fetch-proxy.svg?branch=master)](https://travis-ci.org/kcmerrill/fetch-proxy) [![Join the chat at https://gitter.im/fetch-proxy/Lobby](https://badges.gitter.im/fetch-proxy/Lobby.svg)](https://gitter.im/fetch-proxy/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Fetch is a simple proxy that automagically routes web traffic to running docker containers to host ports. Great for dev/ci environments. Works great in production for hosts that have web sites running on one machine.

* Healthchecks for services
* Secure connections using lets encrypt
* Default service if mapping not found
* Response timeouts
* Automagically maps new containers
* Ideal for dev environments
* Zero downtime deployments out of the box

![fetch-proxy](https://raw.githubusercontent.com/kcmerrill/fetch-proxy/master/assets/fetch-proxy.gif "fetch-proxy gif")

# Binaries || Installation

[![MacOSX](https://raw.githubusercontent.com/kcmerrill/go-dist/master/assets/apple_logo.png "Mac OSX")] (http://go-dist.kcmerrill.com/kcmerrill/fetch-proxy/mac/amd6) [![Linux](https://raw.githubusercontent.com/kcmerrill/go-dist/master/assets/linux_logo.png "Linux")] (http://go-dist.kcmerrill.com/kcmerrill/fetch-proxy/linux/amd64)

or via go:

`$ go get github.com/kcmerrill/fetch-proxy`

via docker:

`$ docker run -d -p 80:80 -p 443:443 -v /var/run/docker.sock:/var/run/docker.sock --restart=always --name fetch-proxy kcmerrill/fetch-proxy --containerized --insecure`

Use the `--insecure` flag if you are working on localhost, else enjoy `https://` urls from letsencrypt.org

## Zero downtime deployments

By deploying containers with `_` in their names, this denotes different versions to `fetch-proxy` along with their start times. If you start a container called `test_v1.0`, test.domain.tld will route traffic to that specific container. If you launch another container named `test_v1.1` test.domain.tld will now start taking in that traffic once the container becomes online(via a healthcheck).

## Custom Mappings

Lets say you are not using docker images for some ports, or perhaps you need to map multiple subdomains to one particular port. Simply pass in the `config` flag, with a location to a config file. This will be used for more later, but for now, create a yaml file with a key of `forward` and as a multi dimensional array, pass in the `subdomain name: port`. Feel free to update this whenever you'd like, give fetch-proxy a few seconds to read in the new changes and then go to the new subdomain. Here is an example of multiple subdomains going to a single port:

```yaml
forward:
    mystaticpage: 1234
    mystaticpage2: 1234
    mystaticpage3: 1234
```

## Questions/Comments/Feedback?

Would love to hear it. Email me at kcmerrill [at] gmail [dot] com
