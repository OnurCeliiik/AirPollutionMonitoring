// Constants for air quality levels based on WHO guidelines
export const AQ_LEVELS = {
  PM25: {
    GOOD: 10,        // WHO guideline
    MODERATE: 25,    // 2.5x guideline
    UNHEALTHY: 55,   // 5.5x guideline
    VERY_UNHEALTHY: 150, // 15x guideline
    HAZARDOUS: 250,  // 25x guideline
  },
  PM10: {
    GOOD: 20,        // WHO guideline
    MODERATE: 50,    // 2.5x guideline
    UNHEALTHY: 100,  // 5x guideline
    VERY_UNHEALTHY: 200, // 10x guideline
    HAZARDOUS: 400,  // 20x guideline
  },
  O3: {
    GOOD: 100,       // WHO guideline (8-hour mean)
    MODERATE: 160,   // 1.6x guideline
    UNHEALTHY: 200,  // 2x guideline
    VERY_UNHEALTHY: 300, // 3x guideline
    HAZARDOUS: 400,  // 4x guideline
  },
  NO2: {
    GOOD: 40,        // WHO guideline (annual mean)
    MODERATE: 100,   // 2.5x guideline
    UNHEALTHY: 200,  // 5x guideline
    VERY_UNHEALTHY: 400, // 10x guideline
    HAZARDOUS: 600,  // 15x guideline
  },
  SO2: {
    GOOD: 40,        // WHO guideline (24-hour mean)
    MODERATE: 100,   // 2.5x guideline
    UNHEALTHY: 200,  // 5x guideline
    VERY_UNHEALTHY: 400, // 10x guideline
    HAZARDOUS: 600,  // 15x guideline
  },
};

// Get color based on pollution level
export const getColorForValue = (parameter, value) => {
  // Default to PM2.5 if parameter not found
  const levels = AQ_LEVELS[parameter] || AQ_LEVELS.PM25;
  
  if (value <= levels.GOOD) {
    return '#00E400'; // Green
  } else if (value <= levels.MODERATE) {
    return '#FFFF00'; // Yellow
  } else if (value <= levels.UNHEALTHY) {
    return '#FF7E00'; // Orange
  } else if (value <= levels.VERY_UNHEALTHY) {
    return '#FF0000'; // Red
  } else {
    return '#7F0023'; // Dark Red
  }
};

// Get quality description based on pollution level
export const getQualityDescription = (parameter, value) => {
  // Default to PM2.5 if parameter not found
  const levels = AQ_LEVELS[parameter] || AQ_LEVELS.PM25;
  
  if (value <= levels.GOOD) {
    return 'Good';
  } else if (value <= levels.MODERATE) {
    return 'Moderate';
  } else if (value <= levels.UNHEALTHY) {
    return 'Unhealthy';
  } else if (value <= levels.VERY_UNHEALTHY) {
    return 'Very Unhealthy';
  } else {
    return 'Hazardous';
  }
};

// Format timestamp to readable date/time
export const formatTimestamp = (timestamp) => {
  const date = new Date(timestamp);
  return date.toLocaleString();
};

// Calculate how long ago a timestamp was
export const timeAgo = (timestamp) => {
  const seconds = Math.floor((new Date() - new Date(timestamp)) / 1000);
  
  let interval = Math.floor(seconds / 31536000);
  if (interval > 1) {
    return `${interval} years ago`;
  }
  if (interval === 1) {
    return '1 year ago';
  }
  
  interval = Math.floor(seconds / 2592000);
  if (interval > 1) {
    return `${interval} months ago`;
  }
  if (interval === 1) {
    return '1 month ago';
  }
  
  interval = Math.floor(seconds / 86400);
  if (interval > 1) {
    return `${interval} days ago`;
  }
  if (interval === 1) {
    return '1 day ago';
  }
  
  interval = Math.floor(seconds / 3600);
  if (interval > 1) {
    return `${interval} hours ago`;
  }
  if (interval === 1) {
    return '1 hour ago';
  }
  
  interval = Math.floor(seconds / 60);
  if (interval > 1) {
    return `${interval} minutes ago`;
  }
  if (interval === 1) {
    return '1 minute ago';
  }
  
  return 'just now';
}; 