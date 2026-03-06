.PHONY: demo clean

# Builds the hotreload engine and runs the demo server
demo:
	@echo "Starting the hot reload demonstration..."
	go run . --root ./testserver --build "go build -o ./bin/server ./testserver" --exec "./bin/server"

# Cleans up the compiled binaries
clean:
	@echo "Cleaning up binaries..."
	rm -rf ./bin