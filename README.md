# HAR file management library

A Golang library for parsing and managing HAR (HTTP Archive) files. Provides easy-to-use structs and functions for loading, inspecting, and manipulating HAR files, making HTTP traffic analysis and debugging simpler.

- [X] Check for case sensitivity in HeaderOrder keys
- [ ] Add Remote Add Address (proxy used)
- [ ] Fix bodySize

## How to compile

### MacOS

```bash
cd example
GOOS=darwin GOARCH=amd64 go build example/main.go
```
