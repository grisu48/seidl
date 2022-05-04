default: all

all: seidl

seidl: cmd/seidl/seidl.go
	go build -o $@ $^

seidl-static: cmd/seidl/seidl.go
	CGO_ENABLED=0 go build -ldflags="-w -s" -o seidl $^

install: seidl
	install seidl ~/bin/
