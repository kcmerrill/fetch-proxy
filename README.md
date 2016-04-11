# Automagicproxy
A proxy written in GO that will allow you rather easily forward HTTP requests along with automagic docker container port forwarding. 

# Usage
docker run -ti -p 80:80 -v /var/run/docker.sock:/var/run/docker.sock --name automagicproxy kcmerrill/automagicproxy

# Features
- Health checks
- Dynamically load container private ips and route accordingly based on container names

# Todo
- Round robin
- Stickied requests
- Debug override console
