const { scrapeAnimeList } = require('../../scraper/list');
const { createSuccessResponse, createErrorResponse } = require('../../utils/helpers');

/**
 * Get complete anime list
 * @param {Object} req - Express request object
 * @param {Object} res - Express response object
 * @param {Function} next - Express next middleware function
 */
const getAnimeList = async (req, res, next) => {
  try {
    const data = await scrapeAnimeList();
    
    if (!data || data.length === 0) {
      return res.status(404).json(
        createErrorResponse('No anime found in the list', 404)
      );
    }
    
    // Calculate total anime count
    const totalAnime = data.reduce((total, group) => total + group.anime_list.length, 0);
    
    return res.status(200).json(
      createSuccessResponse({
        total_anime: totalAnime,
        groups: data
      }, 'Anime list retrieved successfully')
    );
  } catch (error) {
    console.error(`Error in getAnimeList controller: ${error.message}`);
    return res.status(500).json(
      createErrorResponse(`Failed to retrieve anime list: ${error.message}`, 500)
    );
  }
};

module.exports = {
  getAnimeList
};