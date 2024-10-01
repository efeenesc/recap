import { get } from 'svelte/store';
import { reportStore, processedReportStore } from "$lib/stores/ReportStore.ts";
import type { ExtendedReport } from '../../types/ExtendedReport.interface.ts';
import { getReports } from '../../utils/report.ts';

// Pulls screenshot from database
async function pullFromDb(limit: number): Promise<ExtendedReport[] | undefined> {
  try {
    const res = await getReports(limit);
    return res;
  } catch (err) {
    console.error(`Could not get reports with limit ${limit}`, err);
    return undefined;
  }
}

// Tries to find the screenshot in the store
function pullFromStore(store: ExtendedReport[]): ExtendedReport[] | undefined {
  if (!store) return undefined;
  return store;
}

/** @type {import('./$types').PageLoad} */
export const load = async () => {
  return {
    streamed: {
      items: new Promise(async (resolve, reject) => {
        const allScr = get(reportStore);
        let result = pullFromStore(allScr);

        // If not in store, try fetching from database
        if (!result) {
          result = await pullFromDb(30);
        }

        // Handle cases where result is not found
        if (!result) {
          reject(0);
        } else {
          processedReportStore.subscribe(val => resolve(val))
        }
      })
    }
  };
};
