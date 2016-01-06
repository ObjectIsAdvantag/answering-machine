
GOFLAGS = -tags netgo
USERNAME = ObjectIsAdvantag

default: all

.PHONY: all
all : clean build dev

.PHONY: prod
prod:
	./answering-machine.exe -stderrthreshold=FATAL -log_dir=./log -v=0

.PHONY: run
run:
	./answering-machine.exe -stderrthreshold=FATAL -log_dir=./log -v=2

.PHONY: dev
dev:
	./answering-machine.exe -logtostderr=true -v=5

.PHONY: build
build:
	go build $(GOFLAGS)

.PHONY: debug
debug:
	godebug build $(GOFLAGS) -instrument github.com/$(USERNAME)/answering-machine/machine,github.com/$(USERNAME)/answering-machine/tropo
	./answering-machine.debug -logtostderr=true -v=5

.PHONY: linux
linux:
	GOOS=linux GOARCH=amd64 go build $(GOFLAGS)

.PHONY: windows
windows:
	GOOS=windows GOARCH=amd64 go build $(GOFLAGS)

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
