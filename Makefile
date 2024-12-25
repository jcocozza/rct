# Name of the binary
BINARY_NAME=rct

# Output directory for binaries
OUTPUT_DIR=bin

# Compiler and linker flags (optional)
LDFLAGS=
GOFLAGS=

# Targets for different OS/Architecture combinations
build: clean build-linux build-windows build-macos

build-linux:
	GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o $(OUTPUT_DIR)/$(BINARY_NAME)-linux-amd64 cmd/main.go

build-windows:
	GOOS=windows GOARCH=amd64 go build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o $(OUTPUT_DIR)/$(BINARY_NAME)-windows-amd64.exe cmd/main.go

build-macos:
	GOOS=darwin GOARCH=amd64 go build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-amd64 cmd/main.go
	codesign -s - $(OUTPUT_DIR)/$(BINARY_NAME)

# Clean up built binaries
clean:
	rm -rf $(OUTPUT_DIR)

# Create the output directory if it doesn't exist
$(OUTPUT_DIR):
	mkdir -p $(OUTPUT_DIR)

