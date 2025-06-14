# AnoBoy API

API Golang untuk scraping dan menyajikan data anime dari website AnoBoy.

## Fitur

- Mengambil detail episode anime
- Mengambil daftar anime lengkap
- Mendapatkan daftar episode untuk anime tertentu
- Mengambil jadwal tayang anime dengan gambar
- Cache response untuk mengurangi beban website dan meningkatkan performa
- Rate limiting untuk mencegah penyalahgunaan

## Memulai

### Clone dari GitHub

```bash
git clone https://github.com/Pendetot/AnimekApi.git
cd AnimekApi
go mod tidy
go run main.go
```

Untuk development dengan auto-restart:
```bash
go install github.com/cosmtrek/air@latest
air
```

## Endpoint API

| Endpoint | Deskripsi | Parameter Query |
|----------|-----------|-----------------|
| `GET /api` | Informasi API dan endpoint yang tersedia | Tidak ada |
| `GET /api/detail` | Mendapatkan detail episode anime | `slug` atau `path` (wajib): Slug atau path episode anime |
| `GET /api/list` | Mendapatkan daftar anime lengkap | Tidak ada |
| `GET /api/episodes` | Mendapatkan daftar episode untuk anime | `title`, `slug`, atau `path` (wajib): Judul, slug atau path anime |
| `GET /api/schedule` | Mendapatkan jadwal tayang anime dengan gambar | Tidak ada |

## Contoh Penggunaan

### Mendapatkan Detail Episode Anime
```
GET http://localhost:1408/api/detail?slug=2025/03/a-war-between-humans-and-ai-episode-10
```

### Mendapatkan Daftar Anime Lengkap
```
GET http://localhost:1408/api/list
```

### Mendapatkan Daftar Episode untuk Anime
```
GET http://localhost:1408/api/episodes?title=a-war-between-humans-and-ai
```
atau
```
GET http://localhost:1408/api/episodes?slug=anime/a-war-between-humans-and-ai
```

### Mendapatkan Jadwal Tayang Anime dengan Gambar
```
GET http://localhost:1408/api/schedule
```

## Kontrol Cache

API mengimplementasikan caching untuk meningkatkan performa dan mengurangi beban pada website AnoBoy. Durasi cache bervariasi berdasarkan endpoint:

- `/api/detail`: 30 menit (default)
- `/api/list`: 1 jam
- `/api/episodes`: 30 menit (default)
- `/api/schedule`: 3 jam

## Rate Limiting

Untuk mencegah penyalahgunaan, API memiliki rate limiting:
- 100 request per 15 menit per alamat IP

## Dependencies

- gin-gonic/gin: Framework web server
- PuerkitoBio/goquery: HTML parsing dan manipulasi
- gin-contrib/cors: Cross-Origin Resource Sharing
- gin-contrib/cache: Implementasi caching

## Struktur Proyek

```
AnimekApi/
├── main.go              # Entry point aplikasi
├── config/
│   └── config.go        # Konfigurasi aplikasi
├── handlers/
│   └── handlers.go      # Handler API
├── middleware/
│   └── ratelimit.go     # Middleware rate limiting
├── scraper/
│   ├── detail.go        # Scraper detail anime
│   ├── list.go          # Scraper daftar anime
│   ├── episodes.go      # Scraper episode
│   └── schedule.go      # Scraper jadwal
├── utils/
│   └── utils.go         # Utility functions
├── go.mod               # Go module definition
└── README.md            # Dokumentasi
```

## Variabel Environment

Anda dapat mengatur konfigurasi menggunakan variabel environment:

- `PORT`: Port server (default: 1408)
- `ENVIRONMENT`: Environment mode (default: development)
- `BASE_URL`: URL base AnoBoy (default: https://ww1.anoboy.app)
- `USER_AGENT`: User agent untuk request (default: Mozilla/5.0...)
- `CACHE_DURATION`: Durasi cache dalam menit (default: 30)
- `REQUEST_TIMEOUT`: Timeout request dalam detik (default: 10)
- `RATE_LIMIT_WINDOW_MS`: Window rate limit dalam menit (default: 15)
- `RATE_LIMIT_MAX_REQUESTS`: Maksimal request per window (default: 100)

## Build dan Deploy

### Build untuk Production
```bash
go build -o animek-api main.go
./animek-api
```

### Build untuk berbagai platform
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o animek-api-linux main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o animek-api-windows.exe main.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o animek-api-macos main.go
```

### Docker
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 1408
CMD ["./main"]
```

## Disclaimer

API ini hanya untuk tujuan edukasi. Harap hormati terms of service dan robots.txt dari website AnoBoy. Pertimbangkan untuk menambahkan delay yang sesuai antara request dan jangan membebani server mereka.

## Kontribusi

1. Fork repository ini
2. Buat branch fitur (`git checkout -b feature/AmazingFeature`)
3. Commit perubahan (`git commit -m 'Add some AmazingFeature'`)
4. Push ke branch (`git push origin feature/AmazingFeature`)
5. Buka Pull Request

## Lisensi

MIT

## Changelog

### v2.0.0 (Golang Version)
- Migrasi dari Node.js ke Golang
- Peningkatan performa dan efisiensi memori
- Implementasi rate limiting yang lebih baik
- Struktur kode yang lebih terorganisir
- Support untuk concurrent request handling
- Dokumentasi dalam bahasa Indonesia

### v1.0.0 (Node.js Version)
- Implementasi awal dengan Node.js
- Fitur dasar scraping anime
- Cache dan rate limiting sederhana
