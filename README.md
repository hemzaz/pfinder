# pfinder: Cross-Platform Process Finder

`pfinder` is a cross-platform tool written in Go that helps you find processes based on various criteria, such as file usage, PID, command strings, regex patterns, and network ports.  
Designed to be efficient and robust,  
`pfinder` works seamlessly on both macOS and Linux without relying on external system commands.

## Features

- **Cross-Platform Compatibility**: Supports macOS and Linux with platform-specific optimizations.
- **Find Processes by File Usage**: Identify which process is using a specific file or directory.
- **Find Processes by PID**: Retrieve detailed information about a process by its PID.
- **Find Processes by Command String or Regex**: Filter running processes based on command strings or regex patterns.
- **Find Processes by Port**: Identify processes listening on a specific network port.
- **Efficient and Concurrent**: Utilizes Go's concurrency features to handle large numbers of processes efficiently.
- **Pure Go Implementation**: No external dependencies or system commands required.

## Installation

To build `pfinder` from source, follow these steps:

### Prerequisites

1. **Install Go**: Ensure you have Go installed. You can download and install Go from [golang.org](https://golang.org/dl/).

2. **Set Up Go Environment**: Configure your Go environment variables, especially `GOPATH`, which is the directory where your Go workspace is located.

### Steps to Build

1. **Clone the Repository**: Clone the `pfinder` repository from GitHub to your local machine.

    ```bash
    git clone https://github.com/hemzaz/pfinder.git
    cd pfinder
    ```

2. **Install Dependencies**: Install the required Go packages.

    ```bash
    go get github.com/mitchellh/go-ps
    go get github.com/shirou/gopsutil/net
    ```

3. **Build the Executable**: Use the `go build` command to compile the source code into an executable.

    ```bash
    go build -o pfinder
    ```

4. **Verify the Executable**: Ensure the `pfinder` executable is created in the current directory.

    ```bash
    ls -l pfinder
    ```

## Usage

`pfinder` can be used with multiple arguments to find processes based on various criteria. The results will be aggregated and displayed.

```bash
./pfinder <arguments>...
```

### Arguments

- **Path**: If the argument is a path to an existing file or folder, `pfinder` reports the PID locking the resource (if any).
- **PID**: If the argument is a PID, `pfinder` provides detailed information about the process.
- **String**: If the argument is a string, `pfinder` filters running processes case-insensitively based on the string in their command line.
- **Regex**: If the argument is a regex, `pfinder` filters running processes using the regex pattern.
- **Port**: If the argument is a port (prefixed with ':'), `pfinder` reports the processes listening on that port.

### Examples

1. **Find the process using a specific file**:

    ```bash
    ./pfinder /path/to/file
    ```

2. **Get details of a process by PID**:

    ```bash
    ./pfinder 1234
    ```

3. **Find processes by a command string**:

    ```bash
    ./pfinder "go run"
    ```

4. **Find processes by a regex pattern**:

    ```bash
    ./pfinder ".*go.*"
    ```

5. **Find processes by a port**:

    ```bash
    ./pfinder :8080
    ```

6. **Aggregate multiple arguments**:

    ```bash
    ./pfinder "go run" :8080 /path/to/file
    ```

## How It Works

### Path Handling

- **macOS**: Uses native Go libraries to check file usage by iterating over processes and their file descriptors.
- **Linux**: Reads the `/proc` filesystem to check which process is using the specified file or directory.

### PID Handling

- Retrieves detailed information about the process including PID, PPID, user, and command.

### String and Regex Handling

- Filters running processes based on command strings or regex patterns in a case-insensitive manner.

### Port Handling

- Uses the `gopsutil` library to identify processes listening on specified network ports.

## Contributing

We welcome contributions! Here are a few ways you can help:

1. **Report Bugs**: Use the issue tracker to report bugs.
2. **Fix Bugs**: If you find a bug and know how to fix it, please do so! Pull requests are welcome.
3. **Request Features**: If you have ideas for new features, we'd love to hear them.

### How to Contribute

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/fooBar`)
3. Commit your changes (`git commit -am 'Add some fooBar'`)
4. Push to the branch (`git push origin feature/fooBar`)
5. Create a new Pull Request

## Acknowledgements

- [mitchellh/go-ps](https://github.com/mitchellh/go-ps): A library for listing processes on a system.
- [shirou/gopsutil](https://github.com/shirou/gopsutil): A library for retrieving system and process information.

---
## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
---
Authored by: **hemzaz the frogodile** 🐸🐊