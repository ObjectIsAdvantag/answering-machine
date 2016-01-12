# Replace with scratch when done
FROM gliderlabs/alpine:3.3

MAINTAINER "Stève Sfartz" <steve.sfartz@gmail.com>

COPY . /machine

EXPOSE 8080 8081

ENTRYPOINT ["/machine/answering-machine", "-port", "8080", "--env=/machine/env.json", "--messages=/machine/conf/messages-fr.json", "-logtostderr=true", "-v=5"]



