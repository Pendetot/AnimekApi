const express = require('express');
const router = express.Router();
const { getAnimeDetail } = require('../controllers/detail');
const { cacheMiddleware } = require('../middleware/cache');

// Get anime episode details by URL
router.get('/', cacheMiddleware(), getAnimeDetail);

module.exports = router;