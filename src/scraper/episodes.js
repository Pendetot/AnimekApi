const axios = require('axios');
const cheerio = require('cheerio');
const { USER_AGENT, REQUEST_TIMEOUT } = require('../config/config');
const { getRandomUserAgent, extractEpisodeNumber } = require('../utils/helpers');

/**
 * Scrape the episode list from an anime series page
 * @param {string} url - The URL of the anime series page
 * @returns {Promise<Object>} Object containing anime title, episodes, metadata, and description
 */
const scrapeEpisodeList = async (url) => {
  try {
    // Set headers to mimic a browser request
    const headers = {
      'User-Agent': getRandomUserAgent() || USER_AGENT
    };
    
    // Send HTTP request
    const response = await axios.get(url, {
      headers,
      timeout: REQUEST_TIMEOUT
    });
    
    // Parse HTML content
    const $ = cheerio.load(response.data);
    
    // Get anime title
    const titleElement = $('h1');
    const title = titleElement.length ? titleElement.text().trim() : "Title not found";
    
    // Find episode list - check different possible containers
    let episodeSection = null;
    
    // Try different possible containers for episode lists
    const possibleContainers = [
      $('div.singlelink'),
      $('div.hq').next('div'),
      $('ul.lcp_catlist'),
      $('div.episodes')
    ];
    
    for (const container of possibleContainers) {
      if (container && container.length) {
        episodeSection = container;
        break;
      }
    }
    
    const episodes = [];
    
    if (episodeSection && episodeSection.length) {
      // Try to find episodes in list items
      const episodeItems = episodeSection.find('li');
      
      if (episodeItems.length) {
        episodeItems.each((index, item) => {
          const link = $(item).find('a');
          if (link.length) {
            const episodeTitle = link.text().trim();
            const episodeUrl = link.attr('href');
            
            // Try to extract episode number from title using regex
            const episodeNumber = extractEpisodeNumber(episodeTitle);
            
            episodes.push({
              title: episodeTitle,
              number: episodeNumber,
              url: episodeUrl
            });
          }
        });
      } else {
        // If no list items found, try to find direct links
        const links = episodeSection.find('a');
        links.each((i, link) => {
          const episodeTitle = $(link).text().trim();
          const episodeUrl = $(link).attr('href');
          
          // Try to extract episode number from title using regex
          const episodeNumber = extractEpisodeNumber(episodeTitle);
          
          episodes.push({
            title: episodeTitle,
            number: episodeNumber,
            url: episodeUrl
          });
        });
      }
    }
    
    // Sort episodes by number if possible
    episodes.sort((a, b) => {
      // Parse episode numbers to integers for comparison
      const numA = a.number === "Unknown" ? Infinity : parseInt(a.number);
      const numB = b.number === "Unknown" ? Infinity : parseInt(b.number);
      
      // Handle non-numeric values
      if (isNaN(numA) && isNaN(numB)) return 0;
      if (isNaN(numA)) return 1;
      if (isNaN(numB)) return -1;
      
      return numA - numB;
    });
    
    // Get anime metadata
    const metadata = {};
    const infoTable = $('table');
    if (infoTable.length) {
      const rows = infoTable.find('tr');
      rows.each((index, row) => {
        const header = $(row).find('th');
        const data = $(row).find('td');
        if (header.length && data.length) {
          const key = header.text().trim();
          const value = data.text().trim();
          metadata[key] = value;
        }
      });
    }
    
    // Get anime description
    let description = "";
    const descElem = $('div.unduhan');
    if (descElem.length) {
      description = descElem.text().trim();
    }
    
    return {
      title,
      episodes,
      metadata,
      description
    };
    
  } catch (error) {
    console.error(`Error fetching the episode list: ${error.message}`);
    throw new Error(`Failed to scrape episode list: ${error.message}`);
  }
};

module.exports = {
  scrapeEpisodeList
};