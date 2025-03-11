const NodeCache = require('node-cache');
const { CACHE_DURATION } = require('../../config/config');

// Initialize cache
const cache = new NodeCache({ 
  stdTTL: CACHE_DURATION, 
  checkperiod: CACHE_DURATION * 0.2,
  useClones: false 
});

/**
 * Middleware to cache API responses
 * @param {number} duration - Cache duration in seconds (optional)
 * @returns {Function} Express middleware function
 */
const cacheMiddleware = (duration) => {
  return (req, res, next) => {
    // Use custom or default duration
    const ttl = duration || CACHE_DURATION;
    
    // Create a unique key based on the request URL and query params
    const key = req.originalUrl;
    
    // Check if response exists in cache
    const cachedResponse = cache.get(key);
    if (cachedResponse) {
      // Set custom header to indicate cached response
      res.set('X-Cache', 'HIT');
      return res.json(cachedResponse);
    }
    
    // If not in cache, continue with the request
    res.set('X-Cache', 'MISS');
    
    // Store the original res.json function
    const originalJson = res.json;
    
    // Override res.json method to cache the response
    res.json = function(body) {
      // Store response in cache before sending
      cache.set(key, body, ttl);
      
      // Call the original json method
      return originalJson.call(this, body);
    };
    
    next();
  };
};

/**
 * Clear specific cache entry
 * @param {string} key - Cache key to clear
 */
const clearCache = (key) => {
  cache.del(key);
};

/**
 * Clear all cache
 */
const clearAllCache = () => {
  cache.flushAll();
};

module.exports = {
  cacheMiddleware,
  clearCache,
  clearAllCache
};