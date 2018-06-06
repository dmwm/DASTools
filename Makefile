GOPATH:=$(PWD):${GOPATH}
export GOPATH
flags=-ldflags="-s -w"
# flags=-ldflags="-s -w -extldflags -static"
TAG := $(shell git tag | sort -r | head -n 1)

all: build

build:
	sed -i -e "s,{{VERSION}},$(TAG),g" das_cleanup.go
	sed -i -e "s,{{VERSION}},$(TAG),g" dasmaps_parser.go
	sed -i -e "s,{{VERSION}},$(TAG),g" dasmaps_validator.go
	sed -i -e "s,{{VERSION}},$(TAG),g" mongostatus.go
	go clean; rm -rf pkg;
	go build ${flags} das_cleanup.go
	go build ${flags} dasmaps_parser.go
	go build ${flags} dasmaps_validator.go
	go build ${flags} mongostatus.go
	sed -i -e "s,$(TAG),{{VERSION}},g" das_cleanup.go
	sed -i -e "s,$(TAG),{{VERSION}},g" dasmaps_parser.go
	sed -i -e "s,$(TAG),{{VERSION}},g" dasmaps_validator.go
	sed -i -e "s,$(TAG),{{VERSION}},g" mongostatus.go
	sed -i -e "s,$(TAG),{{VERSION}},g" mongoimport.go
	mv das_cleanup dasmaps_parser dasmaps_validator mongostatus mongoimport bin

clean:
	go clean; rm -rf pkg;
	rm bin/{das_cleanup,dasmaps_parser,dasmaps_validator,mongostatus}
