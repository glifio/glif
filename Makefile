SHELL=/bin/bash -o pipefail

all: glif

glif: main.go
	go build -o glif .

clean:
	rm -f glif

config:
	mkdir -p ~/.config/glif
	cp config.toml ~/.config/glif/config.toml

calibnet-config:
	mkdir -p ~/.config/glif
	cp calibnet-config.toml ~/.config/glif/config.toml

install: glif
	cp glif /usr/local/bin/glif
