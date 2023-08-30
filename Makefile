SHELL=/bin/bash -o pipefail

all: glif

glif: main.go
	go build $(GOFLAGS) -o glif .
.PHONY: glif

clean:
	rm -f glif

config:
	mkdir -p ~/.glif
	cp config.toml ~/.glif/config.toml

calibnet-config:
	mkdir -p ~/.glif/calibnet
	cp calibnet-config.toml ~/.glif/calibnet/config.toml

calibnet: GOFLAGS+=-tags=calibnet
calibnet: glif

advanced: GOFLAGS+=-tags=advanced
advanced: glif

install:
	cp glif /usr/local/bin/glif
.PHONY: install
