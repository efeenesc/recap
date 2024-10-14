import { error } from "@sveltejs/kit";
import { get } from 'svelte/store';
import type { ExtendedReport } from "../../../types/ExtendedReport.interface.ts";
import { getReportById } from "../../../utils/report.ts";
import { reportStore } from "$lib/stores/ReportStore.ts";
import { ConvertToHtmlTree } from "$lib/markdown/Markdown.ts";
import type { MdNode } from "$lib/markdown/Markdown.interface.ts";

// Pulls report from database
async function pullFromDb(id: number): Promise<ExtendedReport | null> {
  try {
    const res = await getReportById(id);
    return res
  } catch (err) {
    console.error(`Failed to fetch screenshot with ID ${id}`, err);
    return null;
  }
}

// Tries to find the report in the store
function pullFromStore(id: number, store: ExtendedReport[]): ExtendedReport | null {
  return store?.find((r) => r.ReportID === id) || null;
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
        const allScr = get(reportStore);
        let result = pullFromStore(scrId, allScr);

        // If not in store, try fetching from database
        if (!result) {
          result = await pullFromDb(scrId);
        }

        // Handle cases where result is not found
        if (!result) {
          reject(error(404, 'Report not found'));
        } else {
          if (!result.ParsedMarkdown) {
            result.ParsedMarkdown = ConvertToHtmlTree(result.Content).content as MdNode[];
          } 
          resolve(result);
        }
      })
    }
  };
};
