import { error } from "@sveltejs/kit";
import { get } from 'svelte/store';
import { screenshotStore } from "$lib/stores/ScreenshotStore.ts";
import { GetScreenshotById } from "$lib/wailsjs/go/app/AppMethods.js";
import type { ExtendedScreenshot } from "../../../types/ExtendedScreenshot.interface.ts";
import { getScreenshotById } from "../../../utils/screenshot.ts";

// Pulls screenshot from database
async function pullFromDb(id: number): Promise<ExtendedScreenshot | null> {
  try {
    const res = await getScreenshotById(id);
    return res;
  } catch (err) {
    console.error(`Failed to fetch screenshot with ID ${id}`, err);
    return null;
  }
}

// Tries to find the screenshot in the store
function pullFromStore(id: number, store: ExtendedScreenshot[]): ExtendedScreenshot | null {
  return store?.find((s) => s.CaptureID === id) || null;
}

/** @type {import('./$types').PageLoad} */
export const load = async ({ params }) => {
  const scrId = Number(params.id);

  if (isNaN(scrId)) {
    throw error(404, 'Invalid screenshot ID');
  }

  return {
    streamed: {
      items: new Promise(async (resolve, reject) => {
        const result = await pullFromDb(scrId);

        // Handle cases where result is not found
        if (!result) {
          reject(error(404, 'Screenshot not found'));
        } else {
          resolve(result); // Pass the screenshot result
        }
      })
    }
  };
};
