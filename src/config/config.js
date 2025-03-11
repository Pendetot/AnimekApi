module.exports = {
  PORT: process.env.PORT || 1408,
  BASE_URL: 'https://ww1.anoboy.app',
  USER_AGENT: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36',
  CACHE_DURATION: 30 * 60, // 30 minutes in seconds
  REQUEST_TIMEOUT: 10000, // 10 seconds in milliseconds
  RATE_LIMIT: {
    WINDOW_MS: 15 * 60 * 1000, // 15 minutes
    MAX_REQUESTS: 100
  }
};