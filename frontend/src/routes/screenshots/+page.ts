import { get, writable } from "svelte/store";
import {
    screenshotStore,
    processedScreenshotStore,
} from "$lib/stores/ScreenshotStore.ts";
import type { ExtendedScreenshot } from "../../types/ExtendedScreenshot.interface.ts";
import {
    getScreenshots,
    getScreenshotsNewerThan,
} from "../../utils/screenshot.ts";

// Pulls screenshot from database
async function pullFromDb(
    limit: number
): Promise<ExtendedScreenshot[] | undefined> {
    try {
        const res = await getScreenshots(limit);
        return res;
    } catch (err) {
        console.error(`Could not get screenshots with limit ${limit}`, err);
        return undefined;
    }
}

// Tries to find the screenshot in the store
async function pullFromStore(
    store: ExtendedScreenshot[]
): Promise<ExtendedScreenshot[] | undefined> {
    if (!store) return undefined;

    let oldest = 0;
    store.forEach((scr) => {
        if (scr.CaptureID > oldest) oldest = scr.CaptureID;
    });

    const newScreenshots = await getScreenshotsNewerThan(oldest);
    if (!newScreenshots) return get(screenshotStore);

    screenshotStore.update((prev) => {
        if (!prev) return prev;

        return [...newScreenshots, ...prev];
    });

    return get(screenshotStore);
}

/** @type {import('./$types').PageLoad} */
export const load = async () => {
    const allScr = get(screenshotStore);
    let result = await pullFromStore(allScr);

    // If not in store, try fetching from database
    if (!result) {
        result = await pullFromDb(30);
    }
    
    return {
        streamed: {
            items: {
                subscribe: processedScreenshotStore.subscribe
            }
        },
    };
};
