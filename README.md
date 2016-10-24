# TL;DR
Ever had a bunch of websites and/or docker web containers running on random ports on your machine? Ever needed a reverse proxy to figure out all the containers, and route traffic to them automagically? Enter automagicproxy.

Basically, any docker container with a private port of 80, that has a public port and a specific container name will get routed to it(assuming you have a wildcard DNS entry setup).

A quick example:

`docker run -d --name mycoolapp -P kcmerrill/base`

mycoolapp.kcmerrill.com -> automagicproxy:80 -> mycoolapp:46743

## Binaries
![Mac OSX](https://raw.githubusercontent.com/kcmerrill/go-dist/master/assets/apple_logo.png "Mac OSX") [386](http://go-dist.kcmerrill.com/kcmerrill/automagicproxy/mac/386) | [amd64](http://go-dist.kcmerrill.com/kcmerrill/automagicproxy/mac/amd64)

![Windows](https://raw.githubusercontent.com/kcmerrill/go-dist/master/assets/windows_logo.png "Windows") [386](http://go-dist.kcmerrill.com/kcmerrill/automagicproxy/windows/386) | [amd64](http://go-dist.kcmerrill.com/kcmerrill/automagicproxy/windows/amd64)

![Linux](https://raw.githubusercontent.com/kcmerrill/go-dist/master/assets/linux_logo.png "Linux") [386](http://go-dist.kcmerrill.com/kcmerrill/automagicproxy/linux/386) | [amd64](http://go-dist.kcmerrill.com/kcmerrill/automagicproxy/linux/amd64)


## Why
At work I was previously on the web platform team. Our job is many, but we are responsible for making sure all devs have access/ability to work with all of our products including but not limited to building and maintaining web dev boxes, transitioning our old infrastructure over to docker etc etc. This is where yoda comes in at. I got tired of explaining to everybody the thousand steps it required to get all of our projects up and running. Furthermore, configuring haproxy was a bit of a tall order, as everybody was using different commands, etc to manage their dev boxes. On top of that, if anything went wrong, they came to me, so yoda and automagicproxy is a way to standardize a few things across the company.

## Usage
`$ go get github.com/kcmerrill/automagicproxy`

`$ automagicproxy --dockerized`

OR ... if you want to use it in a docker container

`$ docker run -d -p 80:80 -p 443:443 -v /var/run/docker.sock:/var/run/docker.sock --restart=always --name automagicproxy kcmerrill/automagicproxy`


Now, start up your docker containers as you normally would. Any docker containers that use a private port of 80, will be automagically mapped based on their container names.


## Features
- Health checks
- Dynamically load container private ips and route accordingly based on container names
- Automagic TLS certs thanks to https://github.com/rsc/letsencrypt

## Todo
- Round robin
- Stickied requests
- Debug override console
