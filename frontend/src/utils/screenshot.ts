import { screenshotStore } from "$lib/stores/ScreenshotStore.ts";
import { GetScreenshotById, GetScreenshots, GetScreenshotsOlderThan, GetScreenshotsNewerThan } from "$lib/wailsjs/go/app/AppMethods.js";
import type { ExtendedScreenshot } from "../types/ExtendedScreenshot.interface.ts";
import { db } from "$lib/wailsjs/go/models.ts";

const processScreenshots = (
    screenshots: db.CaptureScreenshotImage[] | null
): ExtendedScreenshot[] => {
    if (!screenshots) return [];
    const cast = screenshots.map<ExtendedScreenshot>(
        (v: any) => {
            v.Time = new Date(v.Timestamp * 1000).toLocaleTimeString();
            return v;
        }
    );
    return cast;
};

export const getScreenshots = async (limit: number) => {
    const res = await GetScreenshots(limit);
    const processed = processScreenshots(res);
    screenshotStore.update((val) => (val = processed));
    return res;
};

export const getScreenshotById = async (id: number): Promise<ExtendedScreenshot> => {
    const res = await GetScreenshotById(id);
    const processed = processScreenshots([res]);
    return processed[0];
};

export const getScreenshotsOlderThan = async (id: number, limit: number) => {
    const res = await GetScreenshotsOlderThan(id, limit);

    const processed = processScreenshots(res);
    return processed;
}

export const getScreenshotsNewerThan = async (id: number) => {
    const res = await GetScreenshotsNewerThan(id);
    if (!res) return null;
    const processed = processScreenshots(res);
    return processed;
}

export const addScreenshotToStore = async (id: number): Promise<ExtendedScreenshot> => {
    const res = await getScreenshotById(id);
    const newScreenshot = res;

    screenshotStore.update((prev) => {
        return [newScreenshot, ...prev];
    });

    return newScreenshot;
}
