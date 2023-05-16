# Do a test build and run of the project
all::
	$(MAKE) build
	$(MAKE) start

# Builds the application for the current system
.PHONY: build
build:: src/*.go
	@go build \
		-o build/xec \
		./src/

# Starts the application from a local build
.PHONY: start
start:: build
	@./build/xec

# Starts the application from a local build with the show parameter
.PHONY: show
show:: build
	@./build/xec --show
