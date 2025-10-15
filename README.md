# xk6-subcommand-httpbin

> Proof of Concept: httpbin subcommand extension for k6

A proof of concept k6 extension that demonstrates how to add custom subcommands to k6. This extension adds a new `httpbin` subcommand that provides a fully functional, local HTTP request and response service right inside your load testing tool.

This project serves as an example of extending k6's functionality through the subcommand extension system, showcasing how developers can create custom tools that integrate seamlessly with the k6 CLI.

## Overview

This extension demonstrates:

- **Subcommand Registration**: How to register a new subcommand with k6 using the extension system
- **HTTP Server Integration**: Running a local httpbin server as part of k6's functionality  
- **Command Line Interface**: Adding custom flags and options to your subcommand

The `httpbin` subcommand provides endpoints for testing HTTP requests, making it useful for:

- Zero-setup, local API testing
- Debugging k6 scripts with predictable responses
- Ensuring consistent behavior in CI/CD pipelines
- Educational purposes to understand HTTP request/response mechanics

## Status

**⚠️ Proof of Concept** - This extension is intended as a demonstration of k6's subcommand extension capabilities.

## Installation

Build k6 with this extension using [xk6](https://github.com/grafana/xk6):

```bash
xk6 build --k6-version subcommand-p1 --with github.com/grafana/xk6-subcommand-httpbin
```

## Usage

Once built, you can use the `httpbin` subcommand:

```bash
# Start httpbin server on default port (localhost:5454)
./k6 httpbin

# Start httpbin server on custom address
./k6 httpbin --bind localhost:8080
```

The server will start and provide various HTTP testing endpoints. You can then use these endpoints in your k6 scripts or test them directly with curl.

### Example k6 Script

```javascript
import http from 'k6/http';
import { check } from 'k6';

export default function () {
  // Test against the local httpbin server
  let response = http.get('http://localhost:5454/get');
  
  check(response, {
    'status is 200': (r) => r.status === 200,
  });
  
  // Test POST endpoint
  response = http.post('http://localhost:5454/post', {
    key: 'value'
  });
  
  check(response, {
    'POST status is 200': (r) => r.status === 200,
  });
}
```

## Available Endpoints

The httpbin server provides numerous endpoints for testing HTTP functionality:

- `/get`, `/post`, `/put`, `/delete`, `/patch` - HTTP method testing
- `/status/:code` - Return specific HTTP status codes
- `/delay/:n` - Add artificial delays to responses
- `/redirect/:n` - Test redirect handling
- `/basic-auth/:user/:passwd` - Test authentication
- `/cookies/*` - Cookie handling endpoints
- `/gzip`, `/deflate` - Compression testing
- `/stream/:n` - Streaming responses
- And many more...

See the built-in help for a complete list:

```bash
./k6 httpbin --help
```

## Extension Development

This project demonstrates key concepts for k6 subcommand extensions:

### 1. Extension Registration

```go
func init() {
    subcommand.RegisterExtension("httpbin", newCommand)
}
```

### 2. Command Creation

```go
func newCommand(gs *state.GlobalState) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "httpbin",
        Short: "A HTTP server for testing purposes",
        Long:  help,
    }
    // Add flags, run logic, etc.
    return cmd
}
```

## Contributing

This is a proof of concept project. Feel free to use it as a reference for creating your own k6 subcommand extensions.
