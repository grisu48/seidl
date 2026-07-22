default: all

all: seidl

seidl: cmd/seidl/seidl.go
	go build -o $@ $^

release:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o seidl_x86_64 ./...
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" -o seidl_arm64 ./...
	CGO_ENABLED=0 GOOS=linux GOARCH=riscv64 go build -ldflags="-w -s" -o seidl_riscv64 ./...

install: seidl
	install seidl ~/bin/
