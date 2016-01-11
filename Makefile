GOFLAGS = -tags netgo
GITHUB_ACCOUNT = ObjectIsAdvantag
CONFIG=config-standalone.private

default: dev

.PHONY: all
all : build build-recorder run

.PHONY: recorder
recorder: build-recorder run-recorder

.PHONY: build-recorder
build-recorder:
	rm -f recorder-server.exe recorder-server
	go build recorder-server.go

.PHONY: run-recorder
run-recorder: build-recorder
	./recorder-server.exe -port 8081 -formID filename -directory "./uploads" -upload "recordings" -download "audio" -logtostderr=true -v=5

.PHONY: run
run:
	rm -f messages.db
	(./answering-machine.exe -port 8080 -logtostderr=true -v=5  --conf=$(CONFIG) &)
	(./recorder-server.exe -port 8081 -formID filename -directory "./uploads" -upload "recordings" -download "audio" -logtostderr=true -v=5 &)
	(lt -p 8080 -s answeringmachine &)
	(lt -p 8081 -s recorder &)

.PHONY: capture
capture:
	(./answering-machine.exe -port 8080 -logtostderr=true -v=5 - --conf=$(CONFIG) &)
	(../smartproxy/smartproxy.exe -capture -port 9090 -serve 127.0.0.1:8080 &)
	(./recorder-server.exe -port 8081 -formID filename -directory "./uploads" -upload "recordings" -download "audio" -logtostderr=true -v=5 &)
	(lt -p 9090 -s answeringmachine &)
	(lt -p 8081 -s recorder &)

.PHONY: dev
dev: clean build
	./answering-machine.exe -logtostderr=true -v=5 --conf=$(CONFIG)

.PHONY: build
build: clean build-recorder
	go build $(GOFLAGS) answering-machine.go

.PHONY: debug
debug:
	godebug build $(GOFLAGS) -instrument github.com/$(GITHUB_ACCOUNT)/answering-machine/machine,github.com/$(GITHUB_ACCOUNT)/answering-machine/tropo answering-machine.go
	./answering-machine.debug -logtostderr=true -v=5

.PHONY: linux
linux:
	GOOS=linux GOARCH=amd64 go build $(GOFLAGS) answering-machine.go

.PHONY: windows
windows:
	GOOS=windows GOARCH=amd64 go build $(GOFLAGS) answering-machine.go

.PHONY: docker
docker: linux
	docker build -t $(GITHUB_ACCOUNT)/answering-machine .

.PHONY: clean
clean:
	rm -f answering-machine answering-machine.exe answering-machine.zip answering-machine.debug
	rm -f recorder-server recorder-server.exe recorder-server.zip recorder-server.debug

.PHONY: erase
erase:
	rm -f *.db
	rm -f ./log/*
	rm -f ./uploads/*

.PHONY: archive
archive:
	git archive --format=zip HEAD > answering-machine.zip
