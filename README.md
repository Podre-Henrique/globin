# Globin

A high-performance URL shortener service written in Go using Fiber framework.

## Features

- Fast URL shortening
- URL validation
- Auto-expiring URLs
- Rate limiting

## Technical Details

- Built with Go 1.24
- Uses Fiber v3 web framework
- Atomic timestamp updates
- Configurable rate limiting

## API Endpoints

### Shorten URL
```http
GET /
Content-Type: application/json

{
    "original": "https://your-long-url.com/..."
}
```

Response:
```json
{
    "shortened": "abc123"
}
```

### Access Shortened URL
```http
GET /:shortened
```
Redirects to the original URL if valid and not expired.

## Configuration

Main constants:
- `URL_LIFETIME`: 3 hours
- `GC_INTERVAL`: 30 minutes
- `URL_LEN`: 6 characters
- Default port: 8000

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