const axios = require('axios');
const cheerio = require('cheerio');
const { USER_AGENT, REQUEST_TIMEOUT } = require('../config/config');
const { getRandomUserAgent } = require('../utils/helpers');

/**
 * Scrape details from an anime episode page
 * @param {string} url - The URL of the anime episode page
 * @returns {Promise<Object>} Object containing anime details
 */
const scrapeAnimeDetail = async (url) => {
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
    
    // Get anime description
    const descElement = $('div.contentdeks');
    const description = descElement.length ? descElement.text().trim() : "Description not found";
    
    // Get anime image
    const imgElement = $('amp-img.gambar');
    const imgSrc = imgElement.length ? imgElement.attr('src') : "Image not found";
    
    // Get streaming sources
    const streamingLinks = [];
    
    // Get main player iframe source
    const playerIframe = $('iframe#mediaplayer');
    if (playerIframe.length && playerIframe.attr('src')) {
      streamingLinks.push({
        name: "Main Player",
        url: playerIframe.attr('src'),
        quality: "Default" // Default quality
      });
    }
    
    // Get alternative streaming sources
    const altPlayers = $('a#allmiror');
    altPlayers.each((index, element) => {
      if ($(element).attr('data-video')) {
        // Try to extract quality info from text
        const playerText = $(element).text().trim() || "Alternative Player";
        let quality = "Unknown";
        
        // Look for common resolution patterns like "360-720" or "240P-1080P"
        const resMatch = playerText.match(/(\d+[pP]?[-]?\d*[pP]?)/);
        if (resMatch) {
          quality = resMatch[1];
        } else {
          // Try to find numbers that might indicate resolution
          const numbers = playerText.match(/\d+/g);
          if (numbers && numbers.length) {
            quality = numbers.join('-') + 'P';
          }
        }
        
        streamingLinks.push({
          name: playerText,
          url: $(element).attr('data-video'),
          quality: quality
        });
      }
    });
    
    // Get download links
    const downloadLinks = [];
    const downloadSection = $('div.download');
    if (downloadSection.length) {
      const linkSpans = downloadSection.find('span.ud');
      linkSpans.each((index, span) => {
        const serverName = $(span).find('span.udj');
        const server = serverName.length ? serverName.text().trim() : "Unknown Server";
        
        const qualityLinks = $(span).find('a.udl');
        qualityLinks.each((i, link) => {
          const qualityLink = $(link);
          if (qualityLink.attr('href') !== 'none' && !qualityLink.attr('style')?.includes('display:none')) {
            downloadLinks.push({
              server: server,
              quality: qualityLink.text().trim(),
              url: qualityLink.attr('href')
            });
          }
        });
      });
    }
    
    // Get metadata
    const metadata = {};
    const infoTable = $('div.contenttable');
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
    
    return {
      title,
      description,
      image_url: imgSrc,
      streaming_links: streamingLinks,
      download_links: downloadLinks,
      metadata
    };
    
  } catch (error) {
    console.error(`Error fetching the anime detail page: ${error.message}`);
    throw new Error(`Failed to scrape anime details: ${error.message}`);
  }
};

module.exports = {
  scrapeAnimeDetail
};