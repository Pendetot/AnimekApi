const express = require('express');
const router = express.Router();
const { getSchedule } = require('../controllers/schedule');
const { cacheMiddleware } = require('../middleware/cache');

// Get anime broadcast schedule
// Cache this for longer as schedule doesn't change frequently
router.get('/', cacheMiddleware(3 * 60 * 60), getSchedule); // 3 hour cache

module.exports = router;