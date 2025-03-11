const express = require('express');
const router = express.Router();
const { getEpisodeList } = require('../controllers/episodes');
const { cacheMiddleware } = require('../middleware/cache');

// Get episode list for an anime
router.get('/', cacheMiddleware(), getEpisodeList);

module.exports = router;