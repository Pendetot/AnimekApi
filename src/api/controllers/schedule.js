const { scrapeJadwal } = require('../../scraper/schedule');
const { createSuccessResponse, createErrorResponse } = require('../../utils/helpers');

/**
 * Get anime broadcast schedule
 * @param {Object} req - Express request object
 * @param {Object} res - Express response object
 * @param {Function} next - Express next middleware function
 */
const getSchedule = async (req, res, next) => {
  try {
    const data = await scrapeJadwal();
    
    if (!data) {
      return res.status(404).json(
        createErrorResponse('Schedule data not found', 404)
      );
    }
    
    // Calculate counts for reporting
    const seasonCount = data.seasons ? data.seasons.length : 0;
    const seasonalAnimeCount = data.seasons ? 
      data.seasons.reduce((total, season) => total + season.anime_list.length, 0) : 0;
    const scheduleCount = data.anoboy_jadwal ? data.anoboy_jadwal.length : 0;
    
    return res.status(200).json(
      createSuccessResponse({
        season_count: seasonCount,
        seasonal_anime_count: seasonalAnimeCount,
        schedule_count: scheduleCount,
        data: data
      }, 'Schedule data retrieved successfully')
    );
  } catch (error) {
    console.error(`Error in getSchedule controller: ${error.message}`);
    return res.status(500).json(
      createErrorResponse(`Failed to retrieve schedule data: ${error.message}`, 500)
    );
  }
};

module.exports = {
  getSchedule
};