import { get } from 'svelte/store';
import { screenshotStore, processedScreenshotStore } from "$lib/stores/ScreenshotStore.ts";
import type { ExtendedScreenshot } from "../../types/ExtendedScreenshot.interface.ts";
import { getScreenshots } from '../../utils/screenshot.ts';

// Pulls screenshot from database
async function pullFromDb(limit: number): Promise<ExtendedScreenshot[] | undefined> {
  try {
    const res = await getScreenshots(limit);
    return res;
  } catch (err) {
    console.error(`Could not get screenshots with limit ${limit}`, err);
    return undefined;
  }
}

// Tries to find the screenshot in the store
function pullFromStore(store: ExtendedScreenshot[]): ExtendedScreenshot[] | undefined {
  if (!store) return undefined;
  return store;
}

/** @type {import('./$types').PageLoad} */
export const load = async () => {
  return {
    streamed: {
      items: new Promise(async (resolve, reject) => {
        const allScr = get(screenshotStore);
        let result = pullFromStore(allScr);

        // If not in store, try fetching from database
        if (!result) {
          result = await pullFromDb(30);
        }

        // Handle cases where result is not found
        if (!result) {
          console.log(result);
          reject(0);
        } else {
          processedScreenshotStore.subscribe(val => resolve(val))
        }
      })
    }
  };
};
