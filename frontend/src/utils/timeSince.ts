export function timeSinceUNIXSeconds(date: number) {
    const seconds = Math.floor((Date.now() - date * 1000) / 1000); // date * 1000 converts a UNIX second timestamp to millis
    let interval = seconds / 31536000;

    if (interval > 1) {
        return Math.floor(interval) + " years";
    }
    interval = seconds / 2592000;
    if (interval > 1) {
        return Math.floor(interval) + " months";
    }
    interval = seconds / 86400;
    if (interval > 1) {
        return Math.floor(interval) + " days";
    }
    interval = seconds / 3600;
    if (interval > 1) {
        return Math.floor(interval) + " hours";
    }
    interval = seconds / 60;
    if (interval > 1) {
        return Math.floor(interval) + " minutes";
    }
    return Math.floor(seconds) + " seconds";
}

export function formatDate(unixSeconds: number): string {
    const date = new Date(unixSeconds * 1000);
    const now = new Date();
    const yesterday = new Date(now);
    yesterday.setDate(yesterday.getDate() - 1);

    if (date.toDateString() === now.toDateString()) {
        return "Today";
    } else if (date.toDateString() === yesterday.toDateString()) {
        return "Yesterday";
    } else {
        const day = date.getDate();
        const month = date.toLocaleString('default', { month: 'long' });
        const year = date.getFullYear();
        const currentYear = now.getFullYear();

        if (year === currentYear) {
            return `${day} ${month}`;
        } else {
            return `${day} ${month} ${year}`;
        }
    }
}