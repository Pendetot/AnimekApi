const express = require('express');
const router = express.Router();
const { getHomeData } = require('../controllers/home');
const { cacheMiddleware } = require('../middleware/cache');

// Get latest anime from homepage
router.get('/', cacheMiddleware(), getHomeData);

module.exports = router;