const axios = require('axios');
const cheerio = require('cheerio');
const { BASE_URL, USER_AGENT, REQUEST_TIMEOUT } = require('../config/config');
const { getRandomUserAgent } = require('../utils/helpers');

/**
 * Scrape image from an anime page
 * @param {string} url - URL of the anime page
 * @returns {Promise<string|null>} URL of the anime image or null if not found
 */
const scrapeAnimeImage = async (url) => {
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
    
    // Look for the main image
    const mainImage = $('amp-img').first();
    if (mainImage.length && mainImage.attr('src')) {
      return mainImage.attr('src');
    }
    
    return null;
  } catch (error) {
    console.error(`Error fetching anime image from ${url}: ${error.message}`);
    return null;
  }
};

/**
 * Extract links and anime data from schedule page
 * @param {string} schedulePageHtml - HTML content of the schedule page
 * @returns {Object} Extracted schedule data with anime links
 */
const extractScheduleData = ($) => {
  // Find all day sections (h1 in div.unduhan)
  const days = [];
  const dayDivs = $('div.unduhan').has('h1');
  
  dayDivs.each((index, dayDiv) => {
    const dayName = $(dayDiv).find('h1').text().trim();
    
    // Find all anime entries (li elements) under this day
    const animeEntries = [];
    const animeItems = $(dayDiv).find('ul.lcp_catlist li');
    
    animeItems.each((i, item) => {
      const animeLink = $(item).find('a');
      if (animeLink.length) {
        const title = animeLink.text().trim();
        const url = animeLink.attr('href');
        
        // Extract broadcast time from the text after the link
        const fullText = $(item).text().trim();
        const timeMatch = fullText.match(/([^,]+,\s*\d+:\d+)$/);
        let broadcastTime = null;
        
        if (timeMatch) {
          broadcastTime = timeMatch[1].trim();
        }
        
        animeEntries.push({
          title,
          url,
          broadcast_time: broadcastTime,
          image_url: null // Will be filled later
        });
      }
    });
    
    if (animeEntries.length) {
      days.push({
        day: dayName,
        anime_list: animeEntries
      });
    }
  });
  
  return days;
};

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
    
    // URL to scrape (changed to the schedule page from your example)
    const url = `${BASE_URL}/2015/05/anime-subtitle-indonesia-ini-adalah-arsip-file-kami/`;
    
    // Send HTTP request
    const response = await axios.get(url, {
      headers,
      timeout: REQUEST_TIMEOUT
    });
    
    // Parse HTML content
    const $ = cheerio.load(response.data);
    
    // Extract schedule data
    const scheduleData = extractScheduleData($);
    
    // Fetch images for anime (limit to first 3 anime per day to avoid too many requests)
    for (const day of scheduleData) {
      const animesToFetch = day.anime_list.slice(0, 3); // Limit to first 3 anime per day
      
      for (const anime of animesToFetch) {
        if (anime.url) {
          const imageUrl = await scrapeAnimeImage(anime.url);
          
          if (imageUrl) {
            // Update all anime with this title to have the same image
            day.anime_list.forEach(entry => {
              if (entry.title === anime.title) {
                entry.image_url = imageUrl;
              }
            });
          }
        }
      }
    }
    
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
          const url = titleElem.length ? titleElem.attr('href') : null;
          const day = cols.length > 1 ? $(cols[1]).text().trim() : "";
          const time = cols.length > 2 ? $(cols[2]).text().trim() : "";
          
          anoboyJadwal.push({
            title,
            url,
            day,
            time,
            image_url: null // Will be filled later if we decide to scrape images for these too
          });
        }
      });
    }
    
    return {
      schedule_by_day: scheduleData,
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