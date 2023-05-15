SHELL=/bin/bash -o pipefail

all: glif

glif: main.go
	go build -o glif .

install: glif
	cp glif /usr/local/bin/glif
