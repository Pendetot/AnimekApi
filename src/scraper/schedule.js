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
      timeout: REQUEST_TIMEOUT * 2 // Double timeout for image scraping
    });
    
    // Parse HTML content
    const $ = cheerio.load(response.data);
    
    // Try multiple selectors in order of specificity
    const selectors = [
      // Main content image - specific to the layout you showed
      '.column-three-fourth amp-img',
      // Featured image that might be in a different position
      'amp-img[width="640"][height="360"]',
      // Try to find image in the content area
      '.column-three-fourth img', 
      // Try any amp-img with dimensions that suggest it's a content image
      'amp-img[width][height]',
      // Default fallback - any amp-img
      'amp-img',
      // Last resort - any img
      'img'
    ];
    
    // Try each selector
    for (const selector of selectors) {
      const element = $(selector).first();
      if (element.length && element.attr('src')) {
        const src = element.attr('src');
        // Make sure the URL is complete
        if (src.startsWith('http')) {
          return src;
        } else if (src.startsWith('/')) {
          // If it's a relative URL, make it absolute
          return `${BASE_URL}${src}`;
        }
      }
    }
    
    // If we couldn't find an image with the selectors, try the meta tags
    const ogImage = $('meta[property="og:image"]').attr('content');
    if (ogImage) {
      return ogImage;
    }
    
    // If we still can't find an image, look for URLs in the HTML that appear to be images
    const htmlStr = $.html();
    const imgRegex = /https?:\/\/[^"']+\.(jpg|jpeg|png|gif|webp)/gi;
    const matches = htmlStr.match(imgRegex);
    if (matches && matches.length > 0) {
      // Find the first URL that contains the word "01as" which is common in the images you showed
      const animeImg = matches.find(url => url.includes('01as'));
      if (animeImg) {
        return animeImg;
      }
      // Otherwise just return the first one that looks like a content image (not a UI element)
      // Filter out common UI images
      const contentImages = matches.filter(url => 
        !url.includes('newlogo') && 
        !url.includes('discord') && 
        !url.includes('qq288') &&
        !url.includes('icon')
      );
      if (contentImages.length > 0) {
        return contentImages[0];
      }
    }
    
    // If we really can't find anything, return null
    return null;
  } catch (error) {
    console.error(`Error fetching anime image from ${url}: ${error.message}`);
    return null;
  }
};

/**
 * Extract links and anime data from schedule page
 * @param {Object} $ - Cheerio object
 * @returns {Array} Extracted schedule data with anime links
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
 * Scrape the anime broadcast schedule from jadwal page
 * @returns {Promise<Object>} Object containing schedule data with images
 */
const scrapeJadwal = async () => {
  try {
    // Set headers to mimic a browser request
    const headers = {
      'User-Agent': getRandomUserAgent() || USER_AGENT
    };
    
    // URL to scrape (schedule page from your example)
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
    
    // Fetch images for all anime with improved error handling and retry logic
    for (const day of scheduleData) {
      for (const anime of day.anime_list) {
        if (anime.url) {
          try {
            // Try to get the image
            const imageUrl = await scrapeAnimeImage(anime.url);
            if (imageUrl) {
              anime.image_url = imageUrl;
            } else {
              // Retry with a delay if first attempt failed
              await new Promise(resolve => setTimeout(resolve, 500)); // Wait 500ms
              const retryImageUrl = await scrapeAnimeImage(anime.url);
              if (retryImageUrl) {
                anime.image_url = retryImageUrl;
              } else {
                // If we still failed, just use a default image
                anime.image_url = 'https://ww1.anoboy.app/wp-content/uploads/2019/02/cropped-512x512-192x192.png'; // Default AnoBoy logo
              }
            }
          } catch (err) {
            console.error(`Error fetching image for ${anime.title}: ${err.message}`);
            // Use default image on error
            anime.image_url = 'https://ww1.anoboy.app/wp-content/uploads/2019/02/cropped-512x512-192x192.png'; // Default AnoBoy logo
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