const { scrapeEpisodeList } = require('../../scraper/episodes');
const { createSuccessResponse, createErrorResponse } = require('../../utils/helpers');
const { BASE_URL } = require('../../config/config');

/**
 * Get episode list for an anime
 * @param {Object} req - Express request object
 * @param {Object} res - Express response object
 * @param {Function} next - Express next middleware function
 */
const getEpisodeList = async (req, res, next) => {
  try {
    // Get the anime title or slug from the query or params
    const { title, slug, path } = req.query;
    let animeUrl;
    
    if (title) {
      // If title is provided, construct the URL (may require additional logic to find correct URL)
      // This is a simplified approach - you might need to search for the anime first
      const cleanTitle = title.toLowerCase().replace(/\s+/g, '-');
      animeUrl = `${BASE_URL}/${cleanTitle}`;
    } else if (slug) {
      // If slug is provided, construct the URL
      animeUrl = `${BASE_URL}/${slug}`;
    } else if (path) {
      // If path is provided, construct the URL
      animeUrl = `${BASE_URL}/${path}`;
    } else {
      return res.status(400).json(
        createErrorResponse('Either title, slug, or path parameter is required', 400)
      );
    }
    
    const data = await scrapeEpisodeList(animeUrl);
    
    if (!data) {
      return res.status(404).json(
        createErrorResponse('Episode list not found', 404)
      );
    }
    
    return res.status(200).json(
      createSuccessResponse({
        title: data.title,
        total_episodes: data.episodes.length,
        description: data.description,
        metadata: data.metadata,
        episodes: data.episodes
      }, 'Episode list retrieved successfully')
    );
  } catch (error) {
    console.error(`Error in getEpisodeList controller: ${error.message}`);
    return res.status(500).json(
      createErrorResponse(`Failed to retrieve episode list: ${error.message}`, 500)
    );
  }
};

module.exports = {
  getEpisodeList
};