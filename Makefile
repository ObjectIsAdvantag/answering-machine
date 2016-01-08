GOFLAGS = -tags netgo
USERNAME = ObjectIsAdvantag

default: dev

.PHONY: all
all : build build-recorder run

.PHONY: build-recorder
build-recorder:
	rm -f recorder-server.exe recorder-server
	go build recorder-server.go

.PHONY: recorder
recorder: build-recorder
	./recorder-server.exe -logtostderr=true -v=5

.PHONY: run
run:
	(./answering-machine.exe -port 8080 -logtostderr=true -v=5 &)
	(./recorder-server.exe -port 8081 -logtostderr=true -v=5 &)
	(lt -p 8080 -s answeringmachine &)
	(lt -p 8081 -s recorder &)

.PHONY: capture
capture:
	(./answering-machine.exe -port 8080 -logtostderr=true -v=5 &)
	(../smartproxy/smartproxy.exe -capture -port 9090 -serve 127.0.0.1:8080 &)
	(./recorder-server.exe -port 8081 -logtostderr=true -v=5 &)
	(lt -p 9090 -s answeringmachine &)
	(lt -p 8081 -s recorder &)

.PHONY: dev
dev: clean build
	./answering-machine.exe -logtostderr=true -v=5

.PHONY: build
build: clean
	go build $(GOFLAGS) answering-machine.go

.PHONY: debug
debug:
	godebug build $(GOFLAGS) -instrument github.com/$(USERNAME)/answering-machine/machine,github.com/$(USERNAME)/answering-machine/tropo answering-machine.go
	./answering-machine.debug -logtostderr=true -v=5

.PHONY: linux
linux:
	GOOS=linux GOARCH=amd64 go build $(GOFLAGS) answering-machine.go

.PHONY: windows
windows:
	GOOS=windows GOARCH=amd64 go build $(GOFLAGS) answering-machine.go

.PHONY: docker
docker: linux
	docker build -t $(USERNAME)/answering-machine .

.PHONY: clean
clean:
	rm -f answering-machine answering-machine.exe answering-machine.zip answering-machine.debug
	rm -f ./log/*

.PHONY: erase
erase:
	rm -f data.db

.PHONY: archive
archive:
	git archive --format=zip HEAD > answering-machine.zip
