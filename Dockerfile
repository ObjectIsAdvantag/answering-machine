# Replace with scratch when done
FROM gliderlabs/alpine:3.3

MAINTAINER "Stève Sfartz" <steve.sfartz@gmail.com>

COPY . /machine

EXPOSE 8080

ENTRYPOINT ["/machine/answering-machine", "--port=8080", "-logtostderr=true", "-v=5", "--messages=/machine/messages/messages-en.json", "--env=/machine/env.json"]



