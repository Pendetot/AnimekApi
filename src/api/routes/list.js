const express = require('express');
const router = express.Router();
const { getAnimeList } = require('../controllers/list');
const { cacheMiddleware } = require('../middleware/cache');

// Get complete anime list
// Cache this for longer as it doesn't change as frequently
router.get('/', cacheMiddleware(60 * 60), getAnimeList); // 1 hour cache

module.exports = router;