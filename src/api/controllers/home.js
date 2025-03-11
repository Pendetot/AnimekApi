const { scrapeAnoboyHome } = require('../../scraper/home');
const { createSuccessResponse, createErrorResponse } = require('../../utils/helpers');

/**
 * Get latest anime from homepage
 * @param {Object} req - Express request object
 * @param {Object} res - Express response object
 * @param {Function} next - Express next middleware function
 */
const getHomeData = async (req, res, next) => {
  try {
    const data = await scrapeAnoboyHome();
    
    if (!data || data.length === 0) {
      return res.status(404).json(
        createErrorResponse('No anime found on homepage', 404)
      );
    }
    
    return res.status(200).json(
      createSuccessResponse(data, 'Latest anime retrieved successfully')
    );
  } catch (error) {
    console.error(`Error in getHomeData controller: ${error.message}`);
    return res.status(500).json(
      createErrorResponse(`Failed to retrieve homepage data: ${error.message}`, 500)
    );
  }
};

module.exports = {
  getHomeData
};