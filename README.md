# TL;DR
Ever had a bunch of websites and/or docker web containers running on random ports on your machine? Ever needed a reverse proxy to figure out all the containers, and route traffic to them automagically? Enter automagicproxy.

Basically, any docker container with a private port of 80, that has a public port of X and a specific container name will get routed to it(assuming you have a wildcard DNS entry setup).

A quick example:

`docker run -d --name mycoolapp -P kcmerrill/base`

mycoolapp.kcmerrill.com -> automagicproxy:80 -> mycoolapp:46743

# Why
At work I was previously on the web platform team. Our job is many, but we are responsible for making sure all devs have access/ability to work with all of our products including but not limited to building and maintaining web dev boxes, transitioning our old infastructure over to docker etc etc. This is where yoda comes in at. I got tired of explaining to everybody the thousand steps it required to get all of our projects up and running. Furthermore, configuring haproxy was a bit of a tall order, as everybody was using different commands, etc to manage their dev boxes. On top of that, if anything went wrong, they came to me, so yoda and automagicproxy is a way to standardize a few things across the company.

# Usage
`$ go get github.com/kcmerrill/automagicproxy`

`$ automagicproxy --dockerized`

Now, start up your docker containers as you normally would. Any docker containers that use a private port of 80, will be automagically mapped based on their container names.


# Features
- Health checks
- Dynamically load container private ips and route accordingly based on container names
- Automagic TSL certs thanks to https://github.com/rsc/letsencrypt

# Todo
- Round robin
- Stickied requests
- Debug override console
