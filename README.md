# go-api-contract

A fast API contract validation and regression testing tool that validates live API responses against OpenAPI/Swagger specifications.

## Features

- 🔍 **OpenAPI/Swagger Support**: Read and validate against OpenAPI 3.0+ specifications
- ✅ **Schema Validation**: Validate API responses against defined schemas
- 🌐 **HTTP Client**: Built-in HTTP client using Go's `net/http`
- 🔐 **Environment Variables**: Support for environment variable substitution and secret injection
- 📊 **Multiple Output Formats**: Colored CLI output or machine-readable JSON
- 🎯 **Diff Reporting**: Detailed error reporting with validation differences
- 🚀 **Fast & Lightweight**: Written in Go for maximum performance

## Installation

### From Source

```bash
git clone https://github.com/BaseMax/go-api-contract.git
cd go-api-contract
go build -o apicontract ./cmd/apicontract
```

### Using Go Install

```bash
go install github.com/BaseMax/go-api-contract/cmd/apicontract@latest
```

## Usage

### Basic Usage

```bash
apicontract -spec openapi.yaml -endpoint https://api.example.com/users/1 -method GET
```

### With Request Body (POST/PUT/PATCH)

```bash
apicontract -spec openapi.yaml \
  -endpoint https://api.example.com/posts \
  -method POST \
  -body '{"title":"Test","body":"Content","userId":1}'
```

### With Headers

```bash
apicontract -spec openapi.yaml \
  -endpoint https://api.example.com/protected \
  -method GET \
  -headers '{"Authorization":"Bearer token123","Content-Type":"application/json"}'
```

### JSON Output

```bash
apicontract -spec openapi.yaml \
  -endpoint https://api.example.com/users/1 \
  -method GET \
  -json
```

### Verbose Mode

```bash
apicontract -spec openapi.yaml \
  -endpoint https://api.example.com/users/1 \
  -method GET \
  -verbose
```

## Environment Variables

The tool supports environment variable substitution in URLs, headers, and request bodies:

```bash
# Set environment variables
export API_BASE_URL=https://api.example.com
export API_TOKEN=your_secret_token

# Use in commands
apicontract -spec openapi.yaml \
  -endpoint '${API_BASE_URL}/users/1' \
  -headers '{"Authorization":"Bearer ${API_TOKEN}"}'

# Alternative syntax
apicontract -spec openapi.yaml \
  -endpoint '$API_BASE_URL/users/1' \
  -headers '{"Authorization":"Bearer $API_TOKEN"}'
```

Both `${VAR_NAME}` and `$VAR_NAME` formats are supported.

## Command-Line Flags

| Flag | Description | Required | Default |
|------|-------------|----------|---------|
| `-spec` | Path to OpenAPI/Swagger specification file | Yes | - |
| `-endpoint` | API endpoint URL to validate | Yes | - |
| `-method` | HTTP method to use | No | `GET` |
| `-body` | Request body (for POST/PUT/PATCH) | No | - |
| `-headers` | Request headers as JSON object | No | - |
| `-json` | Output results as JSON | No | `false` |
| `-verbose` | Enable verbose output | No | `false` |
| `-version` | Print version information | No | - |

## Example OpenAPI Specification

See [examples/openapi.yaml](examples/openapi.yaml) for a complete example.

```yaml
openapi: 3.0.0
info:
  title: Sample API
  version: 1.0.0
paths:
  /users/{id}:
    get:
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      responses:
        '200':
          description: Success
          content:
            application/json:
              schema:
                type: object
                properties:
                  id:
                    type: integer
                  name:
                    type: string
                required:
                  - id
                  - name
```

## Output Examples

### Colored CLI Output (Default)

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  API Contract Validation Report
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✓ All validations passed

Timestamp: 2024-01-15T10:30:00Z
Status Code: 200
Errors Found: 0
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### JSON Output

```json
{
  "valid": true,
  "statusCode": 200,
  "responseBody": {
    "id": 1,
    "name": "John Doe"
  },
  "timestamp": "2024-01-15T10:30:00Z",
  "errorCount": 0
}
```

## Testing with JSONPlaceholder

The included example uses the JSONPlaceholder API for testing:

```bash
# Build the tool
go build -o apicontract ./cmd/apicontract

# Test GET request
./apicontract -spec examples/openapi.yaml \
  -endpoint https://jsonplaceholder.typicode.com/posts/1 \
  -method GET

# Test POST request
./apicontract -spec examples/openapi.yaml \
  -endpoint https://jsonplaceholder.typicode.com/posts \
  -method POST \
  -body '{"title":"foo","body":"bar","userId":1}'
```

## Exit Codes

- `0`: Validation passed successfully
- `1`: Validation failed or error occurred

## Libraries Used

- [kin-openapi](https://github.com/getkin/kin-openapi) - OpenAPI 3.0 implementation for Go
- [fatih/color](https://github.com/fatih/color) - Colored terminal output
- [gorilla/mux](https://github.com/gorilla/mux) - HTTP routing (via kin-openapi)
- Standard Go `net/http` - HTTP client

## Project Structure

```
.
├── cmd/
│   └── apicontract/         # CLI entry point
│       └── main.go
├── pkg/
│   ├── client/              # HTTP client with env var support
│   │   └── client.go
│   ├── loader/              # OpenAPI spec loader
│   │   └── loader.go
│   ├── validator/           # Schema validation logic
│   │   └── validator.go
│   ├── reporter/            # Report generation
│   │   └── reporter.go
│   └── output/              # Output formatting
│       └── output.go
├── examples/
│   └── openapi.yaml         # Sample OpenAPI spec
├── go.mod
├── go.sum
└── README.md
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the GPL-3.0 License - see the [LICENSE](LICENSE) file for details.

## Author

Copyright (c) 2025 Max Base
