
# Leave it as is if you just clone the original repo
# Replace this with your github account if you've forked the original repo,
GITHUB_ACCOUNT = ObjectIsAdvantag

# Leave as if if you 're just using, customizing or extending the machine
# But, change this to your docker account if you plan to package and release your own answering machine
DOCKER_ACCOUNT = objectisadvantag

# Default configuration is english language for Text to Speech, and recordings are backed by Tropo Cloud
CONFIG=--env=env-tropofs.json --messages=messages-en.json
#CONFIG=--env=env-recorder.json --messages=messages-fr.json

# Command line to start the machine for Mac, Linux and Windows
STARTUP=./answering-machine -port 8080 -logtostderr=true -v=5 $(CONFIG)
GOFLAGS = -tags netgo

default: all

.PHONY: all
all : build build-recorder run

.PHONY: goget
goget: 
	go get github.com/golang/glog
	go get github.com/boltdb/bolt
	go get github.com/paked/configure

.PHONY: recorder
recorder: build-recorder run-recorder

.PHONY: build-recorder
build-recorder:
	rm -f recorder-server recorder-server.exe recorder-server.zip recorder-server.debug
	go build recorder-server.go

.PHONY: run-recorder
run-recorder: build-recorder
	./recorder-server -port 8081 -formID filename -directory "./uploads" -upload "recordings" -download "audio" -logtostderr=true -v=5

.PHONY: build
build: goget build-recorder build-machine

.PHONY: build-machine
build-machine:
	rm -f answering-machine answering-machine.exe answering-machine.zip answering-machine.debug
	go build $(GOFLAGS) answering-machine.go

.PHONY: debug
debug:
	godebug build $(GOFLAGS) -instrument github.com/$(GITHUB_ACCOUNT)/answering-machine/machine answering-machine.go
	./answering-machine.debug -logtostderr=true -v=5

.PHONY: run
run:
	rm -f messages.db
	($(STARTUP) &)
	# (./answering-machine -port 8080 -logtostderr=true -v=5  $(CONFIG) &)
	(./recorder-server -port 8081 -formID filename -directory "./uploads" -upload "recordings" -download "audio" -logtostderr=true -v=5 &)
	(lt -p 8080 -s answeringmachine &)
	(lt -p 8081 -s recorder &)

.PHONY: capture
capture:
	rm -f messages.db
	(./answering-machine -port 8080 -logtostderr=true -v=5 $(CONFIG) &)
	(../smartproxy/smartproxy -capture -port 9090 -serve 127.0.0.1:8080 &)
	(./recorder-server -port 8081 -formID filename -directory "./uploads" -upload "recordings" -download "audio" -logtostderr=true -v=5 &)
	(lt -p 9090 -s answeringmachine &)
	(lt -p 8081 -s recorder &)

.PHONY: dev
dev:
	rm -f messages.db
	(./recorder-server -port 8081 -formID filename -directory "./uploads" -upload "recordings" -download "audio" -logtostderr=true -v=5 &)
	./answering-machine -port 8080 -logtostderr=true -v=5  $(CONFIG)

.PHONY: clean
clean:
	rm -f answering-machine answering-machine.exe answering-machine.zip answering-machine.debug
	rm -f recorder-server recorder-server.exe recorder-server.zip recorder-server.debug

.PHONY: erase
erase:
	rm -f *.db
	rm -f ./log/*
	rm -f ./uploads/*
	rm -f ./dist

.PHONY: linux
linux: goget
	GOOS=linux GOARCH=amd64 go build $(GOFLAGS) answering-machine.go
	GOOS=linux GOARCH=amd64 go build $(GOFLAGS) recorder-server.go

.PHONY: mac
mac: goget
	GOOS=darwin GOARCH=amd64 go build $(GOFLAGS) answering-machine.go
	GOOS=darwin GOARCH=amd64 go build $(GOFLAGS) recorder-server.go

.PHONY: windows
windows: goget
	GOOS=windows GOARCH=amd64 go build $(GOFLAGS) answering-machine.go
	GOOS=windows GOARCH=amd64 go build $(GOFLAGS) recorder-server.go

.PHONY: dist
dist: linux
	rm -rf dist
	mkdir dist
	cp answering-machine dist/
	cp recorder-server dist/
	mkdir dist/logs
	mkdir dist/uploads
	mkdir dist/messages
	cp messages-en.json dist/messages
	cp messages-fr.json dist/messages
	cp env-tropofs.json dist/env.json
	cp Dockerfile dist/

.PHONY: docker
docker: dist
	cd dist; docker build -t $(DOCKER_ACCOUNT)/answeringmachine .

.PHONY: archive
archive:
	git archive --format=zip HEAD > answering-machine.zip

.PHONY: pkg
pkg: pkg-windows pkg-linux

.PHONY: pkg-windows
pkg-windows: windows
	rm -rf pkg/windows
	mkdir pkg/windows
	cp answering-machine.exe pkg/windows
	cp recorder-server.exe pkg/windows
	mkdir pkg/windows/logs
	mkdir pkg/windows/uploads
	cp messages-en.json pkg/windows
	cp messages-fr.json pkg/windows
	cp env-tropofs.json pkg/windows/env.json

.PHONY: pkg-linux
pkg-linux: linux
	rm -rf pkg/linux
	mkdir pkg/linu
	cp answering-machine pkg/linux
	cp recorder-server pkg/linux
	mkdir pkg/linux/logs
	mkdir pkg/linux/uploads
	cp messages-en.json pkg/linux
	cp messages-fr.json pkg/linux
	cp env-tropofs.json pkg/linux/env.json
	tar -zcvf linux.tgz pkg/linux

.PHONY: graph
graph:
	godepgraph github.com/ObjectIsAdvantag/answering-machine | dot -Tpng -o answeringmachine.png
