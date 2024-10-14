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
    const reports: DatedReport = {};
    if (!$rep || !$rep.length) return;

    $rep.forEach((r: ExtendedReport) => {
        r.Date = formatDate(r.Timestamp);

        try {
            r.ParsedMarkdown = ConvertToHtmlTree(r.Content).content as MdNode[];
        } catch (err) {
            console.log(ConvertToHtmlTree(r.Content));
            console.error(err);
        }

        if (!Object.prototype.hasOwnProperty.call(reports, r.Date)) {
            reports[r.Date] = [];
        }

        reports[r.Date].push(r);
    });

    return reports;
});
