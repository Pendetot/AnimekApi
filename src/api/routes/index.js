const express = require('express');
const router = express.Router();

const homeRoutes = require('./home');
const detailRoutes = require('./detail');
const listRoutes = require('./list');
const episodesRoutes = require('./episodes');
const scheduleRoutes = require('./schedule');

// Info route
router.get('/', (req, res) => {
  res.json({
    status: 'success',
    message: 'AnoBoy API is running',
    version: '1.0.0',
    endpoints: [
      { path: '/api/home', description: 'Get latest anime from homepage' },
      { path: '/api/detail?slug=<anime_slug>', description: 'Get anime episode details by slug' },
      { path: '/api/list', description: 'Get complete anime list' },
      { path: '/api/episodes?title=<anime_title>', description: 'Get episodes list for an anime by title' },
      { path: '/api/schedule', description: 'Get anime broadcast schedule' }
    ]
  });
});

// Register routes
router.use('/home', homeRoutes);
router.use('/detail', detailRoutes);
router.use('/list', listRoutes);
router.use('/episodes', episodesRoutes);
router.use('/schedule', scheduleRoutes);

module.exports = router;