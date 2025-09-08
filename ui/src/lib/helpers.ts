// Function to calculate time ago
export function getTimeAgo(createdAt: string | Date): string {
  const now = new Date();
  const createdDate = new Date(createdAt);
  const diffInMs = now.getTime() - createdDate.getTime();

  const diffInMinutes = Math.floor(diffInMs / (1000 * 60));
  const diffInHours = Math.floor(diffInMs / (1000 * 60 * 60));
  const diffInDays = Math.floor(diffInMs / (1000 * 60 * 60 * 24));
  const diffInWeeks = Math.floor(diffInDays / 7);
  const diffInMonths = Math.floor(diffInDays / 30);
  const diffInYears = Math.floor(diffInDays / 365);

  if (diffInMinutes < 1) {
    return "now";
  } else if (diffInMinutes < 60) {
    return `${diffInMinutes} min${diffInMinutes > 1 ? "s" : ""} ago`;
  } else if (diffInHours < 24) {
    return `${diffInHours} hr${diffInHours > 1 ? "s" : ""} ago`;
  } else if (diffInDays < 7) {
    return `${diffInDays} day${diffInDays > 1 ? "s" : ""} ago`;
  } else if (diffInWeeks < 4) {
    return `${diffInWeeks} week${diffInWeeks > 1 ? "s" : ""} ago`;
  } else if (diffInMonths < 12) {
    return `${diffInMonths} month${diffInMonths > 1 ? "s" : ""} ago`;
  } else {
    return `${diffInYears} year${diffInYears > 1 ? "s" : ""} ago`;
  }
}

const urlRegex =
  /(\b(https?|ftp|file):\/\/[-A-Z0-9+&@#\/%?=~_|!:,.;]*[-A-Z0-9+&@#\/%=~_|])|(\bwww\.[-A-Z0-9+&@#\/%?=~_|!:,.;]*[-A-Z0-9+&@#\/%=~_|])|(\b[-A-Z0-9+&@#\/%?=~_|!:,.;]*[-A-Z0-9+&@#\/%=~_|]\.(com|org|net|edu|gov|mil|biz|info|mobi|name|aero|jobs|museum|coop|asia|cat|int|io|pro|tel|travel|xxx))\b/gi;

export function linkify(text: string): string {
  return text.replace(urlRegex, function (url) {
    if (url.startsWith("www.")) {
      return (
        '<a href="http://' +
        url +
        '" class="link" target="_blank" rel="noopener noreferrer">' +
        url +
        "</a>"
      );
    } else if (
      !url.startsWith("http") &&
      !url.startsWith("ftp") &&
      !url.startsWith("file")
    ) {
      return (
        '<a href="http://' +
        url +
        '" class="link" target="_blank" rel="noopener noreferrer">' +
        url +
        "</a>"
      );
    } else {
      return (
        '<a href="' +
        url +
        '" class="link" target="_blank" rel="noopener noreferrer">' +
        url +
        "</a>"
      );
    }
  });
}

export function extractLinks(text: string): string[] {
  const matches = text.match(urlRegex);
  return matches ? matches : [];
}
