# TL;DR
Ever had a bunch of docker containers running on random ports on your machine? Ever needed a reverse proxy to figure out all the containers, and route traffic to them automagically? Enter automagicproxy.

Basically, any docker container with a private port of 80, that has a public port of X and a specific container name will get routed to it.

A quick example:

docker run -d --name mycoolapp -P kcmerrill/base

mycoolapp.kcmerrill.com -> automagicproxy:80 -> mycoolapp:46743

# Usage
Automagicproxy can listen to whatever port you'd like, but by default, lets let it listen on port 80.

`docker run -ti -p 80:80 -v /var/run/docker.sock:/var/run/docker.sock --name automagicproxy kcmerrill/automagicproxy`

Now, start up your docker containers as you normally would. Any docker containers that use a private port of 80, will be automagically mapped based on their container names.

# Features
- Health checks
- Dynamically load container private ips and route accordingly based on container names

# Todo
- Round robin
- Stickied requests
- Debug override console
