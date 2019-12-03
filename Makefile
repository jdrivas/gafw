package := github.com/jdrivas/gafw
now := $(shell date +%s)
timeflag := -X ${package}/version.unixtime=$(now)
hash := $(shell git rev-parse --short HEAD)
hashflag := -X $(package)/version.gitHash=$(hash)
tag := $(shell git tag -l --points-at HEAD)
tagflag := -X $(package)/version.gitTag=$(tag)
ld_args :=  $(timeflag) $(hashflag) $(tagflag)

build:
	go build "-ldflags=$(ld_args)"

install: 
	go install "-ldflags=$(ld_args)"
