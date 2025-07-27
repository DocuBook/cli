# DocuBook CLI

DocuBook CLI written on GO! Initialize, Update, Push and Deploy your Docs direct into Terminal.

## Installation

### Option 1: Using Go Install (Recommended)

1. Ensure Go (version 1.24 or newer) is installed on your system
2. Run the following command to install globally:
   ```bash
   go install github.com/DocuBook/cli@latest
   ```
3. Make sure `$GOPATH/bin` is in your system's PATH:
   ```bash
   # Check if $GOPATH/bin is in PATH
   echo $PATH | grep -q $(go env GOPATH)/bin && echo "✓ GOPATH/bin is in PATH" || echo "GOPATH/bin is NOT in PATH"

   # If not in PATH, add it (for current session only)
   export PATH=$PATH:$(go env GOPATH)/bin
   ```

4. Verify the installation was successful:
   ```bash
   docubook --version
   ```

### Option 2: Build from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/DocuBook/cli.git
   cd cli
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Build the binary:
   ```bash
   go build -o docubook
   ```
4. (Optional) Move the binary to your system's PATH

## Usage

1. Open a terminal and navigate to your project directory
2. Run the command:
   ```bash
   docubook cli
   ```
3. Follow the on-screen instructions to create a new project

## System Requirements

- Go 1.24 or newer
- Git (for version control)
- Internet access (for downloading dependencies)

## Optional Configuration

### Adding $GOPATH/bin to PATH

To make `docubook` command available globally, add the following line to your shell configuration file:

For Zsh (recommended for macOS):
```bash
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc
```

For Bash:
```bash
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bash_profile
source ~/.bash_profile
```

## Troubleshooting

### Command not found: docubook
If you get this error after installation:

1. Verify Go is installed and in your PATH:
   ```bash
   go version
   ```

2. Check if the binary was installed:
   ```bash
   ls -l $(go env GOPATH)/bin/docubook
   ```

3. If the file exists but you still get the error, add the Go bin directory to your PATH (see Optional Configuration above).

### Permission Denied
If you encounter permission issues:

```bash
chmod +x $(go env GOPATH)/bin/docubook
```

### Outdated Version
To update to the latest version:
```bash
go install github.com/DocuBook/cli@latest
```

### Still having issues?
- Make sure your Go environment is set up correctly:
  ```bash
  go env
  ```
- Check for error messages during installation
- Ensure you have write permissions to the Go bin directory
