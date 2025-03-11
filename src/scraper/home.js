const axios = require('axios');
const cheerio = require('cheerio');
const { BASE_URL, USER_AGENT, REQUEST_TIMEOUT } = require('../config/config');
const { getRandomUserAgent } = require('../utils/helpers');

/**
 * Scrape the homepage of anoBoy website to get the latest anime episodes
 * @returns {Promise<Array>} Array of anime objects with title, update time, image URL, and episode URL
 */
const scrapeAnoboyHome = async () => {
  try {
    // Set headers to mimic a browser request
    const headers = {
      'User-Agent': getRandomUserAgent() || USER_AGENT
    };
    
    // Send HTTP request
    const response = await axios.get(BASE_URL, { 
      headers,
      timeout: REQUEST_TIMEOUT
    });
    
    // Parse HTML content
    const $ = cheerio.load(response.data);
    
    // Find all anime entries (each in a div with class "amv")
    const animeEntries = $('.amv');
    
    // Create an array to store the extracted data
    const animeData = [];
    
    // Extract information from each entry
    animeEntries.each((index, element) => {
      // Find the title
      const titleElement = $(element).find('h3.ibox1');
      const title = titleElement.length ? titleElement.text().trim() : "No title found";
      
      // Find the update time
      const timeElement = $(element).find('div.jamup');
      const updateTime = timeElement.length ? timeElement.text().trim() : "No update time found";
      
      // Find the image
      const imgElement = $(element).find('amp-img');
      let imgSrc = imgElement.length ? imgElement.attr('src') : "No image found";
      
      // Find the URL to the episode page
      const parentA = $(element).closest('a');
      const episodeUrl = parentA.length ? parentA.attr('href') : null;
      
      // Convert relative image URL to absolute URL
      if (imgSrc && !imgSrc.startsWith(('http://', 'https://'))) {
        const url = new URL(imgSrc, BASE_URL);
        imgSrc = url.href;
      }
      
      // Add to our data array
      animeData.push({
        title,
        update_time: updateTime,
        image_url: imgSrc,
        episode_url: episodeUrl
      });
    });
    
    return animeData;
    
  } catch (error) {
    console.error(`Error fetching the homepage: ${error.message}`);
    throw new Error(`Failed to scrape AnoBoy homepage: ${error.message}`);
  }
};

module.exports = {
  scrapeAnoboyHome
};