**Quick Summary**

If you're working on multiple projects within docker, chances are likely you are either going to an ip:port instead of a domain name. Or you've configured some sort of reverse proxy(which is what this is, but it's automagic with no configuration needed).

You might even have a DNS entry that is *.domain.com on your local network that points to your docker host(if not, it's easy to create). This is where this container fits in.

By simply running this container, it will find all web containers with a private port of 80, and using haproxy will forward all traffic based on the containers name to the appropriate port.

Example:

    CONTAINER ID        IMAGE                             COMMAND                CREATED             STATUS              PORTS                         NAMES
    cf01ff8ceb11        kcmerrill/automagicproxy:latest   "\"/bin/sh -c 'hapro   2 seconds ago       Up 2 seconds        443/tcp, 0.0.0.0:80->80/tcp   haproxy
    047b37031f66        sequenceiq/socat:latest           "/start"               44 minutes ago      Up 44 minutes       0.0.0.0:2375->2375/tcp        docker-http
    1171fdf7d5b0        kcmerrill/base:latest             "/bin/sh -c '/usr/sb   56 minutes ago      Up 56 minutes       0.0.0.0:1201->80/tcp          hello
    6373866cead4        kcmerrill/base:latest             "/bin/sh -c '/usr/sb   56 minutes ago      Up 56 minutes       0.0.0.0:1200->80/tcp          base
    6cab74ba4220        kcmerrill/vitaminc:latest         "/bin/sh -c 'cd /opt   58 minutes ago      Up 58 minutes       0.0.0.0:9999->9999/tcp        vitaminc


As you can see, I have multiple websites running. I have containers named "hello" and "base" running. They are simple websites that are configured to run. If you were to go to my host machine http://192.168.59.103:1200/ with the appropriate port you'd see the contents of the pages.

If you have a wildcard DNS entry that points to 192.168.59.103, lets call it dev.com, this container will automatically  forward traffic to the appropriate container.

So ... http://base.dev.com would route traffic to the base container, http://hello.dev.com will automatically forward traffic to the hello container.

To run, you need to ensure that socat is running. socat allows the remote docker api to be accessed via regular http. This will sadly make this project only useful in dev, however I'll be working on getting the certs working, or coming up with an alternate solution in the meantime.

**Running Instructions:**

Step 1:

    $(docker run sequenceiq/socat)

Step 2:

    docker run -d -p 80:80 --name automagicproxy kcmerrill/automagicproxy

Step 3:

Setup a wildcard dns entry ... or edit your host file. While editing the host file somewhat defeats the purpose of this, it's still much quicker as you don't have to reconfigure the proxy.

    192.168.59.103  vitaminc.dev.com
    192.168.59.103  base.dev.com
    192.168.59.103  hello.dev.com

Step 4:

Anytime you add new containers, simply restart the automagic proxy to discover the new containers. If you call it via docker exec, you should have 0 downtime!

    docker exec -t -i automagicproxy /automagicproxy/bin/reload

**Additional Goodies**

If you set an environment variable before running the container, or before running the reload command called "AUTOMAGIC_PORTS" with a comma separated list of ports. It will use those ports along with port 80, in haproxy's configuration. 

You can see from the container list above, that there is another container called "vitaminc" running on port 9999. It's done this way for other various reasons. However, when we first start the automagicproxy container, it will not configure this particular container, because it's private port is not port 80. By setting the environment variable "AUTOMAGIC_PORTS=9999,1234,<other_ports_here>" will allow vitaminc.dev.com to be configured along with those on private port 80. 

An example setup to include port 9999 in the configuration would look like this:

    docker run -d -p 80:80 -e AUTOMAGIC_PORTS=9999 --name automagicproxy kcmerrill/automagicproxy
