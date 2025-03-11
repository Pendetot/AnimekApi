const UserAgent = require('user-agents');

/**
 * Get a random user agent
 * @returns {string} A random user agent string
 */
const getRandomUserAgent = () => {
  return new UserAgent().toString();
};

/**
 * Extract episode number from title
 * @param {string} title - The episode title
 * @returns {string} The extracted episode number or "Unknown"
 */
const extractEpisodeNumber = (title) => {
  const match = title.match(/Episode\s+(\d+)/i);
  return match ? match[1] : "Unknown";
};

/**
 * Clean title for use in filenames
 * @param {string} title - The title to clean
 * @returns {string} The cleaned title
 */
const cleanTitle = (title) => {
  return title
    .replace(/[^\w\s-]/g, '')
    .replace(/\s+/g, '_');
};

/**
 * Format current date to YYYY-MM-DD
 * @returns {string} Formatted date
 */
const getCurrentFormattedDate = () => {
  const date = new Date();
  return date.toISOString().split('T')[0];
};

/**
 * Creates an error response object
 * @param {string} message - Error message
 * @param {number} statusCode - HTTP status code
 * @returns {Object} Formatted error object
 */
const createErrorResponse = (message, statusCode = 500) => {
  return {
    status: 'error',
    statusCode,
    message
  };
};

/**
 * Creates a success response object
 * @param {any} data - Response data
 * @param {string} message - Success message
 * @returns {Object} Formatted success object
 */
const createSuccessResponse = (data, message = 'Success') => {
  return {
    status: 'success',
    message,
    data
  };
};

module.exports = {
  getRandomUserAgent,
  extractEpisodeNumber,
  cleanTitle,
  getCurrentFormattedDate,
  createErrorResponse,
  createSuccessResponse
};