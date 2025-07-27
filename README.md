# 🚀 DocuBook CLI

DocuBook CLI is a Go-based tool that helps you initialize, update, and deploy documentation directly from your terminal.

<div align="center">
  <img src="https://media.giphy.com/media/JtjAa4u5GTHd1dph0G/giphy.gif" alt="Coding Panda" width="600" />
</div>


## 📋 Table of Contents
- [System Requirements](#-system-requirements)
- [Installation](#-installation)
  - [Method 1: Using Go Install (Recommended)](#method-1-using-go-install-recommended)
  - [Method 2: Build from Source](#method-2-build-from-source)
- [Configuration](#️-configuration)
  - [Adding to PATH](#adding-to-path)
  - [Verifying Installation](#verifying-installation)
- [Usage](#-usage)
- [Troubleshooting](#-troubleshooting)
  - [Command Not Found](#command-not-found)
  - [Permission Denied](#permission-denied)
  - [Update to Latest Version](#update-to-latest-version)
  - [Uninstalling DocuBook CLI](#uninstalling-docubook-cli)
  - [Still Having Issues?](#still-having-issues)

## 💻 System Requirements

- Go version 1.24 or newer
- Git (for version control)
- Internet connection (for downloading dependencies)

## 📥 Installation

### Method 1: Using Go Install (Recommended)

1. Make sure Go is installed on your system:
   ```bash
   go version
   ```
   Ensure the installed version is 1.24 or newer.

2. Install DocuBook CLI globally:
   ```bash
   go install github.com/DocuBook/cli/docubook@latest
   ```

   This will install the binary as `docubook` in `$GOPATH/bin`.

3. Ensure `$GOPATH/bin` is in your PATH:
   - For Zsh (macOS):
     ```bash
     echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
     source ~/.zshrc
     ```
   - For Bash:
     ```bash
     echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bash_profile
     source ~/.bash_profile
     ```

4. Verify the installation was successful:
   ```bash
   docubook --version
   ```

### Method 2: Build from Source

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
   go build -o docubook .
   ```

4. (Optional) Move the binary to your system's PATH

## ⚙️ Configuration

### Adding to PATH

Make sure `$GOPATH/bin` is in your PATH. If not, follow the installation steps above.

### Verifying Installation

To ensure the installation was successful, run:
```bash
docubook --version
```

You should see a version message.

## 🚀 Usage

1. Open a terminal and navigate to your project directory
2. To start a new project, run:
   ```bash
   docubook cli
   ```
3. Follow the on-screen instructions to complete project configuration

Complete example:
```bash
# Create a new project directory
mkdir my-documentation
cd my-documentation

# Initialize DocuBook project
docubook cli
```

## 🔧 Troubleshooting

### Command Not Found

If you see `command not found: docubook`:

1. Check if Go is installed:
   ```bash
   go version
   ```

2. Verify the binary was installed:
   ```bash
   ls -l $(go env GOPATH)/bin/docubook
   ```

3. If the file exists but can't be executed, ensure `$GOPATH/bin` is in your PATH.

### Permission Denied

If you encounter a permission error:
```bash
chmod +x $(go env GOPATH)/bin/docubook
```

### Update to Latest Version

To update to the latest version of DocuBook CLI, run the following command in your terminal:
```bash
go install github.com/DocuBook/cli/docubook@latest
```
This command will update the `docubook` binary in your `$GOPATH/bin` directory.

### Uninstalling DocuBook CLI

To uninstall DocuBook CLI:

1. Remove the binary:
   ```bash
   rm $(go env GOPATH)/bin/docubook
   ```

2. (Optional) If you added the export line to your shell configuration, you can remove it:
   - For Zsh (macOS):
     ```bash
     sed -i '' '/export PATH=\$PATH:$(go env GOPATH)\/bin/d' ~/.zshrc
     ```
   - For Bash:
     ```bash
     sed -i '' '/export PATH=\$PATH:$(go env GOPATH)\/bin/d' ~/.bash_profile
     ```

3. Verify the uninstallation:
   ```bash
   which docubook || echo "DocuBook CLI has been successfully uninstalled"
   ```

### Still Having Issues?

- Check your Go configuration:
  ```bash
  go env
  ```
- Look for any error messages during installation
- Verify you have write permissions to the Go bin directory
- Ensure `$GOPATH` is properly configured
