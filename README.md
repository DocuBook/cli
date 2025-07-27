# DocuBook CLI

DocuBook CLI is a command-line tool for creating and managing beautiful documentation sites.

## Installation

You can install the DocuBook CLI using Go's `install` command:

   ```bash
   go install github.com/DocuBook/cli@latest
   ```

2. Build the binary:
   ```bash
   go build -o docubook
   ```

3. Run the CLI:
   ```bash
   ./docubook cli
   ```

4. Navigate to your desired directory:
   ```bash
   cd your-project-directory
   ```

5. Create a new project:
   ```bash
   create
   ```

## Requirements

- Go 1.24 or later
- Git (for version control)

## Development

1. Clone the repository:
   ```bash
   git clone https://github.com/DocuBook/cli.git
   cd cli
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Build and run:
   ```bash
   go build -o docubook
   ./docubook cli
   ```

## License

[MIT](LICENSE)
