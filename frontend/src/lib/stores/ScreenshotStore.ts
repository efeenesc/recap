import { type Writable, writable } from "svelte/store";
import { derived } from "svelte/store";
import { formatDate } from "../../utils/timeSince.ts";
import type {
    DatedScreenshot,
    ExtendedScreenshot,
} from "../../types/ExtendedScreenshot.interface.ts";

export const screenshotStore = writable<ExtendedScreenshot[]>();
export const processedScreenshotStore = derived<
    Writable<ExtendedScreenshot[]>,
    DatedScreenshot | undefined
>(screenshotStore, ($scr) => {
    const screenshots: { [key: string]: ExtendedScreenshot[] } = {};
    if (!$scr) return;
    $scr.map((s: ExtendedScreenshot) => {
        s.Date = formatDate(s.Timestamp);
        if (!Object.hasOwn(screenshots, s.Date)) {
            screenshots[s.Date] = [];
        }

        screenshots[s.Date].push(s);
        return s;
    });
    return screenshots;
});
