const { getAnimeScheduleWithImages } = require('../../scraper/schedule');
const { createSuccessResponse, createErrorResponse } = require('../../utils/helpers');

/**
 * Get anime broadcast schedule with images
 * @param {Object} req - Express request object
 * @param {Object} res - Express response object
 * @param {Function} next - Express next middleware function
 */
const getSchedule = async (req, res, next) => {
  try {
    const data = await getAnimeScheduleWithImages();
    
    if (!data) {
      return res.status(404).json(
        createErrorResponse('Schedule data not found', 404)
      );
    }
    
    // Calculate counts for reporting
    const dayCount = data.schedule_by_day ? data.schedule_by_day.length : 0;
    const animeCount = data.schedule_by_day ? 
      data.schedule_by_day.reduce((total, day) => total + day.anime_list.length, 0) : 0;
    const jadwalCount = data.anoboy_jadwal ? data.anoboy_jadwal.length : 0;
    
    return res.status(200).json(
      createSuccessResponse({
        days_count: dayCount,
        anime_count: animeCount,
        jadwal_count: jadwalCount,
        data: data
      }, 'Schedule data with images retrieved successfully')
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