# go-genetic-ml makefile
# Copyright (C) 2018 Daniel Wilson
# GNU GPL v3.0 License - See LICENSE
# https://github.com/Danw33/go-genetic-ml

# Default: Make the "build" version without debug symbols
all: build

# Build without debug symbols and pack using UPX (smallest possible output)
packed: build pack

# Build without debug symbols (Smaller output executable) for the current OS and Arch
build:
	go build -ldflags "-s -w" ./src/go-genetic-ml.go
	if [ -a ./go-genetic-ml ]; then chmod +X ./go-genetic-ml; fi;

# Debug build with debug symbols (Larger output executable) for the current OS and Arch
debug:
	go build ./src/go-genetic-ml.go
	if [ -a ./go-genetic-ml ]; then chmod +X ./go-genetic-ml; fi;

# Pack the compiled file using UPX
pack:
	if [ -a ./go-genetic-ml ]; then upx -9 -v ./go-genetic-ml; fi;

# Install the build (with systemd service if the host OS uses systemd)
install:
	cp ./go-genetic-ml /usr/local/bin/go-genetic-ml

# Clean the working directory of any existing builds
clean:
	if [ -a ./go-genetic-ml ]; then rm ./go-genetic-ml; fi;
