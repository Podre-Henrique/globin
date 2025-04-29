# Globin

A high-performance URL shortener service written in Go using Fiber framework.

## Features

- Fast URL shortening
- URL validation
- Auto-expiring URLs
- Rate limiting
- Thread-safe operations

## Technical Details

- Built with Go 1.24
- Uses Fiber v3 web framework
- Atomic timestamp updates
- Configurable rate limiting

## API Endpoints

### Shorten URL
```http
POST /
Content-Type: application/json

{
    "original": "https://something.com.cu/awid?p=2&..."
}
```

Response:
```json
{
    "url": <url>
}
```

### Access Shortened URL
```http
GET /<url>
```
Redirects to the original URL if valid and not expired.

## Examples using curl

### Create short URL:
```bash
curl -X POST \
     -H "Content-Type: application/json" \
     -d '{"original":"https://www.pudim.com.br"}' \
     http://localhost:8000/
```
Then copy the shortened URL from the response.
### Access shortened URL:
```bash
curl -L http://localhost:8000/ <-paste_here
```

## Running the Project

1. Clone the repository
2. Install dependencies:
```bash
go mod download
```
3. Run the server:
```bash
go run main.go
```

The server will start at `localhost:8000`