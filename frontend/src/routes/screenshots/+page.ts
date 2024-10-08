import { get } from "svelte/store";
import {
    screenshotStore,
    processedScreenshotStore,
} from "$lib/stores/ScreenshotStore.ts";
import type { ExtendedScreenshot } from "../../types/ExtendedScreenshot.interface.ts";
import {
    getScreenshots,
    getScreenshotsNewerThan,
} from "../../utils/screenshot.ts";
import { navigating } from '$app/stores'
import type { PageLoad } from "./$types.js";

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

async function pullFromStore(
    store: ExtendedScreenshot[]
): Promise<ExtendedScreenshot[] | undefined> {
    if (!store) return undefined;

    let oldest = 0;
    store.forEach((scr) => {
        if (scr.CaptureID > oldest) oldest = scr.CaptureID;
    });

    const newScreenshots = await getScreenshotsNewerThan(oldest);

    screenshotStore.update((prev) => {
        if (!prev) return prev;

        let newVal;
        if (!newScreenshots) newVal = [...prev];
        else newVal = [...newScreenshots, ...prev];

        // If navigating back from /screenshots/[id], don't trim the screenshot store down to 30 items
        if (get(navigating)!.from!.route.id === "/screenshots/[id]")
            return newVal;

        return newVal.splice(0, 30);
    });

    return get(screenshotStore);
}

/** @type {import('./$types').PageLoad} */
export const load: PageLoad = async () => {
    const allScr = get(screenshotStore);
    let result = await pullFromStore(allScr);

    // If not in store, try fetching from database
    if (!result) {
        result = await pullFromDb(30);
    }

    return {
        data: get(processedScreenshotStore)
    };
};
