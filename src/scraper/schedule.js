const axios = require('axios');
const cheerio = require('cheerio');
const { BASE_URL, USER_AGENT, REQUEST_TIMEOUT } = require('../config/config');
const { getRandomUserAgent } = require('../utils/helpers');

/**
 * Scrape the anime broadcast schedule from jadwal.php page
 * @returns {Promise<Object>} Object containing seasonal anime and anoboy internal schedule
 */
const scrapeJadwal = async () => {
  try {
    // Set headers to mimic a browser request
    const headers = {
      'User-Agent': getRandomUserAgent() || USER_AGENT
    };
    
    // URL to scrape
    const url = `${BASE_URL}/uploads/jadwal.php`;
    
    // Send HTTP request
    const response = await axios.get(url, {
      headers,
      timeout: REQUEST_TIMEOUT
    });
    
    // Parse HTML content
    const $ = cheerio.load(response.data);
    
    // Find all anime seasons/categories
    const seasons = [];
    const animeHeaders = $('div.anime-header');
    
    animeHeaders.each((index, header) => {
      const seasonName = $(header).text().trim();
      
      // Find all anime entries that follow this header until the next header
      const animeEntries = [];
      let current = $(header).next();
      
      while (current.length && !current.hasClass('anime-header')) {
        if (current.find('div.title').length) {
          // Extract anime information
          const animeData = {};
          
          // Title
          const titleElem = current.find('h2.h2_anime_title');
          if (titleElem.length && titleElem.find('a').length) {
            animeData.title = titleElem.find('a').text().trim();
          }
          
          // Broadcast info
          const infoElem = current.find('div.info');
          if (infoElem.length) {
            const items = infoElem.find('span.item');
            if (items.length > 0) {
              animeData.broadcast_time = $(items[0]).text().trim();
            }
            
            if (items.length > 1) {
              const details = $(items[1]).text().trim().replace(/\n/g, ' ').replace(/  +/g, ' ');
              animeData.details = details;
            }
          }
          
          // Genres
          const genresElem = current.find('div.genres-inner');
          if (genresElem.length) {
            const genreLinks = genresElem.find('a');
            animeData.genres = [];
            genreLinks.each((i, genre) => {
              animeData.genres.push($(genre).text().trim());
            });
          }
          
          if (Object.keys(animeData).length) {
            animeEntries.push(animeData);
          }
        }
        
        // Move to the next element
        current = current.next();
      }
      
      if (animeEntries.length) {
        seasons.push({
          season_name: seasonName,
          anime_list: animeEntries
        });
      }
    });
    
    // Also scrape the anoBoy internal schedule table if available
    const anoboyJadwal = [];
    const jadwalTable = $('table').first();
    
    if (jadwalTable.length) {
      const rows = jadwalTable.find('tr');
      
      // Skip header row
      rows.slice(1).each((index, row) => {
        const cols = $(row).find('td');
        if (cols.length >= 3) {
          const titleElem = $(cols[0]).find('a');
          const title = titleElem.length ? titleElem.text().trim() : $(cols[0]).text().trim();
          const day = cols.length > 1 ? $(cols[1]).text().trim() : "";
          const time = cols.length > 2 ? $(cols[2]).text().trim() : "";
          
          anoboyJadwal.push({
            title,
            day,
            time
          });
        }
      });
    }
    
    return {
      seasons,
      anoboy_jadwal: anoboyJadwal
    };
    
  } catch (error) {
    console.error(`Error fetching the schedule page: ${error.message}`);
    throw new Error(`Failed to scrape schedule: ${error.message}`);
  }
};

module.exports = {
  scrapeJadwal
};