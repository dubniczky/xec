# Builds the application for the current system
build::
	@go build \
		-o build/xec \
		./src/

# Starts the application from a local build
start:: build
	@./build/xec
