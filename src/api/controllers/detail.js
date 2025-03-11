const { scrapeAnimeDetail } = require('../../scraper/detail');
const { createSuccessResponse, createErrorResponse } = require('../../utils/helpers');
const { BASE_URL } = require('../../config/config');

/**
 * Get anime episode details by slug or path
 * @param {Object} req - Express request object
 * @param {Object} res - Express response object
 * @param {Function} next - Express next middleware function
 */
const getAnimeDetail = async (req, res, next) => {
  try {
    // Get the slug/path from the query or params
    const { slug, path } = req.query;
    let animeUrl;
    
    if (slug) {
      // If slug is provided, construct the URL
      animeUrl = `${BASE_URL}/${slug}`;
    } else if (path) {
      // If path is provided, construct the URL
      animeUrl = `${BASE_URL}/${path}`;
    } else {
      return res.status(400).json(
        createErrorResponse('Either slug or path parameter is required', 400)
      );
    }
    
    const data = await scrapeAnimeDetail(animeUrl);
    
    if (!data) {
      return res.status(404).json(
        createErrorResponse('Anime details not found', 404)
      );
    }
    
    return res.status(200).json(
      createSuccessResponse(data, 'Anime details retrieved successfully')
    );
  } catch (error) {
    console.error(`Error in getAnimeDetail controller: ${error.message}`);
    return res.status(500).json(
      createErrorResponse(`Failed to retrieve anime details: ${error.message}`, 500)
    );
  }
};

module.exports = {
  getAnimeDetail
};