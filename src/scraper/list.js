const axios = require('axios');
const cheerio = require('cheerio');
const { BASE_URL, USER_AGENT, REQUEST_TIMEOUT } = require('../config/config');
const { getRandomUserAgent } = require('../utils/helpers');

/**
 * Scrape the overall anime list from the website
 * @returns {Promise<Array>} Array of anime groups organized by first letter
 */
const scrapeAnimeList = async () => {
  try {
    // Set headers to mimic a browser request
    const headers = {
      'User-Agent': getRandomUserAgent() || USER_AGENT
    };
    
    // URL to scrape
    const url = `${BASE_URL}/anime-list/`;
    
    // Send HTTP request
    const response = await axios.get(url, {
      headers,
      timeout: REQUEST_TIMEOUT
    });
    
    // Parse HTML content
    const $ = cheerio.load(response.data);
    
    // Find anime list structure - typically organized by first letter
    const animeGroups = [];
    
    // Find all divs with class "letter-group" or similar
    const letterGroups = $(['div[class*="letter-group"]', 'ul.lcp_catlist']);
    
    // If the standard structure isn't found, try to find all list items
    if (letterGroups.length === 0) {
      const allListItems = $('li');
      const animeList = [];
      
      allListItems.each((index, item) => {
        const link = $(item).find('a');
        if (link.length) {
          const title = link.text().trim();
          const url = link.attr('href');
          
          animeList.push({
            title,
            url
          });
        }
      });
      
      animeGroups.push({
        letter: 'All',
        anime_list: animeList
      });
    } else {
      // Process each letter group
      letterGroups.each((index, group) => {
        // Try to find the letter/heading
        const heading = $(group).find(['h2', 'h3']);
        const letter = heading.length ? heading.text().trim() : "Unknown";
        
        // Find all anime links in this group
        const links = $(group).find('a');
        const animeList = [];
        
        links.each((i, link) => {
          const title = $(link).text().trim();
          const url = $(link).attr('href');
          
          animeList.push({
            title,
            url
          });
        });
        
        if (animeList.length) {
          animeGroups.push({
            letter,
            anime_list: animeList
          });
        }
      });
    }
    
    return animeGroups;
    
  } catch (error) {
    console.error(`Error fetching the anime list page: ${error.message}`);
    throw new Error(`Failed to scrape anime list: ${error.message}`);
  }
};

module.exports = {
  scrapeAnimeList
};