Start a local HTTP test server for k6 testing and development purposes.
  
The httpbin server provides a comprehensive set of endpoints for testing HTTP clients, including various response types, authentication methods, redirects, and streaming capabilities.

## AVAILABLE ENDPOINTS

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

Endpoint                           | Description
-----------------------------------|-------------------------------------------
`/basic-auth/:user/:passwd`        | Challenges HTTP Basic Auth
`/hidden-basic-auth/:user/:passwd` | Returns 404 for failed Basic Auth
`/digest-auth/:qop/:user/:passwd`  | Challenges HTTP Digest Auth
`/digest-auth/:qop/:user/:passwd/:algorithm` | Challenges HTTP Digest Auth (with algorithm)

### Dynamic Data

Endpoint                              | Description
--------------------------------------|-------------------------------------------
`/stream/:n`                          | Streams min(n, 100) lines of JSON
`/delay/:n`                           | Delays response for min(n, 10) seconds
`/bytes/:n`                           | Generates n random bytes (accepts seed param)
`/stream-bytes/:n`                    | Streams n random bytes (accepts seed, chunk_size)
`/range/1024?duration=s&chunk_size=n` | Streams n bytes with Range header support
`/drip?numbytes=n&duration=s&delay=s&code=code` | Drips data over duration with optional delay

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
