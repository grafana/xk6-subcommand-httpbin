# xk6-subcommand-httpbin

**httpbin subcommand extension for k6**

A k6 subcommand extension that runs a local [httpbin](https://httpbin.org/) server directly from k6, providing HTTP testing endpoints for your load tests without external dependencies.

```bash
# Start httpbin server on default port (localhost:5454)
./k6 x httpbin

# Start httpbin server on custom address
./k6 x httpbin --bind localhost:8080

# Stop the server with Ctrl-C
```

The server provides various HTTP testing endpoints (GET, POST, status codes, headers, auth, etc.) that you can use in your k6 scripts or test directly with curl.

**Benefits:**
- **Zero-setup testing** - No external services or Docker required
- **Fast local responses** - Test HTTP logic without network latency
- **Offline development** - Work on k6 scripts without internet connectivity
- **Predictable behavior** - Consistent, deterministic responses for debugging
- **Full httpbin API** - Complete implementation with all standard endpoints

## Example Usage

This example demonstrates how to use the httpbin server with a k6 load test. The test exercises various HTTP methods and validates responses against the local httpbin instance.

### Example k6 Script

```javascript file=script.js
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

### Running the Example

You'll need two terminal windows to run this example:

**Terminal 1 - Start the httpbin server:**
```shell
./k6 x httpbin
```

Wait for the server to start (you'll see a message with the server address).

**Terminal 2 - Run the k6 test:**
```shell
./k6 run script.js
```

The test will execute and display results showing successful checks against the httpbin endpoints. Press `Ctrl-C` in Terminal 1 to stop the server when done.

## Motivation

This project serves as an example of extending k6's functionality through the subcommand extension system, showcasing how developers can create custom tools that integrate seamlessly with the k6 CLI.

This extension demonstrates:

- **Subcommand Registration**: How to register a new subcommand with k6 using the extension system
- **HTTP Server Integration**: Running a local httpbin server as part of k6's functionality  
- **Command Line Interface**: Adding custom flags and options to your subcommand

## Download

Building a custom k6 binary with the `xk6-subcommand-httpbin` extension is necessary for its use. You can download pre-built k6 binaries from the [Releases page](https://github.com/grafana/xk6-subcommand-httpbin/releases/).

## Build

Use the [xk6](https://github.com/grafana/xk6) tool to build a custom k6 binary with the `xk6-subcommand-httpbin` extension. Refer to the [xk6 documentation](https://github.com/grafana/xk6) for more information.

```bash
xk6 build --with github.com/grafana/xk6-subcommand-httpbin@latest
```

## Endpoints

### Information & Inspection

Endpoint      | Description
--------------|--------------------------------------------
`/`           | Help page listing all available endpoints
`/ip`         | Returns the client's origin IP address
`/user-agent` | Returns the client's user-agent string
`/headers`    | Returns all request headers as JSON
`/get`        | Returns GET request data

### HTTP Methods

Endpoint  | Description
----------|-----------------------------
`/post`   | Returns POST request data
`/put`    | Returns PUT request data
`/patch`  | Returns PATCH request data
`/delete` | Returns DELETE request data

### Response Formats

Endpoint         | Description
-----------------|----------------------------------
`/encoding/utf8` | Returns a page containing UTF-8 data
`/html`          | Renders an HTML page
`/xml`           | Returns XML content
`/json`          | Returns JSON content

### Compression

Endpoint   | Description
-----------|-------------------------------------------
`/gzip`    | Returns gzip-encoded data
`/deflate` | Returns deflate-encoded data
`/brotli`  | Returns brotli-encoded data (not implemented)

### Status Codes

Endpoint        | Description
----------------|------------------------------------
`/status/:code` | Returns the specified HTTP status code

### Response Headers

Endpoint                    | Description
----------------------------|----------------------------------
`/response-headers?key=val` | Returns response with custom headers

### Redirects

Endpoint                             | Description
-------------------------------------|-------------------------------
`/redirect/:n`                       | 302 redirect n times (max 10)
`/redirect-to?url=foo`               | 302 redirect to specified URL
`/redirect-to?url=foo&status_code=307` | Custom status code redirect
`/relative-redirect/:n`              | 302 relative redirects n times
`/absolute-redirect/:n`              | 302 absolute redirects n times

### Cookies

Endpoint                  | Description
--------------------------|--------------------------------
`/cookies`                | Returns all cookie data
`/cookies/set?name=value` | Sets one or more simple cookies
`/cookies/delete?name`    | Deletes one or more simple cookies

### Authentication

Endpoint                                  | Description
------------------------------------------|-------------------------------------------
`/basic-auth/:user/:passwd`               | Challenges HTTP Basic Auth
`/hidden-basic-auth/:user/:passwd`        | Returns 404 for failed Basic Auth
`/digest-auth/:qop/:user/:passwd/:algorithm` | Challenges HTTP Digest Auth (with algorithm)
`/digest-auth/:qop/:user/:passwd`         | Challenges HTTP Digest Auth

### Dynamic Data

Endpoint                                    | Description
--------------------------------------------|-------------------------------------------
`/stream/:n`                                | Streams min(n, 100) lines of JSON
`/delay/:n`                                 | Delays response for min(n, 10) seconds
`/drip?numbytes=n&duration=s&delay=s&code=code` | Drips data over duration with optional delay
`/range/1024?duration=s&chunk_size=n`       | Streams n bytes with Range header support
`/bytes/:n`                                 | Generates n random bytes (accepts seed param)
`/stream-bytes/:n`                          | Streams n random bytes (accepts seed, chunk_size)

### Caching

Endpoint      | Description
--------------|-------------------------------------------
`/cache`      | Returns 200 or 304 based on cache headers
`/cache/:n`   | Sets Cache-Control header for n seconds
`/etag/:etag` | Handles ETag validation (200/304/412)

### Images

Endpoint      | Description
--------------|--------------------------------------
`/image`      | Returns image based on Accept header
`/image/png`  | Returns a PNG image
`/image/jpeg` | Returns a JPEG image
`/image/webp` | Returns a WEBP image
`/image/svg`  | Returns an SVG image

### Other

Endpoint      | Description
--------------|--------------------------------
`/links/:n`   | Returns page with n HTML links
`/forms/post` | HTML form that submits to /post
`/robots.txt` | Returns robots.txt rules
`/deny`       | Endpoint denied by robots.txt

For detailed documentation, visit https://httpbingo.org/

## Contribute

If you wish to contribute to this project, please start by reading the [Contributing Guidelines](CONTRIBUTING.md).
