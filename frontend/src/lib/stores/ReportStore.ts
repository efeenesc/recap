import { derived, type Writable, writable } from "svelte/store";
import type { DatedReport, ExtendedReport } from "../../types/ExtendedReport.interface.ts";
import { formatDate } from "../../utils/timeSince.ts";
import type { MdNode } from "$lib/markdown/Markdown.interface.ts";
import { ConvertToHtmlTree } from "$lib/markdown/Markdown.ts"

export const reportStore = writable<ExtendedReport[]>();
export const processedReportStore = derived<
    Writable<ExtendedReport[]>,
    DatedReport | undefined
>(reportStore, ($rep) => {
    const reports: DatedReport = {}
    if (!$rep) return;
    $rep.map((r: ExtendedReport) => {
        r.Date = formatDate(r.Timestamp);
        r.ParsedMarkdown = ConvertToHtmlTree(r.Content).content as MdNode[];
        
        if (!Object.hasOwn(reports, r.Date)) {
            reports[r.Date] = [];
        }

        reports[r.Date].push(r);
        return r;
    });
    return reports;
});
