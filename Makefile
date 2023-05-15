SHELL=/bin/bash -o pipefail

all: glif

glif: main.go
	go build -o glif .

clean:
	rm -f glif

config:
	mkdir -p ~/.glif
	cp config.toml ~/.glif/config.toml

calibnet-config:
	mkdir -p ~/.glif
	cp calibnet-config.toml ~/.glif/config.toml

install: glif
	cp glif /usr/local/bin/glif
