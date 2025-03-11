# AnoBoy API

A Node.js API for scraping and serving anime data from AnoBoy website.

## Features

- Fetches anime episode details
- Retrieves complete anime list
- Gets episode list for specific anime
- Fetches anime broadcast schedule with images
- Response caching to reduce website load and improve performance
- Rate limiting to prevent abuse

## Getting Started

### Option 1: Clone from GitHub

```bash
git clone https://github.com/Pendetot/AnimekApi.git
cd AnimekApi
npm install
npm start
```

For development with auto-restart:
```bash
npm run dev
```

## API Endpoints

| Endpoint | Description | Query Parameters |
|----------|-------------|------------------|
| `GET /api` | API information and available endpoints | None |
| `GET /api/detail` | Get anime episode details | `slug` or `path` (required): Anime episode slug or path |
| `GET /api/list` | Get complete anime list | None |
| `GET /api/episodes` | Get episode list for an anime | `title`, `slug`, or `path` (required): Anime title, slug or path |
| `GET /api/schedule` | Get anime broadcast schedule with images | None |

## Usage Examples

### Get Anime Episode Details
```
GET http://localhost:1408/api/detail?slug=2025/03/a-war-between-humans-and-ai-episode-10
```

### Get Complete Anime List
```
GET http://localhost:1408/api/list
```

### Get Episode List for an Anime
```
GET http://localhost:1408/api/episodes?title=a-war-between-humans-and-ai
```
or
```
GET http://localhost:1408/api/episodes?slug=anime/a-war-between-humans-and-ai
```

### Get Anime Broadcast Schedule with Images
```
GET http://localhost:1408/api/schedule
```

## Cache Control

The API implements caching to improve performance and reduce load on the AnoBoy website. Cache duration varies by endpoint:

- `/api/detail`: 30 minutes (default)
- `/api/list`: 1 hour
- `/api/episodes`: 30 minutes (default)
- `/api/schedule`: 3 hours

## Rate Limiting

To prevent abuse, the API has rate limiting in place:
- 100 requests per 15 minutes per IP address

## Dependencies

- express: Web server framework
- axios: HTTP client for making requests
- cheerio: HTML parsing and manipulation
- node-cache: Caching implementation
- cors: Cross-Origin Resource Sharing
- helmet: Security headers
- morgan: HTTP request logger
- express-rate-limit: Rate limiting middleware
- dotenv: Environment variables management
- user-agents: Random user agent generator

## Disclaimer

This API is for educational purposes only. Please respect the AnoBoy website's terms of service and robots.txt file. Consider adding appropriate delays between requests and do not overload their servers.

## License

MIT