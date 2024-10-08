import { get } from "svelte/store";
import { reportStore, processedReportStore } from "$lib/stores/ReportStore.ts";
import type { ExtendedReport } from "../../types/ExtendedReport.interface.ts";
import { getReports, getReportsNewerThan } from "../../utils/report.ts";
import { navigating } from "$app/stores";

// Pulls screenshot from database
async function pullFromDb(
    limit: number
): Promise<ExtendedReport[] | undefined> {
    try {
        const res = await getReports(limit);
        return res;
    } catch (err) {
        console.error(`Could not get reports with limit ${limit}`, err);
        return undefined;
    }
}

// Tries to find the screenshot in the store
async function pullFromStore(
    store: ExtendedReport[]
): Promise<ExtendedReport[] | undefined> {
    if (!store) return undefined;

    let oldest = 0;
    store.forEach((scr) => {
        if (scr.ReportID > oldest) oldest = scr.ReportID;
    });

    const newScreenshots = await getReportsNewerThan(oldest);

    reportStore.update((prev) => {
        if (!prev) return prev;

        let newVal;
        if (!newScreenshots) newVal = [...prev];
        else newVal = [...newScreenshots, ...prev];

        // If navigating back from /screenshots/[id], don't trim the screenshot store down to 30 items
        if (get(navigating)!.from!.route.id === "/reports/[id]")
            return newVal;

        return newVal.splice(0, 30);
    });

    return get(reportStore);
}

/** @type {import('./$types').PageLoad} */
export const load = async () => {
    const allScr = get(reportStore);
    let result = await pullFromStore(allScr);

    // If not in store, try fetching from database
    if (!result) {
        result = await pullFromDb(30);
    }

    console.log(get(processedReportStore));

    return {
        data: get(processedReportStore)
    };
};
